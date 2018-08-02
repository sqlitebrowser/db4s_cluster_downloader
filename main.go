package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

const (
	baseDir = "/go/src/github.com/justinclift/db4s_cluster_downloader"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// TODO - Log the downloads, so we don't lose the ability to count download numbers over time
	switch r.URL.Path {
	case "/favicon.ico":
		http.ServeFile(w, r, filepath.Join(baseDir, "favicon.ico"))
	case "/DB.Browser.for.SQLite-3.10.1-win64.exe":
		fmt.Fprintf(w, "Not yet available")
	case "/DB.Browser.for.SQLite-3.10.1-win32.exe":
		http.ServeFile(w, r, filepath.Join(baseDir, "DB.Browser.for.SQLite-3.10.1-win32.exe"))
	case "/DB.Browser.for.SQLite-3.10.1.dmg":
		fmt.Fprintf(w, "Not yet available")
	case "/SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe":
		fmt.Fprintf(w, "Not yet available")
	default:
		// TODO - Use a template instead
		fmt.Fprintf(w, "<html><head><title>DB Browser for SQLite download cluster</title></head>")
		fmt.Fprintf(w, "<body>Welcome to the DB Browser for SQLite download cluster.")
		fmt.Fprintf(w, "<br /><br />")
		fmt.Fprintf(w, "Requested path: %s", r.URL.Path)
		fmt.Fprintf(w, "<br /><br />")
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
	// TODO - Open log file

	// TODO - Connect to PG

	http.HandleFunc("/", handler)
	fmt.Println("Listening on port 443...")
	err := http.ListenAndServeTLS(":443", filepath.Join(baseDir, "cert.pem"), filepath.Join(baseDir, "key.pem"), nil)
	if err != nil {
		// TODO - Log the problem somewhere outside the container
		log.Fatal(err)
	}

	// TODO - Close the PG connection gracefully? (is this really needed?)

	// TODO - Close and flush the log file.  Also not sure if this is really needed.
}
