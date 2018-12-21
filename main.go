package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
)

const (
	// Directories to load things from
	baseDir = "/root/git_repos/db4s_cluster_downloader"          // Location of the git source
	certDir = "/etc/letsencrypt/live/download.sqlitebrowser.org" // Location of the TLS certificates.  Shared with the host.
	dataDir = "/root/data"                                       // Directory where the downloads are located.  Shared with the host.

	// Application config settings
	configFile = "/root/data/config.toml"

	// Port to listen on
	listenPort = 443
)

const (
	DB4S_3_10_1_WIN32 = iota // The order needs to match the ramCache entries in the global var section
	DB4S_3_10_1_WIN64
	DB4S_3_10_1_OSX
	DB4S_3_10_1_PORTABLE
)

// Configuration file
type TomlConfig struct {
	Pg PGInfo
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

type download struct {
	lastRFC1123 string // Pre-rendered string
	disposition string // Pre-rendered string
	mem         []byte
	ready       bool
	reader      *bytes.Reader
	size        string // Pre-rendered string
}

var (
	// PostgreSQL Connection pool
	pg *pgx.ConnPool

	// Cached downloads
	ramCache = [4]download{
		// These hard coded last modified timestamps are because we're working with existing files, so we provide the
		// same ones already being used
		{ // Win32
			lastRFC1123: time.Date(2017, time.September, 20, 14, 59, 44, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.10.1-win32.exe"),
				time.Date(2017, time.September, 20, 14, 59, 44, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // Win64
			lastRFC1123: time.Date(2017, time.September, 20, 14, 59, 59, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.10.1-win64.exe"),
				time.Date(2017, time.September, 20, 14, 59, 59, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // OSX
			lastRFC1123: time.Date(2017, time.September, 20, 15, 23, 27, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.10.1.dmg"),
				time.Date(2017, time.September, 20, 15, 23, 27, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // PortableApp
			lastRFC1123: time.Date(2017, time.September, 28, 19, 32, 48, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe"),
				time.Date(2017, time.September, 28, 19, 32, 48, 0, time.UTC).Format(time.RFC3339)),
		},
	}
	tmpl *template.Template
)

// Populates a cache entry
func cache(cacheEntry download) {
	cacheEntry.reader = bytes.NewReader(cacheEntry.mem)
	cacheEntry.size = fmt.Sprintf("%d", len(cacheEntry.mem))
	cacheEntry.ready = true
}

// Handler for download requests
func handler(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.URL.Path {
	case "/favicon.ico":
		http.ServeFile(w, r, filepath.Join(baseDir, "favicon.ico"))
		logDownload(r, 90022, http.StatusOK) // 90022 is the file size of favicon.ico in bytes
	case "/currentrelease":
		bytesSent, err := fmt.Fprint(w, "3.10.1\nhttps://github.com/sqlitebrowser/sqlitebrowser/releases/tag/v3.10.1\n")
		if err != nil {
			log.Printf("Error serving currentrelease: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logDownload(r, int64(bytesSent), http.StatusOK)
	case "/DB.Browser.for.SQLite-3.10.1-win64.exe":
		serveDownload(w, r, ramCache[DB4S_3_10_1_WIN64], "DB.Browser.for.SQLite-3.10.1-win64.exe")
	case "/DB.Browser.for.SQLite-3.10.1-win32.exe":
		serveDownload(w, r, ramCache[DB4S_3_10_1_WIN32], "DB.Browser.for.SQLite-3.10.1-win32.exe")
	case "/DB.Browser.for.SQLite-3.10.1.dmg":
		serveDownload(w, r, ramCache[DB4S_3_10_1_OSX], "DB.Browser.for.SQLite-3.10.1.dmg")
	case "/SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe":
		serveDownload(w, r, ramCache[DB4S_3_10_1_PORTABLE], "SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe")
	default:
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			log.Printf("Error: %s", err)
		}
		// Send the index page listing
		err = tmpl.Execute(w, nil)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			log.Printf("Error: %s", err)
		}
		logDownload(r, 771, http.StatusOK) // The index page is 771 bytes in length
	}
}

func logDownload(r *http.Request, bytesSent int64, status int) (err error) {
	// Use the new v3 pgx/pgtype structures
	ref := &pgtype.Text{
		String: r.Referer(),
		Status: pgtype.Present,
	}
	if r.Referer() == "" {
		ref.Status = pgtype.Null
	}

	// Grab the client IP address
	IPv6 := false
	var clientIP pgtype.Text
	var clientPort int
	tempIP := r.Header.Get("X-Forwarded-For")
	if tempIP == "" {
		tempIP = r.RemoteAddr
	}
	if tempIP != "" {
		// Determine if client IP address is IPv4 or IPv6, and split out the port number
		if strings.HasPrefix(tempIP, "[") {
			IPv6 = true
			s := strings.SplitN(tempIP, "]:", 2)
			if s[1] != "" {
				clientPort, _ = strconv.Atoi(s[1])
			}
			clientIP.String = strings.TrimLeft(s[0], "[")
		} else {
			// Client IP address seems to be IPv4
			s := strings.SplitN(tempIP, ":", 2)
			if s[1] != "" {
				clientPort, _ = strconv.Atoi(s[1])
			}
			clientIP.String = s[0]
		}
	} else {
		// Can't determine client IP address
		clientIP.Status = pgtype.Null
	}

fmt.Printf("ClientIP = '%v' Port = %d; IPv6 = %v\n", clientIP, clientPort, IPv6)

	dbField := "client_ipv4"
	if IPv6 {
		dbField = "client_ipv6"
	}

	dbQuery := fmt.Sprintf("INSERT INTO download_log (%s, remote_user, request_time, request_type, request, protocol, status, body_bytes_sent, http_referer, http_user_agent, client_port)"+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", dbField)
	res, err := pg.Exec(dbQuery,
		// client_ipv4 or client_ipv6
		clientIP,
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
		r.Header.Get("User-Agent"),
		// client port
		clientPort)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return fmt.Errorf("something went wrong when inserting a new download entry.  # of entries affected != 1")
	}
	return
}

func main() {
	// TODO: Investigate the logging drivers, to ensure problems get recorded somewhere more useful than the default

	// Read our configuration settings
	var err error
	var Conf TomlConfig
	if _, err = toml.DecodeFile(configFile, &Conf); err != nil {
		log.Fatal(err)
	}
	// PostgreSQL configuration info
	pgConfig := new(pgx.ConnConfig)

	// Set the PostgreSQL configuration values
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

	// Load our HTML template
	tmpl = template.Must(template.New("downloads").ParseFiles(filepath.Join(baseDir, "template.html"))).Lookup("downloads")

	// Load the files into ram from the data directory
	ramCache[DB4S_3_10_1_WIN32].mem, err = ioutil.ReadFile(filepath.Join(dataDir, "DB.Browser.for.SQLite-3.10.1-win32.exe"))
	if err == nil {
		cache(ramCache[DB4S_3_10_1_WIN32])
	}
	ramCache[DB4S_3_10_1_WIN64].mem, err = ioutil.ReadFile(filepath.Join(dataDir, "DB.Browser.for.SQLite-3.10.1-win64.exe"))
	if err == nil {
		cache(ramCache[DB4S_3_10_1_WIN64])
	}
	ramCache[DB4S_3_10_1_OSX].mem, err = ioutil.ReadFile(filepath.Join(dataDir, "DB.Browser.for.SQLite-3.10.1.dmg"))
	if err == nil {
		cache(ramCache[DB4S_3_10_1_OSX])
	}
	ramCache[DB4S_3_10_1_PORTABLE].mem, err = ioutil.ReadFile(filepath.Join(dataDir, "SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe"))
	if err == nil {
		cache(ramCache[DB4S_3_10_1_PORTABLE])
	}

	http.HandleFunc("/", handler)
	fmt.Printf("Listening on port %d...\n", listenPort)
	err = http.ListenAndServeTLS(fmt.Sprintf(":%d", listenPort), filepath.Join(certDir, "fullchain.pem"), filepath.Join(certDir, "privkey.pem"), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Close the PG connection gracefully
	pg.Close()
}

// Serves downloads from cache
func serveDownload(w http.ResponseWriter, r *http.Request, cacheEntry download, fileName string) {
	// If the file isn't cached, check if it's ready to be cached yet
	var err error
	if !cacheEntry.ready {
		cacheEntry.mem, err = ioutil.ReadFile(filepath.Join(dataDir, fileName))
		if err == nil {
			// TODO: It'd probably be a good idea to check the SHA256 of the file contents before marking the cache as valid
			// Add the download to the cache
			cacheEntry.reader = bytes.NewReader(cacheEntry.mem)
			cacheEntry.size = fmt.Sprintf("%d", len(cacheEntry.mem))
			cacheEntry.ready = true
		}
	}

	// Send the file (if cached)
	if cacheEntry.ready {
		w.Header().Set("Last-Modified", cacheEntry.lastRFC1123)
		w.Header().Set("Content-Disposition", cacheEntry.disposition)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", cacheEntry.size)
		bytesSent, err := cacheEntry.reader.WriteTo(w)
		if err != nil {
			log.Printf("Error serving %s: %v\n", fileName, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			// TODO: Catch when an error occurs (eg client aborts download), and log that too.  Probably need extra
			//       fields added to the database for holding the info.
			//logDownload(r, bytesSent, http.StatusBadRequest)
			return
		}
		_ = logDownload(r, bytesSent, http.StatusOK)
	} else {
		// Warn the user
		_, _ = fmt.Fprintf(w, "Not yet available")
		_ = logDownload(r, 17, http.StatusNotFound)
	}
}
