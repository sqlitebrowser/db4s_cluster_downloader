package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/minio/go-homedir"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	DB4S_3_10_1_WIN32 = iota // The order needs to match the ramCache entries in the global var section
	DB4S_3_10_1_WIN64
	DB4S_3_10_1_OSX
	DB4S_3_10_1_PORTABLE
)

// Configuration file
type TomlConfig struct {
	Jaeger JaegerInfo
	Paths  PathInfo
	Pg     PGInfo
	Server ServerInfo
}
type JaegerInfo struct {
	CollectorEndPoint string
}
type PathInfo struct {
	BaseDir string // Location of the git source
	CertDir string // Location of the TLS certificates
	DataDir string // Directory where the downloads are located
}
type PGInfo struct {
	Database       string
	NumConnections int `toml:"num_connections"`
	Port           int
	Password       string
	Server         string
	SSL            bool
	Username       string
}
type ServerInfo struct {
	Port int
}

// dbEntry is used for storing the new database entries
type dbEntry struct {
	ipv4      pgtype.Text
	ipv6      pgtype.Text
	ipstrange pgtype.Text
	port      pgtype.Int4
}

// Stores each of the files available for download, along with its metadata
type cacheEntry struct {
	lastRFC1123 string // Pre-rendered string
	disposition string // Pre-rendered string
	mem         []byte
	ready       bool
	reader      *bytes.Reader
	size        string // Pre-rendered string
}

