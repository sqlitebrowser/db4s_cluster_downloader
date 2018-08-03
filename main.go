package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

const (
	baseDir = "/go/src/github.com/justinclift/db4s_cluster_downloader" // Location of the go program
	dataDir = "/data" // Directory where the certs and downloads are located.  Shared with the host
)

const (
	DB4S_3_10_1_WIN32 = iota
	DB4S_3_10_1_WIN64
	DB4S_3_10_1_OSX
	DB4S_3_10_1_PORTABLE
)

type download struct {
	lastRFC1123 string // Pre-rendered string
	lastRFC3339 string // Pre-rendered string
	mem         []byte
	ready       bool
	reader      *bytes.Reader
	size        string // Pre-rendered string
}

var (
	ramCache = [4]download{
		{ // Win32
			lastRFC1123: time.Date(2017, time.September, 20, 14, 59, 44, 0, time.UTC).Format(time.RFC1123),
			lastRFC3339: time.Date(2017, time.September, 20, 14, 59, 44, 0, time.UTC).Format(time.RFC3339),
		},
		{ // Win64
			lastRFC1123: time.Date(2017, time.September, 20, 14, 59, 59, 0, time.UTC).Format(time.RFC1123),
			lastRFC3339: time.Date(2017, time.September, 20, 14, 59, 59, 0, time.UTC).Format(time.RFC3339),
		},
		{ // OSX
			lastRFC1123: time.Date(2017, time.September, 20, 15, 23, 27, 0, time.UTC).Format(time.RFC1123),
			lastRFC3339: time.Date(2017, time.September, 20, 15, 23, 27, 0, time.UTC).Format(time.RFC3339),
		},
		{ // PortableApp
			lastRFC1123: time.Date(2017, time.September, 28, 19, 32, 48, 0, time.UTC).Format(time.RFC1123),
			lastRFC3339: time.Date(2017, time.September, 28, 19, 32, 48, 0, time.UTC).Format(time.RFC3339),
		},
	}
)

func handler(w http.ResponseWriter, r *http.Request) {
	// TODO: Log the downloads, so we don't lose the ability to count download numbers over time
	switch r.URL.Path {
	case "/favicon.ico":
		http.ServeFile(w, r, filepath.Join(baseDir, "favicon.ico"))
	case "/DB.Browser.for.SQLite-3.10.1-win64.exe":
		serveDownload(w, ramCache[DB4S_3_10_1_WIN64], "DB.Browser.for.SQLite-3.10.1-win64.exe")
	case "/DB.Browser.for.SQLite-3.10.1-win32.exe":
		serveDownload(w, ramCache[DB4S_3_10_1_WIN32], "DB.Browser.for.SQLite-3.10.1-win32.exe")
	case "/DB.Browser.for.SQLite-3.10.1.dmg":
		serveDownload(w, ramCache[DB4S_3_10_1_OSX], "DB.Browser.for.SQLite-3.10.1.dmg")
	case "/SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe":
		serveDownload(w, ramCache[DB4S_3_10_1_PORTABLE], "SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe")
	default:
		// TODO: Use a template instead
		fmt.Fprintf(w, "<html><head><title>DB Browser for SQLite download cluster</title></head>")
		fmt.Fprintf(w, "<body>Welcome to the DB Browser for SQLite downloads.")
		fmt.Fprintf(w, "<br /><br />")
		//fmt.Fprintf(w, "Requested path: %s", r.URL.Path)
		//fmt.Fprintf(w, "<br /><br />")
		fmt.Fprintf(w, "Available downloads:")
		fmt.Fprintf(w, "<ul>")
		fmt.Fprintf(w, "<li><a href=\"/DB.Browser.for.SQLite-3.10.1.dmg\">DB.Browser.for.SQLite-3.10.1.dmg</a> - For macOS</li>")
		fmt.Fprintf(w, "<li><a href=\"/DB.Browser.for.SQLite-3.10.1-win32.exe\">DB.Browser.for.SQLite-3.10.1-win32.exe</a> - For Windows 32-bit</li>")
		fmt.Fprintf(w, "<li><a href=\"/DB.Browser.for.SQLite-3.10.1-win64.exe\">DB.Browser.for.SQLite-3.10.1-win64.exe</a> - For Windows 64-bit</li>")
		fmt.Fprintf(w, "<li><a href=\"/SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe\">SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe</a> - PortableApp for Windows</li>")
		fmt.Fprintf(w, "</ul></body></html>")
	}
}

