package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
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
	"github.com/mitchellh/go-homedir"
)

const (
	DB4S_3_10_1_WIN32 = iota // The order needs to match the ramCache entries in the global var section
	DB4S_3_10_1_WIN64
	DB4S_3_10_1_OSX
	DB4S_3_10_1_PORTABLE
	DB4S_3_11_0_WIN32_MSI
	DB4S_3_11_0_WIN32_ZIP
	DB4S_3_11_0_WIN64_MSI
	DB4S_3_11_0_WIN64_ZIP
	DB4S_3_11_0_OSX
	DB4S_3_11_1_WIN32_MSI
	DB4S_3_11_1_WIN32_ZIP
	DB4S_3_11_1_WIN64_MSI
	DB4S_3_11_1_WIN64_ZIP
	DB4S_3_11_1_OSX
	DB4S_3_11_1V2_OSX
	DB4S_3_11_2_WIN32_MSI
	DB4S_3_11_2_WIN32_ZIP
	DB4S_3_11_2_WIN64_MSI
	DB4S_3_11_2_WIN64_ZIP
	DB4S_3_11_2_OSX
	DB4S_3_11_2_PORTABLE
	DB4S_3_11_2_PORTABLE_V2
	DB4S_3_12_0_WIN32_MSI
	DB4S_3_12_0_WIN32_ZIP
	DB4S_3_12_0_WIN64_MSI
	DB4S_3_12_0_WIN64_ZIP
	DB4S_3_12_0_OSX
	DB4S_3_12_0_PORTABLE
	DB4S_3_12_2_WIN32_MSI
	DB4S_3_12_2_WIN32_ZIP
	DB4S_3_12_2_WIN64_MSI
	DB4S_3_12_2_WIN64_ZIP
	DB4S_3_12_2_OSX
	DB4S_3_12_2_PORTABLE
	DB4S_3_12_2_APPIMAGE
	DB4S_3_12_2_OSX_ARM64
)