var (
	// Application config
	Conf TomlConfig

	// Use Jaeger?
	enableJaeger = false

	// PostgreSQL Connection pool
	pg *pgx.ConnPool

	// Cached downloads
	ramCache = [4]cacheEntry{
		// These hard coded last modified timestamps are because we're working with existing files, so we provide the
		// same ones already being used
		{ // DB4S 3.10.1 Win32
			lastRFC1123: time.Date(2017, time.September, 20, 14, 59, 44, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.10.1-win32.exe"),
				time.Date(2017, time.September, 20, 14, 59, 44, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.10.1 Win64
			lastRFC1123: time.Date(2017, time.September, 20, 14, 59, 59, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.10.1-win64.exe"),
				time.Date(2017, time.September, 20, 14, 59, 59, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.10.1 OSX
			lastRFC1123: time.Date(2017, time.September, 20, 15, 23, 27, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.10.1.dmg"),
				time.Date(2017, time.September, 20, 15, 23, 27, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.10.1 PortableApp
			lastRFC1123: time.Date(2017, time.September, 28, 19, 32, 48, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe"),
				time.Date(2017, time.September, 28, 19, 32, 48, 0, time.UTC).Format(time.RFC3339)),
		},
	}
	tracer opentracing.Tracer
	tmpl   *template.Template
)

func main() {
	// Override config file location via environment variables
	var err error
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		// TODO: Might be a good idea to add permission checks of the dir & conf file, to ensure they're not
		//       world readable.  Similar in concept to what ssh does for its config files.
		userHome, err := homedir.Dir()
		if err != nil {
			log.Fatalf("User home directory couldn't be determined: %s", "\n")
		}
		configFile = filepath.Join(userHome, ".db4s", "downloader_config.toml")
	}

	// Read our configuration settings
	if _, err = toml.DecodeFile(configFile, &Conf); err != nil {
		log.Fatal(err)
	}

	// Set up initial Jaeger service and span
	var closer io.Closer
	tracer, closer = initJaeger("db4s_cluster_downloader")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// * Connect to PG database *

	pgSpan := tracer.StartSpan("connect postgres")

	// Setup the PostgreSQL config
	pgConfig := new(pgx.ConnConfig)
	pgConfig.Host = Conf.Pg.Server
	pgConfig.Port = uint16(Conf.Pg.Port)
	pgConfig.User = Conf.Pg.Username
	pgConfig.Password = Conf.Pg.Password
	pgConfig.Database = Conf.Pg.Database
	clientTLSConfig := tls.Config{InsecureSkipVerify: true}
	if Conf.Pg.SSL {
		// TODO: Likely need to add the PG TLS cert file info here
		pgConfig.TLSConfig = &clientTLSConfig
	} else {
		pgConfig.TLSConfig = nil
	}

	// Connect to PG
	pgPoolConfig := pgx.ConnPoolConfig{*pgConfig, Conf.Pg.NumConnections, nil, 5 * time.Second}
	pg, err = pgx.NewConnPool(pgPoolConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Log successful connection
	log.Printf("Connected to PostgreSQL server: %v:%v\n", Conf.Pg.Server, uint16(Conf.Pg.Port))
	pgSpan.Finish()

	// Load our HTML template
	// TODO: Embed the template in the compiled binary
	tmpl = template.Must(template.New("downloads").ParseFiles(filepath.Join(Conf.Paths.BaseDir, "template.html"))).Lookup("downloads")

	// Load the files into ram from the data directory
	cacheSpan := tracer.StartSpan("create cache entries")
	ctx := opentracing.ContextWithSpan(context.Background(), cacheSpan)
	ramCache[DB4S_3_10_1_WIN32].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.10.1-win32.exe"))
	if err == nil {
		cache(ctx, ramCache[DB4S_3_10_1_WIN32])
	}
	ramCache[DB4S_3_10_1_WIN64].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.10.1-win64.exe"))
	if err == nil {
		cache(ctx, ramCache[DB4S_3_10_1_WIN64])
	}
	ramCache[DB4S_3_10_1_OSX].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.10.1.dmg"))
	if err == nil {
		cache(ctx, ramCache[DB4S_3_10_1_OSX])
	}
	ramCache[DB4S_3_10_1_PORTABLE].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe"))
	if err == nil {
		cache(ctx, ramCache[DB4S_3_10_1_PORTABLE])
	}
	cacheSpan.Finish()

	http.HandleFunc("/", handler)
	fmt.Printf("Listening on port %d...\n", Conf.Server.Port)
	err = http.ListenAndServeTLS(fmt.Sprintf(":%d", Conf.Server.Port), filepath.Join(Conf.Paths.CertDir, "fullchain.pem"), filepath.Join(Conf.Paths.CertDir, "privkey.pem"), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Close the PG connection gracefully
	pg.Close()
}

// Populates a cache entry
func cache(ctx context.Context, cacheEntry cacheEntry) {
	span, _ := opentracing.StartSpanFromContext(ctx, "populate cache entry")
	defer span.Finish()
	span.SetTag("Entry", cacheEntry.disposition)

	cacheEntry.reader = bytes.NewReader(cacheEntry.mem)
	cacheEntry.size = fmt.Sprintf("%d", len(cacheEntry.mem))
	cacheEntry.ready = true
}

// Handler for download requests
func handler(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpan("page handler")
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	var err error
	switch r.URL.Path {
	case "/favicon.ico":
		span.SetTag("Request", "favicon.ico")
		http.ServeFile(w, r, filepath.Join(Conf.Paths.BaseDir, "favicon.ico"))
		err = logRequest(ctx, r, 90022, http.StatusOK) // 90022 is the file size of favicon.ico in bytes
		if err != nil {
			log.Printf("Error: %s", err)
		}
	case "/currentrelease":
		span.SetTag("Request", "currentrelease")
		bytesSent, err := fmt.Fprint(w, "3.10.1\nhttps://github.com/sqlitebrowser/sqlitebrowser/releases/tag/v3.10.1\n")
		if err != nil {
			log.Printf("Error serving currentrelease: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = logRequest(ctx, r, int64(bytesSent), http.StatusOK)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	case "/DB.Browser.for.SQLite-3.10.1-win64.exe":
		span.SetTag("Request", "DB.Browser.for.SQLite-3.10.1-win64.exe")
		serveDownload(ctx, w, r, ramCache[DB4S_3_10_1_WIN64], "DB.Browser.for.SQLite-3.10.1-win64.exe")
	case "/DB.Browser.for.SQLite-3.10.1-win32.exe":
		span.SetTag("Request", "DB.Browser.for.SQLite-3.10.1-win32.exe")
		serveDownload(ctx, w, r, ramCache[DB4S_3_10_1_WIN32], "DB.Browser.for.SQLite-3.10.1-win32.exe")
	case "/DB.Browser.for.SQLite-3.10.1.dmg":
		span.SetTag("Request", "DB.Browser.for.SQLite-3.10.1.dmg")
		serveDownload(ctx, w, r, ramCache[DB4S_3_10_1_OSX], "DB.Browser.for.SQLite-3.10.1.dmg")
	case "/SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe":
		span.SetTag("Request", "SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe")
		serveDownload(ctx, w, r, ramCache[DB4S_3_10_1_PORTABLE], "SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe")
	default:
		span.SetTag("Request", "index page")
		if err != nil {
			_, e := fmt.Fprintf(w, "Error: %v", err)
			log.Printf("Error: %s", e)
			log.Printf("Error: %s", err)
		}
		// Send the index page listing
		err = tmpl.Execute(w, nil)
		if err != nil {
			_, e := fmt.Fprintf(w, "Error: %v", err)
			log.Printf("Error: %s", e)
			log.Printf("Error: %s", err)
		}
		err = logRequest(ctx, r, 771, http.StatusOK) // The index page is 771 bytes in length
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
}

// initJaeger returns an instance of Jaeger Tracer
func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	samplerConst := 1.0
	if !enableJaeger {
		samplerConst = 0.0
	}
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: samplerConst,
		},
		Reporter: &config.ReporterConfig{
			CollectorEndpoint: Conf.Jaeger.CollectorEndPoint,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func logRequest(ctx context.Context, r *http.Request, bytesSent int64, status int) (err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "log request")
	defer span.Finish()

	// Use the new v3 pgx/pgtype structures
	ref := &pgtype.Text{
		String: r.Referer(),
		Status: pgtype.Present,
	}
	if r.Referer() == "" {
		ref.Status = pgtype.Null
	}

	// Grab the client IP address
	clientIP := dbEntry{
		ipv4:      pgtype.Text{Status: pgtype.Null},
		ipv6:      pgtype.Text{Status: pgtype.Null},
		ipstrange: pgtype.Text{Status: pgtype.Null},
		port:      pgtype.Int4{Status: pgtype.Null},
	}
	tempIP := r.Header.Get("X-Forwarded-For")
	if tempIP == "" {
		tempIP = r.RemoteAddr
	}
	if tempIP != "" {
		// Determine if client IP address is IPv4 or IPv6, and split out the port number
		if strings.HasPrefix(tempIP, "[") {
			// * This is an IPv6 address *

			// When the string starts with "[", it seems to mean this is an IPv6 address with the port number on the
			// end.  Along the lines of (say) [1:2:3:4:5:6]:789
			s := strings.SplitN(tempIP, "]:", 2)
			ip := strings.TrimLeft(s[0], "[")

			// Check for unexpected values
			if len(s) != 2 {
				// TODO: We should probably serialise the entire request, and store that instead, to allow for better
				//       analysis of any weirdness we need to be aware of and/or adjust for

				// We either have too many, or not enough fields.  Either way, store the IP address in the "strange"
				// database field for future analysis
				clientIP.ipstrange.String = ip
				clientIP.ipstrange.Status = pgtype.Present

			} else {
				// * In theory, we should have the port number in [s]1 *

				// Convert the port string into a number
				if s[1] != "" {
					p, e := strconv.ParseInt(s[1], 10, 32)
					if e == nil {
						clientIP.port.Status = pgtype.Present
						clientIP.port.Int = int32(p)

						// Double check the port number conversion was correct
						tst := fmt.Sprintf("%d", clientIP.port.Int)
						if tst != s[1] {
							log.Printf("String conversion failed! s[1] = %v, p = %v, int32(p) = %v, port = %v\n",
								s[1], p, int32(p), clientIP.port.Int)
							clientIP.port.Status = pgtype.Null
							clientIP.port.Int = 0
						}
					} else {
						log.Printf("Conversion error: %v\n", e)
					}
				}

				// Validate the likely IPv6 address
				tmp := net.ParseIP(ip)
				if tmp == nil {
					// Not a valid IP address, so store the (complete) weird address in the client_ip_strange field
					// for future investigation.
					log.Printf("Strange address '%v'", tempIP)
					clientIP.ipstrange.String = tempIP
					clientIP.ipstrange.Status = pgtype.Present
				}

				// Double check the IP address, by seeing if converting it to a string matches what we were given
				if tmp.String() != ip {
					// Something seems a bit off with the IP address, so store it in the client_ip_strange field for
					// future investigation
					log.Printf("Strange address '%v'", tempIP)
					clientIP.ipstrange.String = tempIP
					clientIP.ipstrange.Status = pgtype.Present
				}

				clientIP.ipv6.String = ip
				clientIP.ipv6.Status = pgtype.Present
			}
		} else {
			// Client IP address doesn't seem to be in "[IPv6]:port" format, so it's likely IPv4, or IPv6 without a
			// port number.  The occasional strange value does also turn up (eg 'unknown' or multiple IP's), so we
			// need to recognise those and handle them as well.

			s := strings.SplitN(tempIP, ":", 2)
			switch len(s) {
			case 1:
				// Most likely an IPv4 address without a port number.  Validate it to make sure.
				if net.ParseIP(s[0]) == nil {
					// Not a valid IP address, so store the (complete) weird address in the client_ip_strange field
					// for future investigation.
					log.Printf("Strange address '%v'", tempIP)
					clientIP.ipstrange.String = tempIP
					clientIP.ipstrange.Status = pgtype.Present
					break
				}

				// Yep, it's just a standard IPv4 address missing a port number
				clientIP.ipv4.String = s[0]
				clientIP.ipv4.Status = pgtype.Present

			case 2:
				// Most likely an IPv4 address with a port number.  eg 1.2.3.4:56789

				// Validate the IPv4 address
				if net.ParseIP(s[0]) == nil {
					// Not a valid IP address, so store the (complete) weird address in the client_ip_strange field
					// for future investigation.
					log.Printf("Strange address '%v'", tempIP)
					clientIP.ipstrange.String = tempIP
					clientIP.ipstrange.Status = pgtype.Present
					break
				}

				// The IPv4 address is valid, so record that
				clientIP.ipv4.String = s[0]
				clientIP.ipv4.Status = pgtype.Present

				// Validate the port number
				if s[1] != "" {
					p, e := strconv.ParseInt(s[1], 10, 32)
					if e != nil {
						log.Printf("Conversion error: %v\n", e)
						break
					}

					// Double check the port number conversion was correct
					tst := fmt.Sprintf("%d", p)
					if tst != s[1] {
						log.Printf("String conversion failed! s[1] = %v, p = %v, int32(p) = %v, port = %v\n",
							s[1], p, int32(p), clientIP.port.Int)
						break
					}

					// Ensure the port number is in the valid port range (0-65535)
					if p < 0 || p > 65535 {
						log.Printf("Port number %v outside valid port range\n", p)
						break
					}

					// Port number seems ok, so record it
					clientIP.port.Status = pgtype.Present
					clientIP.port.Int = int32(p)
				}
			default:
				// * Most likely an IPv6 address without a port number.  Along the lines of (say) "1:2:3:4:5:6" *

				// Validate the address
				ip := net.ParseIP(s[0])
				if ip == nil {
					// Not a valid IP address, so store the (complete) weird address in the client_ip_strange field
					// for future investigation.
					log.Printf("Strange address '%v'", tempIP)
					clientIP.ipstrange.String = tempIP
					clientIP.ipstrange.Status = pgtype.Present
					break
				}

				// Double check the IP address, by seeing if converting it to a string matches what we were given
				if ip.String() != s[0] {
					// Something seems a bit off with the IP address, so store it in the client_ip_strange field for
					// future investigation
					log.Printf("Strange address '%v'", tempIP)
					clientIP.ipstrange.String = tempIP
					clientIP.ipstrange.Status = pgtype.Present
					break
				}

				// Seems ok so far, so lets store it in the IPv6 field
				clientIP.ipv6.String = s[0]
				clientIP.ipv6.Status = pgtype.Present
			}
		}
	} else {
		// Can't determine client IP address
		log.Printf("Unknown client IP address. :(")
	}

	dbQuery := `
		INSERT INTO download_log (
			client_ipv4, client_ipv6, client_ip_strange, client_port, remote_user, request_time, request_type, request,
			protocol, status, body_bytes_sent, http_referer, http_user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	res, err := pg.Exec(dbQuery,
		// IP address
		&clientIP.ipv4, &clientIP.ipv6, &clientIP.ipstrange,
		// client port
		&clientIP.port,
		// remote_user
		&pgtype.Text{String: "", Status: pgtype.Null}, // Hard coded empty string for now
		// request_time
		time.Now().Format(time.RFC3339Nano),
		// request_type
		r.Method,
		// request
		fmt.Sprintf("%s", r.URL),
		// protocol
		r.Proto,
		// status
		status,
		// body_bytes_sent
		bytesSent,
		// http_referer
		ref,
		// http_user_agent
		r.Header.Get("User-Agent"))
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return fmt.Errorf("something went wrong when inserting a new download entry.  # of entries affected != 1")
	}
	return
}

// Serves downloads from cache
func serveDownload(ctx context.Context, w http.ResponseWriter, r *http.Request, download cacheEntry, fileName string) {
	span, newCtx := opentracing.StartSpanFromContext(ctx, "serve download")
	defer span.Finish()
	span.SetTag("Request", fileName)

	// If the file isn't cached, check if it's ready to be cached yet
	var err error
	if !download.ready {
		download.mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, fileName))
		if err == nil {
			// TODO: It'd probably be a good idea to check the SHA256 of the file contents before marking the cache as valid
			// Add the download to the cache
			download.reader = bytes.NewReader(download.mem)
			download.size = fmt.Sprintf("%d", len(download.mem))
			download.ready = true
		}
	}

	// Send the file (if cached)
	if download.ready {
		w.Header().Set("Last-Modified", download.lastRFC1123)
		w.Header().Set("Content-Disposition", download.disposition)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", download.size)
		bytesSent, err := download.reader.WriteTo(w)
		if err != nil {
			log.Printf("Error serving %s: %v\n", fileName, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			// TODO: Catch when an error occurs (eg client aborts download), and log that too.  Probably need extra
			//       fields added to the database for holding the info.
			err = logRequest(newCtx, r, bytesSent, http.StatusBadRequest)
			if err != nil {
				log.Printf("Error: %s", err)
			}
			return
		}
		err = logRequest(newCtx, r, bytesSent, http.StatusOK)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	} else {
		// Warn the user
		_, err = fmt.Fprintf(w, "Not yet available")
		if err != nil {
			log.Printf("Error: %s", err)
		}
		err = logRequest(newCtx, r, 17, http.StatusNotFound)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
}
