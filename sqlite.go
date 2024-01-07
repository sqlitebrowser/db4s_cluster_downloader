package main

import (
	"log"

	sqlite "github.com/gwenn/gosqlite"
)

func connectSQLite(fileName string) (err error) {
	sdb, err = sqlite.Open(fileName, sqlite.OpenReadWrite|sqlite.OpenCreate|sqlite.OpenFullMutex)
	if err != nil {
		log.Printf("Couldn't open SQLite database: %s", err)
		return
	}
	err2 := sdb.EnableExtendedResultCodes(true)
	if err2 != nil {
		log.Printf("Couldn't enable extended result codes! Error: %v", err2.Error())
		return
	}

	// If it doesn't already exist, create the table to log downloads
	dbQuery := `
		CREATE TABLE IF NOT EXISTS download_log (
			download_id INTEGER PRIMARY KEY,
			remote_user text,
			request_time timestamp with time zone,
			request_type text,
			request text,
			protocol text,
			status integer,
			body_bytes_sent bigint,
			http_referer text,
			http_user_agent text,
			client_ipv4 text,
			client_ipv6 text,
			client_ip_strange text,
			client_port integer
		)`
	err = sdb.Exec(dbQuery)
	if err != nil {
		log.Printf("Something went wrong when creating the SQLite table for recording downloads: %v", err)
		return
	}
	return
}