// Configuration file
type TomlConfig struct {
	Paths  PathInfo
	Pg     PGInfo
	Server ServerInfo
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
	Debug bool
	Port  int
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

	// Should debugging info be displayed?
	debug = false

	// PostgreSQL Connection pool
	pg *pgx.ConnPool

	// Cached downloads
	ramCache = [36]cacheEntry{
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
		{ // DB4S 3.11.0 Win32 MSI
			lastRFC1123: time.Date(2019, time.February, 5, 17, 33, 47, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.0-win32.msi"),
				time.Date(2019, time.February, 5, 17, 33, 47, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.0 Win32 zip
			lastRFC1123: time.Date(2019, time.February, 5, 17, 34, 1, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.0-win32.zip"),
				time.Date(2019, time.February, 5, 17, 34, 1, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.0 Win64 MSI
			lastRFC1123: time.Date(2019, time.February, 5, 17, 34, 21, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.0-win64.msi"),
				time.Date(2019, time.February, 5, 17, 34, 21, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.0 Win64 zip
			lastRFC1123: time.Date(2019, time.February, 5, 17, 34, 44, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.0-win64.zip"),
				time.Date(2019, time.February, 5, 17, 34, 44, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.0 OSX
			lastRFC1123: time.Date(2019, time.February, 7, 9, 50, 18, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.0.dmg"),
				time.Date(2019, time.February, 7, 9, 50, 18, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.1 Win32 MSI
			lastRFC1123: time.Date(2019, time.February, 18, 16, 28, 5, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.1-win32.msi"),
				time.Date(2019, time.February, 18, 16, 28, 5, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.1 Win32 zip
			lastRFC1123: time.Date(2019, time.February, 18, 16, 28, 16, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.1-win32.zip"),
				time.Date(2019, time.February, 18, 16, 28, 16, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.1 Win64 MSI
			lastRFC1123: time.Date(2019, time.February, 18, 16, 28, 35, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.1-win64.msi"),
				time.Date(2019, time.February, 18, 16, 28, 35, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.1 Win64 zip
			lastRFC1123: time.Date(2019, time.February, 18, 16, 28, 50, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.1-win64.zip"),
				time.Date(2019, time.February, 18, 16, 28, 50, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.1 OSX
			lastRFC1123: time.Date(2019, time.February, 18, 10, 37, 48, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.1.dmg"),
				time.Date(2019, time.February, 18, 10, 37, 48, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.1v2 OSX
			lastRFC1123: time.Date(2019, time.February, 23, 9, 15, 10, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.1v2.dmg"),
				time.Date(2019, time.February, 23, 9, 15, 10, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.2 Win32 MSI
			lastRFC1123: time.Date(2019, time.April, 3, 18, 13, 2, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.2-win32.msi"),
				time.Date(2019, time.April, 3, 18, 13, 2, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.2 Win32 zip
			lastRFC1123: time.Date(2019, time.April, 3, 18, 13, 16, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.2-win32.zip"),
				time.Date(2019, time.April, 3, 18, 13, 16, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.2 Win64 MSI
			lastRFC1123: time.Date(2019, time.April, 3, 18, 13, 35, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.2-win64.msi"),
				time.Date(2019, time.April, 3, 18, 13, 35, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.2 Win64 zip
			lastRFC1123: time.Date(2019, time.April, 3, 18, 14, 8, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.2-win64.zip"),
				time.Date(2019, time.April, 3, 18, 14, 8, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.2 OSX
			lastRFC1123: time.Date(2019, time.April, 3, 14, 48, 13, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.11.2.dmg"),
				time.Date(2019, time.April, 3, 14, 48, 13, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.2 Portable
			lastRFC1123: time.Date(2019, time.May, 7, 10, 48, 35, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("SQLiteDatabaseBrowserPortable_3.11.2_English.paf.exe"),
				time.Date(2019, time.May, 7, 10, 48, 35, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.11.2 Portable v2
			lastRFC1123: time.Date(2019, time.May, 14, 22, 59, 52, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("SQLiteDatabaseBrowserPortable_3.11.2_Rev_2_English.paf.exe"),
				time.Date(2019, time.May, 14, 22, 59, 52, 0, time.UTC).Format(time.RFC3339)),
		},

		// *** 3.12.0 release ***
		{ // DB4S 3.12.0 Win32 MSI
			lastRFC1123: time.Date(2020, time.June, 15, 18, 18, 1, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.12.0-win32.msi"),
				time.Date(2020, time.June, 15, 18, 18, 1, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.12.0 Win32 zip
			lastRFC1123: time.Date(2020, time.June, 15, 18, 18, 9, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.12.0-win32.zip"),
				time.Date(2020, time.June, 15, 18, 18, 9, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.12.0 Win64 MSI
			lastRFC1123: time.Date(2020, time.June, 15, 18, 18, 19, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.12.0-win64.msi"),
				time.Date(2020, time.June, 15, 18, 18, 19, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.12.0 Win64 zip
			lastRFC1123: time.Date(2020, time.June, 15, 18, 18, 37, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.12.0-win64.zip"),
				time.Date(2020, time.June, 15, 18, 18, 37, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.12.0 OSX
			lastRFC1123: time.Date(2020, time.June, 14, 7, 24, 20, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.12.0.dmg"),
				time.Date(2020, time.June, 14, 7, 24, 20, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.12.0 Portable
			lastRFC1123: time.Date(2020, time.June, 18, 4, 59, 35, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("SQLiteDatabaseBrowserPortable_3.12.0_English.paf.exe"),
				time.Date(2020, time.June, 18, 4, 59, 35, 0, time.UTC).Format(time.RFC3339)),
		},

		// *** 3.12.2 release ***
		{ // DB4S 3.12.2 Win32 MSI
			lastRFC1123: time.Date(2021, time.May, 17, 12, 39, 2, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.12.2-win32.msi"),
				time.Date(2021, time.May, 17, 12, 39, 2, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.12.2 Win32 zip
			lastRFC1123: time.Date(2021, time.May, 16, 20, 0, 6, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.12.2-win32.zip"),
				time.Date(2021, time.May, 16, 20, 0, 6, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.12.2 Win64 MSI
			lastRFC1123: time.Date(2021, time.May, 17, 12, 39, 16, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.12.2-win64.msi"),
				time.Date(2021, time.May, 17, 12, 39, 16, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.12.2 Win64 zip
			lastRFC1123: time.Date(2021, time.May, 16, 20, 0, 21, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.12.2-win64.zip"),
				time.Date(2021, time.May, 16, 20, 0, 21, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.12.2 OSX Intel
			lastRFC1123: time.Date(2021, time.May, 9, 11, 14, 6, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-3.12.2.dmg"),
				time.Date(2021, time.May, 9, 11, 14, 6, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.12.2 Portable
			lastRFC1123: time.Date(2021, time.May, 19, 16, 42, 57, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("SQLiteDatabaseBrowserPortable_3.12.2_English.paf.exe"),
				time.Date(2021, time.May, 19, 16, 42, 57, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.12.2 AppImage
			lastRFC1123: time.Date(2021, time.July, 7, 6, 55, 29, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB_Browser_for_SQLite-v3.12.2-x86_64.AppImage"),
				time.Date(2021, time.July, 7, 6, 55, 29, 0, time.UTC).Format(time.RFC3339)),
		},
		{ // DB4S 3.12.2 OSX ARM64
			lastRFC1123: time.Date(2022, time.October, 23, 16, 16, 06, 0, time.UTC).Format(time.RFC1123),
			disposition: fmt.Sprintf(`attachment; filename="%s"; modification-date="%s";`,
				url.QueryEscape("DB.Browser.for.SQLite-arm64-3.12.2.dmg"),
				time.Date(2022, time.October, 23, 16, 16, 06, 0, time.UTC).Format(time.RFC3339)),
		},
	}
	tmpl *template.Template

	// Is the PostgreSQL connection working?
	usePG = true
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

	// Enable debugging output, if the option is set in the config file
	debug = Conf.Server.Debug

	// * Connect to PG database *

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
	pg, err = pgx.NewConnPool(pgPoolConfig)
	if err != nil {
		log.Printf("Couldn't connect to PostgreSQL: '%s'. Continuing, but downloads won't be recorded.", err)
		usePG = false
	} else {
		// Log successful connection
		log.Printf("Connected to PostgreSQL server: %v:%v\n", Conf.Pg.Server, uint16(Conf.Pg.Port))
	}

	// Load our HTML template
	// TODO: Embed the template in the compiled binary
	tmpl = template.Must(template.New("downloads").ParseFiles(filepath.Join(Conf.Paths.BaseDir, "template.html"))).Lookup("downloads")

	// Load the files into ram from the data directory
	ramCache[DB4S_3_10_1_WIN32].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.10.1-win32.exe"))
	if err == nil {
		cache(ramCache[DB4S_3_10_1_WIN32])
	}
	ramCache[DB4S_3_10_1_WIN64].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.10.1-win64.exe"))
	if err == nil {
		cache(ramCache[DB4S_3_10_1_WIN64])
	}
	ramCache[DB4S_3_10_1_OSX].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.10.1.dmg"))
	if err == nil {
		cache(ramCache[DB4S_3_10_1_OSX])
	}
	ramCache[DB4S_3_10_1_PORTABLE].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe"))
	if err == nil {
		cache(ramCache[DB4S_3_10_1_PORTABLE])
	}
	ramCache[DB4S_3_11_0_WIN32_MSI].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.0-win32.msi"))
	if err == nil {
		cache(ramCache[DB4S_3_11_0_WIN32_MSI])
	}
	ramCache[DB4S_3_11_0_WIN32_ZIP].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.0-win32.zip"))
	if err == nil {
		cache(ramCache[DB4S_3_11_0_WIN32_ZIP])
	}
	ramCache[DB4S_3_11_0_WIN64_MSI].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.0-win64.msi"))
	if err == nil {
		cache(ramCache[DB4S_3_11_0_WIN64_MSI])
	}
	ramCache[DB4S_3_11_0_WIN64_ZIP].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.0-win64.zip"))
	if err == nil {
		cache(ramCache[DB4S_3_11_0_WIN64_ZIP])
	}
	ramCache[DB4S_3_11_0_OSX].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.0.dmg"))
	if err == nil {
		cache(ramCache[DB4S_3_11_0_OSX])
	}
	ramCache[DB4S_3_11_1_WIN32_MSI].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.1-win32.msi"))
	if err == nil {
		cache(ramCache[DB4S_3_11_1_WIN32_MSI])
	}
	ramCache[DB4S_3_11_1_WIN32_ZIP].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.1-win32.zip"))
	if err == nil {
		cache(ramCache[DB4S_3_11_1_WIN32_ZIP])
	}
	ramCache[DB4S_3_11_1_WIN64_MSI].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.1-win64.msi"))
	if err == nil {
		cache(ramCache[DB4S_3_11_1_WIN64_MSI])
	}
	ramCache[DB4S_3_11_1_WIN64_ZIP].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.1-win64.zip"))
	if err == nil {
		cache(ramCache[DB4S_3_11_1_WIN64_ZIP])
	}
	ramCache[DB4S_3_11_1_OSX].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.1.dmg"))
	if err == nil {
		cache(ramCache[DB4S_3_11_1_OSX])
	}
	ramCache[DB4S_3_11_1V2_OSX].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.1v2.dmg"))
	if err == nil {
		cache(ramCache[DB4S_3_11_1V2_OSX])
	}
	ramCache[DB4S_3_11_2_WIN32_MSI].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.2-win32.msi"))
	if err == nil {
		cache(ramCache[DB4S_3_11_2_WIN32_MSI])
	}
	ramCache[DB4S_3_11_2_WIN32_ZIP].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.2-win32.zip"))
	if err == nil {
		cache(ramCache[DB4S_3_11_2_WIN32_ZIP])
	}
	ramCache[DB4S_3_11_2_WIN64_MSI].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.2-win64.msi"))
	if err == nil {
		cache(ramCache[DB4S_3_11_2_WIN64_MSI])
	}
	ramCache[DB4S_3_11_2_WIN64_ZIP].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.2-win64.zip"))
	if err == nil {
		cache(ramCache[DB4S_3_11_2_WIN64_ZIP])
	}
	ramCache[DB4S_3_11_2_OSX].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.11.2.dmg"))
	if err == nil {
		cache(ramCache[DB4S_3_11_2_OSX])
	}
	ramCache[DB4S_3_11_2_PORTABLE].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "SQLiteDatabaseBrowserPortable_3.11.2_English.paf.exe"))
	if err == nil {
		cache(ramCache[DB4S_3_11_2_PORTABLE])
	}
	ramCache[DB4S_3_11_2_PORTABLE_V2].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "SQLiteDatabaseBrowserPortable_3.11.2_Rev_2_English.paf.exe"))
	if err == nil {
		cache(ramCache[DB4S_3_11_2_PORTABLE_V2])
	}

	// 3.12.0 release
	ramCache[DB4S_3_12_0_WIN32_MSI].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.12.0-win32.msi"))
	if err == nil {
		cache(ramCache[DB4S_3_12_0_WIN32_MSI])
	}
	ramCache[DB4S_3_12_0_WIN32_ZIP].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.12.0-win32.zip"))
	if err == nil {
		cache(ramCache[DB4S_3_12_0_WIN32_ZIP])
	}
	ramCache[DB4S_3_12_0_WIN64_MSI].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.12.0-win64.msi"))
	if err == nil {
		cache(ramCache[DB4S_3_12_0_WIN64_MSI])
	}
	ramCache[DB4S_3_12_0_WIN64_ZIP].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.12.0-win64.zip"))
	if err == nil {
		cache(ramCache[DB4S_3_12_0_WIN64_ZIP])
	}
	ramCache[DB4S_3_12_0_OSX].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.12.0.dmg"))
	if err == nil {
		cache(ramCache[DB4S_3_12_0_OSX])
	}
	ramCache[DB4S_3_12_0_PORTABLE].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "SQLiteDatabaseBrowserPortable_3.12.0_English.paf.exe"))
	if err == nil {
		cache(ramCache[DB4S_3_12_0_PORTABLE])
	}

	// 3.12.2 release
	ramCache[DB4S_3_12_2_WIN32_MSI].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.12.2-win32.msi"))
	if err == nil {
		cache(ramCache[DB4S_3_12_2_WIN32_MSI])
	}
	ramCache[DB4S_3_12_2_WIN32_ZIP].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.12.2-win32.zip"))
	if err == nil {
		cache(ramCache[DB4S_3_12_2_WIN32_ZIP])
	}
	ramCache[DB4S_3_12_2_WIN64_MSI].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.12.2-win64.msi"))
	if err == nil {
		cache(ramCache[DB4S_3_12_2_WIN64_MSI])
	}
	ramCache[DB4S_3_12_2_WIN64_ZIP].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.12.2-win64.zip"))
	if err == nil {
		cache(ramCache[DB4S_3_12_2_WIN64_ZIP])
	}
	ramCache[DB4S_3_12_2_OSX].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-3.12.2.dmg"))
	if err == nil {
		cache(ramCache[DB4S_3_12_2_OSX])
	}
	ramCache[DB4S_3_12_2_PORTABLE].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "SQLiteDatabaseBrowserPortable_3.12.2_English.paf.exe"))
	if err == nil {
		cache(ramCache[DB4S_3_12_2_PORTABLE])
	}
	ramCache[DB4S_3_12_2_APPIMAGE].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB_Browser_for_SQLite-v3.12.2-x86_64.AppImage"))
	if err == nil {
		cache(ramCache[DB4S_3_12_2_APPIMAGE])
	}
	ramCache[DB4S_3_12_2_OSX_ARM64].mem, err = ioutil.ReadFile(filepath.Join(Conf.Paths.DataDir, "DB.Browser.for.SQLite-arm64-3.12.2.dmg"))
	if err == nil {
		cache(ramCache[DB4S_3_12_2_OSX_ARM64])
	}

	http.HandleFunc("/", handler)
	fmt.Printf("Listening on port %d...\n", Conf.Server.Port)
	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", Conf.Server.Port),
		ErrorLog: HttpErrorLog(),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12, // TLS 1.2 is now the lowest acceptable level
		},
	}
	err = srv.ListenAndServeTLS(filepath.Join(Conf.Paths.CertDir, "fullchain.pem"), filepath.Join(Conf.Paths.CertDir, "privkey.pem"))
	if err != nil {
		log.Fatal(err)
	}

	// Close the PG connection gracefully
	if usePG {
		pg.Close()
	}
}

// Populates a cache entry
func cache(cacheEntry cacheEntry) {
	cacheEntry.reader = bytes.NewReader(cacheEntry.mem)
	cacheEntry.size = fmt.Sprintf("%d", len(cacheEntry.mem))
	cacheEntry.ready = true
}

// Handler for download requests
func handler(w http.ResponseWriter, r *http.Request) {
	// Set the maximum accepted http request size, for safety
	r.Body = http.MaxBytesReader(w, r.Body, 4096) // 4k seems like a reasonable max size

	var err error
	switch r.URL.Path {
	case "/favicon.ico":
		http.ServeFile(w, r, filepath.Join(Conf.Paths.BaseDir, "favicon.ico"))
		err = logRequest(r, 90022, http.StatusOK) // 90022 is the file size of favicon.ico in bytes
		if err != nil {
			log.Printf("Error: %s", err)
		}
		if debug {
			log.Printf("Successful favicon.ico request, Client: %s\n", r.RemoteAddr)
		}
	case "/currentrelease":
		bytesSent, err := fmt.Fprint(w, "3.12.2\nhttps://sqlitebrowser.org/blog/version-3-12-2-released\n")
		if err != nil {
			log.Printf("Error serving currentrelease: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = logRequest(r, int64(bytesSent), http.StatusOK)
		if err != nil {
			log.Printf("Error: %s", err)
		}
		if debug {
			log.Printf("Successful /currentrelease request, Client: %s\n", r.RemoteAddr)
		}
	case "/DB.Browser.for.SQLite-3.10.1-win64.exe":
		serveDownload(w, r, ramCache[DB4S_3_10_1_WIN64], "DB.Browser.for.SQLite-3.10.1-win64.exe")
	case "/DB.Browser.for.SQLite-3.10.1-win32.exe":
		serveDownload(w, r, ramCache[DB4S_3_10_1_WIN32], "DB.Browser.for.SQLite-3.10.1-win32.exe")
	case "/DB.Browser.for.SQLite-3.10.1.dmg":
		serveDownload(w, r, ramCache[DB4S_3_10_1_OSX], "DB.Browser.for.SQLite-3.10.1.dmg")
	case "/SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe":
		serveDownload(w, r, ramCache[DB4S_3_10_1_PORTABLE], "SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe")
	case "/DB.Browser.for.SQLite-3.11.0-win32.msi":
		serveDownload(w, r, ramCache[DB4S_3_11_0_WIN32_MSI], "DB.Browser.for.SQLite-3.11.0-win32.msi")
	case "/DB.Browser.for.SQLite-3.11.0-win32.zip":
		serveDownload(w, r, ramCache[DB4S_3_11_0_WIN32_ZIP], "DB.Browser.for.SQLite-3.11.0-win32.zip")
	case "/DB.Browser.for.SQLite-3.11.0-win64.msi":
		serveDownload(w, r, ramCache[DB4S_3_11_0_WIN64_MSI], "DB.Browser.for.SQLite-3.11.0-win64.msi")
	case "/DB.Browser.for.SQLite-3.11.0-win64.zip":
		serveDownload(w, r, ramCache[DB4S_3_11_0_WIN64_ZIP], "DB.Browser.for.SQLite-3.11.0-win64.zip")
	case "/DB.Browser.for.SQLite-3.11.0.dmg":
		serveDownload(w, r, ramCache[DB4S_3_11_0_OSX], "DB.Browser.for.SQLite-3.11.0.dmg")
	case "/DB.Browser.for.SQLite-3.11.1-win32.msi":
		serveDownload(w, r, ramCache[DB4S_3_11_1_WIN32_MSI], "DB.Browser.for.SQLite-3.11.1-win32.msi")
	case "/DB.Browser.for.SQLite-3.11.1-win32.zip":
		serveDownload(w, r, ramCache[DB4S_3_11_1_WIN32_ZIP], "DB.Browser.for.SQLite-3.11.1-win32.zip")
	case "/DB.Browser.for.SQLite-3.11.1-win64.msi":
		serveDownload(w, r, ramCache[DB4S_3_11_1_WIN64_MSI], "DB.Browser.for.SQLite-3.11.1-win64.msi")
	case "/DB.Browser.for.SQLite-3.11.1-win64.zip":
		serveDownload(w, r, ramCache[DB4S_3_11_1_WIN64_ZIP], "DB.Browser.for.SQLite-3.11.1-win64.zip")
	case "/DB.Browser.for.SQLite-3.11.1.dmg":
		serveDownload(w, r, ramCache[DB4S_3_11_1_OSX], "DB.Browser.for.SQLite-3.11.1.dmg")
	case "/DB.Browser.for.SQLite-3.11.1v2.dmg":
		serveDownload(w, r, ramCache[DB4S_3_11_1V2_OSX], "DB.Browser.for.SQLite-3.11.1v2.dmg")
	case "/DB.Browser.for.SQLite-3.11.2-win32.msi":
		serveDownload(w, r, ramCache[DB4S_3_11_2_WIN32_MSI], "DB.Browser.for.SQLite-3.11.2-win32.msi")
	case "/DB.Browser.for.SQLite-3.11.2-win32.zip":
		serveDownload(w, r, ramCache[DB4S_3_11_2_WIN32_ZIP], "DB.Browser.for.SQLite-3.11.2-win32.zip")
	case "/DB.Browser.for.SQLite-3.11.2-win64.msi":
		serveDownload(w, r, ramCache[DB4S_3_11_2_WIN64_MSI], "DB.Browser.for.SQLite-3.11.2-win64.msi")
	case "/DB.Browser.for.SQLite-3.11.2-win64.zip":
		serveDownload(w, r, ramCache[DB4S_3_11_2_WIN64_ZIP], "DB.Browser.for.SQLite-3.11.2-win64.zip")
	case "/DB.Browser.for.SQLite-3.11.2.dmg":
		serveDownload(w, r, ramCache[DB4S_3_11_2_OSX], "DB.Browser.for.SQLite-3.11.2.dmg")
	case "/SQLiteDatabaseBrowserPortable_3.11.2_English.paf.exe":
		serveDownload(w, r, ramCache[DB4S_3_11_2_PORTABLE], "SQLiteDatabaseBrowserPortable_3.11.2_English.paf.exe")
	case "/SQLiteDatabaseBrowserPortable_3.11.2_Rev_2_English.paf.exe":
		serveDownload(w, r, ramCache[DB4S_3_11_2_PORTABLE_V2], "SQLiteDatabaseBrowserPortable_3.11.2_Rev_2_English.paf.exe")
	case "/DB.Browser.for.SQLite-3.12.0-win32.msi":
		serveDownload(w, r, ramCache[DB4S_3_12_0_WIN32_MSI], "DB.Browser.for.SQLite-3.12.0-win32.msi")
	case "/DB.Browser.for.SQLite-3.12.0-win32.zip":
		serveDownload(w, r, ramCache[DB4S_3_12_0_WIN32_ZIP], "DB.Browser.for.SQLite-3.12.0-win32.zip")
	case "/DB.Browser.for.SQLite-3.12.0-win64.msi":
		serveDownload(w, r, ramCache[DB4S_3_12_0_WIN64_MSI], "DB.Browser.for.SQLite-3.12.0-win64.msi")
	case "/DB.Browser.for.SQLite-3.12.0-win64.zip":
		serveDownload(w, r, ramCache[DB4S_3_12_0_WIN64_ZIP], "DB.Browser.for.SQLite-3.12.0-win64.zip")
	case "/DB.Browser.for.SQLite-3.12.0.dmg":
		serveDownload(w, r, ramCache[DB4S_3_12_0_OSX], "DB.Browser.for.SQLite-3.12.0.dmg")
	case "/SQLiteDatabaseBrowserPortable_3.12.0_English.paf.exe":
		serveDownload(w, r, ramCache[DB4S_3_12_0_PORTABLE], "SQLiteDatabaseBrowserPortable_3.12.0_English.paf.exe")

	// 3.12.2 release
	case "/DB.Browser.for.SQLite-3.12.2-win32.msi":
		serveDownload(w, r, ramCache[DB4S_3_12_2_WIN32_MSI], "DB.Browser.for.SQLite-3.12.2-win32.msi")
	case "/DB.Browser.for.SQLite-3.12.2-win32.zip":
		serveDownload(w, r, ramCache[DB4S_3_12_2_WIN32_ZIP], "DB.Browser.for.SQLite-3.12.2-win32.zip")
	case "/DB.Browser.for.SQLite-3.12.2-win64.msi":
		serveDownload(w, r, ramCache[DB4S_3_12_2_WIN64_MSI], "DB.Browser.for.SQLite-3.12.2-win64.msi")
	case "/DB.Browser.for.SQLite-3.12.2-win64.zip":
		serveDownload(w, r, ramCache[DB4S_3_12_2_WIN64_ZIP], "DB.Browser.for.SQLite-3.12.2-win64.zip")
	case "/DB.Browser.for.SQLite-3.12.2.dmg":
		serveDownload(w, r, ramCache[DB4S_3_12_2_OSX], "DB.Browser.for.SQLite-3.12.2.dmg")
	case "/SQLiteDatabaseBrowserPortable_3.12.2_English.paf.exe":
		serveDownload(w, r, ramCache[DB4S_3_12_2_PORTABLE], "SQLiteDatabaseBrowserPortable_3.12.2_English.paf.exe")
	case "/DB_Browser_for_SQLite-v3.12.2-x86_64.AppImage":
		serveDownload(w, r, ramCache[DB4S_3_12_2_APPIMAGE], "DB_Browser_for_SQLite-v3.12.2-x86_64.AppImage")
	case "/DB.Browser.for.SQLite-arm64-3.12.2.dmg":
		serveDownload(w, r, ramCache[DB4S_3_12_2_OSX_ARM64], "DB.Browser.for.SQLite-arm64-3.12.2.dmg")
	default:

		// Send the index page listing
		err = tmpl.Execute(w, nil)
		if err != nil {
			_, e := fmt.Fprintf(w, "Error: %v", err)
			log.Printf("Error: %s", e)
			log.Printf("Error: %s", err)
		}
		err = logRequest(r, 4634, http.StatusOK) // The index page is 4634 bytes in length
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
}

func logRequest(r *http.Request, bytesSent int64, status int) (err error) {
	// Only log the download if the PostgreSQL connection is present
	if usePG {
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
	}
	return
}

// Serves downloads from cache
func serveDownload(w http.ResponseWriter, r *http.Request, download cacheEntry, fileName string) {
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
			err = logRequest(r, bytesSent, http.StatusBadRequest)
			if err != nil {
				log.Printf("Error: %s", err)
			}
			return
		}
		err = logRequest(r, bytesSent, http.StatusOK)
		if err != nil {
			log.Printf("Error: %s", err)
		}
		if debug {
			log.Printf("Successful download: %s, Client: %s\n", fileName, r.RemoteAddr)
		}
	} else {
		// Warn the user
		_, err = fmt.Fprintf(w, "Not yet available")
		if err != nil {
			log.Printf("Error: %s", err)
		}
		err = logRequest(r, 17, http.StatusNotFound)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
}
