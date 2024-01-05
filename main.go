package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	sqlite "github.com/gwenn/gosqlite"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/mitchellh/go-homedir"
)

var (
	// Conf holds the application configuration values
	Conf TomlConfig

	// Should debugging info be displayed?
	debug bool

	// PostgreSQL Connection pool
	pg *pgx.ConnPool

	// SQLite connection, used as fallback if PostgreSQL isn't available
	sdb *sqlite.Conn

	// Timestamps for the files.  We use hard coded values that match GitHub
	timeStamps = map[string]time.Time{
		// *** 3.10.1 release ***
		"DB.Browser.for.SQLite-3.10.1-win32.exe":               time.Date(2017, time.September, 20, 14, 59, 44, 0, time.UTC),
		"DB.Browser.for.SQLite-3.10.1-win64.exe":               time.Date(2017, time.September, 20, 14, 59, 59, 0, time.UTC),
		"DB.Browser.for.SQLite-3.10.1.dmg":                     time.Date(2017, time.September, 20, 15, 23, 27, 0, time.UTC),
		"SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe": time.Date(2017, time.September, 28, 19, 32, 48, 0, time.UTC),

		// *** 3.11.0 release ***
		"DB.Browser.for.SQLite-3.11.0-win32.msi": time.Date(2019, time.February, 5, 17, 33, 47, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.0-win32.zip": time.Date(2019, time.February, 5, 17, 34, 1, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.0-win64.msi": time.Date(2019, time.February, 5, 17, 34, 21, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.0-win64.zip": time.Date(2019, time.February, 5, 17, 34, 44, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.0.dmg":       time.Date(2019, time.February, 7, 9, 50, 18, 0, time.UTC),

		// *** 3.11.1 release ***
		"DB.Browser.for.SQLite-3.11.1-win32.msi": time.Date(2019, time.February, 18, 16, 28, 5, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.1-win32.zip": time.Date(2019, time.February, 18, 16, 28, 16, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.1-win64.msi": time.Date(2019, time.February, 18, 16, 28, 35, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.1-win64.zip": time.Date(2019, time.February, 18, 16, 28, 50, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.1.dmg":       time.Date(2019, time.February, 18, 10, 37, 48, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.1v2.dmg":     time.Date(2019, time.February, 23, 9, 15, 10, 0, time.UTC),

		// *** 3.11.2 release ***
		"DB.Browser.for.SQLite-3.11.2-win32.msi":                     time.Date(2019, time.April, 3, 18, 13, 2, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.2-win32.zip":                     time.Date(2019, time.April, 3, 18, 13, 16, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.2-win64.msi":                     time.Date(2019, time.April, 3, 18, 13, 35, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.2-win64.zip":                     time.Date(2019, time.April, 3, 18, 14, 8, 0, time.UTC),
		"DB.Browser.for.SQLite-3.11.2.dmg":                           time.Date(2019, time.April, 3, 14, 48, 13, 0, time.UTC),
		"SQLiteDatabaseBrowserPortable_3.11.2_English.paf.exe":       time.Date(2019, time.May, 7, 10, 48, 35, 0, time.UTC),
		"SQLiteDatabaseBrowserPortable_3.11.2_Rev_2_English.paf.exe": time.Date(2019, time.May, 14, 22, 59, 52, 0, time.UTC),

		// *** 3.12.0 release ***
		"DB.Browser.for.SQLite-3.12.0-win32.msi":               time.Date(2020, time.June, 15, 18, 18, 1, 0, time.UTC),
		"DB.Browser.for.SQLite-3.12.0-win32.zip":               time.Date(2020, time.June, 15, 18, 18, 9, 0, time.UTC),
		"DB.Browser.for.SQLite-3.12.0-win64.msi":               time.Date(2020, time.June, 15, 18, 18, 19, 0, time.UTC),
		"DB.Browser.for.SQLite-3.12.0-win64.zip":               time.Date(2020, time.June, 15, 18, 18, 37, 0, time.UTC),
		"DB.Browser.for.SQLite-3.12.0.dmg":                     time.Date(2020, time.June, 14, 7, 24, 20, 0, time.UTC),
		"SQLiteDatabaseBrowserPortable_3.12.0_English.paf.exe": time.Date(2020, time.June, 18, 4, 59, 35, 0, time.UTC),

		// *** 3.12.2 release ***
		"DB.Browser.for.SQLite-3.12.2-win32.msi":               time.Date(2021, time.May, 17, 12, 39, 2, 0, time.UTC),
		"DB.Browser.for.SQLite-3.12.2-win32.zip":               time.Date(2021, time.May, 16, 20, 0, 6, 0, time.UTC),
		"DB.Browser.for.SQLite-3.12.2-win64.msi":               time.Date(2021, time.May, 17, 12, 39, 16, 0, time.UTC),
		"DB.Browser.for.SQLite-3.12.2-win64.zip":               time.Date(2021, time.May, 16, 20, 0, 21, 0, time.UTC),
		"DB.Browser.for.SQLite-3.12.2.dmg":                     time.Date(2021, time.May, 9, 11, 14, 6, 0, time.UTC),
		"SQLiteDatabaseBrowserPortable_3.12.2_English.paf.exe": time.Date(2021, time.May, 19, 16, 42, 57, 0, time.UTC),
		"DB_Browser_for_SQLite-v3.12.2-x86_64.AppImage":        time.Date(2021, time.July, 7, 6, 55, 29, 0, time.UTC),
		"DB.Browser.for.SQLite-arm64-3.12.2.dmg":               time.Date(2022, time.October, 23, 16, 16, 06, 0, time.UTC),
	}

	// RecordDownloadsLocation controls where downloads are recorded
	RecordDownloadsLocation = RECORD_NOWHERE
)

func main() {
	// Read the config file
	err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database for recording downloads
	connectDatabase()

	// Set up Gin
	router, err := setupRouter(false)
	if err != nil {
		log.Fatal(err)
	}

	// Create the basic HTTP server configuration
	s := &http.Server{
		ErrorLog:     HttpErrorLog(),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// If TLS Cert and key file paths are given, then we're using TLS
	if Conf.TLS.CertFile != "" && Conf.TLS.KeyFile != "" {
		s.Addr = fmt.Sprintf(":%d", Conf.Server.SSLPort)
		s.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12, // TLS 1.2 is now the lowest acceptable level
		}
		log.Printf("Listening on port %d...", Conf.Server.SSLPort)

		// Start the server
		err = s.ListenAndServeTLS(filepath.Join(Conf.TLS.CertFile), filepath.Join(Conf.TLS.KeyFile))
	} else {
		// Not using TLS (eg for testing on GitHub Actions)
		s.Addr = fmt.Sprintf(":%d", Conf.Server.Port)
		log.Printf("Listening on port %d...", Conf.Server.Port)

		// Start the server
		err = s.ListenAndServe()
	}
	if err != nil {
		log.Fatal(err)
	}

	// Close the database connection gracefully
	if RecordDownloadsLocation == RECORD_IN_PG {
		pg.Close()
	} else if RecordDownloadsLocation == RECORD_IN_SQLITE {
		sdb.Close()
	}
}

// connectDatabase attempts to connect to the backend PostgreSQL database.  If that fails, it connects to a local
// SQLite database instead.  If *that* fails as well, it just doesn't bother recording downloads.
func connectDatabase() {
	// Setup the PostgreSQL config
	pgConfig := new(pgx.ConnConfig)
	pgConfig.Host = Conf.Pg.Server
	pgConfig.Port = uint16(Conf.Pg.Port)
	pgConfig.User = Conf.Pg.Username
	pgConfig.Password = Conf.Pg.Password
	pgConfig.Database = Conf.Pg.Database
	clientTLSConfig := tls.Config{InsecureSkipVerify: true}
	if Conf.Pg.SSL {
		pgConfig.TLSConfig = &clientTLSConfig
	} else {
		pgConfig.TLSConfig = nil
	}

	// Connect to PG
	pgPoolConfig := pgx.ConnPoolConfig{*pgConfig, Conf.Pg.NumConnections, nil, 5 * time.Second}
	var err error
	pg, err = pgx.NewConnPool(pgPoolConfig)
	if err != nil {
		fileName := "DB4S_downloads.sqlite"
		err = connectSQLite(fileName)
		if err != nil {
			// Something went wrong with the SQLite database, so we just turn off recording downloads
			RecordDownloadsLocation = RECORD_NOWHERE
		} else {
			log.Printf("Connecting to PostgreSQL failed, so recording downloads to local SQLite file '%s' instead", fileName)
			RecordDownloadsLocation = RECORD_IN_SQLITE
		}
	} else {
		// Log successful connection
		log.Printf("Recording downloads to PostgreSQL server: %v:%v", Conf.Pg.Server, uint16(Conf.Pg.Port))
		RecordDownloadsLocation = RECORD_IN_PG
	}
	return
}

// currentReleaseHandler serves the "current release" information to users
func currentReleaseHandler(c *gin.Context) {
	resp := "3.12.2\nhttps://sqlitebrowser.org/blog/version-3-12-2-released\n"
	c.String(200, resp)
}

// Handler for download requests
func fileHandler(c *gin.Context) {
	// If the requested file is unknown, then abort
	fileName := c.Param("filename")
	ts, ok := timeStamps[fileName]
	if !ok {
		fmt.Fprintf(c.Writer, "Unknown file requested")
		log.Printf("Unknown file '%s' requested by '%s', aborting", fileName, c.Request.RemoteAddr)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Retrieve the file size
	fullPath := filepath.Join(Conf.Paths.DataDir, fileName)
	info, err := os.Stat(fullPath)
	if err != nil {
		fmt.Fprintf(c.Writer, "Internal server error")
		log.Printf("Error occured when trying to stat local file '%s': %s", fileName, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	sz := strconv.FormatInt(info.Size(), 10)

	// Create the format disposition string
	disp := fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`, fileName, ts.Format(time.RFC3339))

	// Set the headers
	c.Header("Content-Disposition", disp)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", sz)

	// Send the file contents
	// We use http.ServeContent() here as it allows setting the desired "last modified" timestamp.  The other functions
	// we could have used instead - c.File() and http.ServeFile() - don't allow this.  Those ones just read the date of
	// the file on disk, whereas we want to use timestamp entries matching the GitHub release files
	z, err := os.Open(fullPath)
	if err != nil {
		fmt.Fprintf(c.Writer, "Internal server error")
		log.Printf("Error occured when trying to open local file '%s': %s", fileName, err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer z.Close()
	http.ServeContent(c.Writer, c.Request, fileName, ts, z)
}

// logRequest records a download in the backend database
func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Execute the other middleware handlers first
		c.Next()

		fileName := c.Request.URL.String()

		if debug {
			log.Printf("Logging download of '%s' (%d bytes) by '%s'", fileName, c.Writer.Size(), c.ClientIP())
		}

		// If we're recording downloads, then figure out the details
		if RecordDownloadsLocation != RECORD_NOWHERE {
			// Use the new v3 pgx/pgtype structures
			ref := &pgtype.Text{
				String: c.Request.Referer(),
				Status: pgtype.Present,
			}
			if c.Request.Referer() == "" {
				ref.Status = pgtype.Null
			}

			// Grab the client IP address
			clientIP := dbEntry{
				ipv4:      pgtype.Text{Status: pgtype.Null},
				ipv6:      pgtype.Text{Status: pgtype.Null},
				ipstrange: pgtype.Text{Status: pgtype.Null},
				port:      pgtype.Int4{Status: pgtype.Null},
			}
			tempIP := c.Request.Header.Get("X-Forwarded-For")
			if tempIP == "" {
				tempIP = c.Request.RemoteAddr
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
									log.Printf("String conversion failed! s[1] = %v, p = %v, int32(p) = %v, port = %v",
										s[1], p, int32(p), clientIP.port.Int)
									clientIP.port.Status = pgtype.Null
									clientIP.port.Int = 0
								}
							} else {
								log.Printf("Conversion error: %v", e)
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
								log.Printf("Conversion error: %v", e)
								break
							}

							// Double check the port number conversion was correct
							tst := fmt.Sprintf("%d", p)
							if tst != s[1] {
								log.Printf("String conversion failed! s[1] = %v, p = %v, int32(p) = %v, port = %v",
									s[1], p, int32(p), clientIP.port.Int)
								break
							}

							// Ensure the port number is in the valid port range (0-65535)
							if p < 0 || p > 65535 {
								log.Printf("Port number %v outside valid port range", p)
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

			if RecordDownloadsLocation == RECORD_IN_PG {
				// Record the download to PostgreSQL
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
					c.Request.Method,
					// request
					fileName,
					// protocol
					c.Request.Proto,
					// status
					c.Writer.Status(),
					// body_bytes_sent
					c.Writer.Size(),
					// http_referer
					ref,
					// http_user_agent
					c.Request.Header.Get("User-Agent"))
				if err != nil {
					log.Printf("error when inserting download entry in PostgreSQL: %v", err)
					return
				}
				numRows := res.RowsAffected()
				if numRows != 1 {
					log.Printf("something went wrong when inserting a new download entry.  # of entries affected = %d instead of 1", numRows)
					return
				}
			} else {
				// Record the download in SQLite
				// Note there's no need to convert the PG data types before hand, as the SQLite library seems ok with them
				dbQuery := `
					INSERT INTO download_log (
						client_ipv4, client_ipv6, client_ip_strange, client_port, remote_user, request_time, request_type, request,
						protocol, status, body_bytes_sent, http_referer, http_user_agent)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
				err := sdb.Exec(dbQuery,
					// IP address
					&clientIP.ipv4, &clientIP.ipv6, &clientIP.ipstrange,
					// client port
					&clientIP.port,
					// remote_user
					&pgtype.Text{String: "", Status: pgtype.Null}, // Hard coded empty string for now
					// request_time
					time.Now().Format(time.RFC3339Nano),
					// request_type
					c.Request.Method,
					// request
					fileName,
					// protocol
					c.Request.Proto,
					// status
					c.Writer.Status(),
					// body_bytes_sent
					c.Writer.Size(),
					// http_referer
					ref,
					// http_user_agent
					c.Request.Header.Get("User-Agent"))
				if err != nil {
					log.Printf("error when inserting download entry in SQLite: %v", err)
					return
				}
			}
		}
		return
	}
}

func readConfig() (err error) {
	// Override config file location via environment variables
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		// TODO: Might be a good idea to add permission checks of the dir & conf file, to ensure they're not
		//       world readable.  Similar in concept to what ssh does for its config files.
		userHome, err := homedir.Dir()
		if err != nil {
			log.Fatalf("User home directory couldn't be determined: %s", err)
		}
		configFile = filepath.Join(userHome, ".db4s", "downloader_config.toml")
	}

	// Read our configuration settings
	if _, err = toml.DecodeFile(configFile, &Conf); err != nil {
		log.Fatal(err)
	}

	// Apply any environment variable configuration overrides
	z := os.Getenv("DEBUG")
	if z != "" {
		Conf.Server.Debug, err = strconv.ParseBool(z)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Enable debugging output, if the option is set in the config file
	debug = Conf.Server.Debug

	return
}

// maxSizeMiddleware limits the maximum request size, to help prevent DOS attacks
func maxSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

// rootHandler serves the html index page that lists the available downloads
func rootHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "downloads", nil)
}

func setupRouter(testingMode bool) (router *gin.Engine, err error) {
	// We turn off Gins' debug mode when testing, and when debug mode is turned off in normal operation
	if testingMode || !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Set up Gin
	router = gin.New()
	router.Use(gin.Recovery())
	if !testingMode {
		// We don't use the Gin Logger middleware when running go test
		router.Use(gin.Logger())
	}

	// Limit the maximum size (in bytes) of incoming requests
	router.Use(maxSizeMiddleware(8192)) // 8k seems like a reasonable max size

	// Add gzip middleware
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Add CORS middleware
	// The default configuration allows all origins
	router.Use(cors.Default())

	// Log requests to PostgreSQL
	router.Use(logRequest())

	// Load our HTML template
	router.LoadHTMLGlob(filepath.Join(Conf.Paths.BaseDir, "template.html"))

	// Register handlers
	router.GET("/", rootHandler)
	router.GET("/:filename", fileHandler)
	router.GET("/currentrelease", currentReleaseHandler)
	router.StaticFile("/favicon.ico", filepath.Join(Conf.Paths.BaseDir, "favicon.ico"))
	return
}