func main() {
	// TODO: Open log file

	// TODO: Connect to PG

	// Load the files into ram from the data directory
	var err error
	ramCache[DB4S_3_10_1_WIN32].mem, err = ioutil.ReadFile(filepath.Join(dataDir, "DB.Browser.for.SQLite-3.10.1-win32.exe"))
	if err == nil {
		ramCache[DB4S_3_10_1_WIN32].reader = bytes.NewReader(ramCache[DB4S_3_10_1_WIN32].mem)
		ramCache[DB4S_3_10_1_WIN32].size = fmt.Sprintf("%d", len(ramCache[DB4S_3_10_1_WIN32].mem))
		ramCache[DB4S_3_10_1_WIN32].ready = true
	}
	ramCache[DB4S_3_10_1_WIN64].mem, err = ioutil.ReadFile(filepath.Join(dataDir, "DB.Browser.for.SQLite-3.10.1-win64.exe"))
	if err == nil {
		ramCache[DB4S_3_10_1_WIN64].reader = bytes.NewReader(ramCache[DB4S_3_10_1_WIN64].mem)
		ramCache[DB4S_3_10_1_WIN64].size = fmt.Sprintf("%d", len(ramCache[DB4S_3_10_1_WIN64].mem))
		ramCache[DB4S_3_10_1_WIN64].ready = true
	}
	ramCache[DB4S_3_10_1_OSX].mem, err = ioutil.ReadFile(filepath.Join(dataDir, "DB.Browser.for.SQLite-3.10.1.dmg"))
	if err == nil {
		ramCache[DB4S_3_10_1_OSX].reader = bytes.NewReader(ramCache[DB4S_3_10_1_OSX].mem)
		ramCache[DB4S_3_10_1_OSX].size = fmt.Sprintf("%d", len(ramCache[DB4S_3_10_1_OSX].mem))
		ramCache[DB4S_3_10_1_OSX].ready = true
	}
	ramCache[DB4S_3_10_1_PORTABLE].mem, err = ioutil.ReadFile(filepath.Join(dataDir, "SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe"))
	if err == nil {
		ramCache[DB4S_3_10_1_PORTABLE].reader = bytes.NewReader(ramCache[DB4S_3_10_1_PORTABLE].mem)
		ramCache[DB4S_3_10_1_PORTABLE].size = fmt.Sprintf("%d", len(ramCache[DB4S_3_10_1_PORTABLE].mem))
		ramCache[DB4S_3_10_1_PORTABLE].ready = true
	}

	http.HandleFunc("/", handler)
	fmt.Println("Listening on port 443...")
	err = http.ListenAndServeTLS(":443", filepath.Join(dataDir, "cert.pem"), filepath.Join(dataDir, "key.pem"), nil)
	if err != nil {
		// TODO - Make sure problems get recorded somewhere outside the container
		log.Fatal(err)
	}

	// TODO - Close the PG connection gracefully? (is this really needed?)

	// TODO - Close and flush the log file.  Also not sure if this is really needed.
}

func serveDownload(w http.ResponseWriter, cacheEntry download, fileName string) {
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
		w.Header().Set("Content-Disposition", cacheEntry.lastRFC3339)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", cacheEntry.size)
		_, err = cacheEntry.reader.WriteTo(w)
		if err != nil {
			log.Printf("Error serving DB.Browser.for.SQLite-3.10.1-win32.exe: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Warn the user
		fmt.Fprintf(w, "Not yet available")
	}
}