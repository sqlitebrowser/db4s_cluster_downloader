package main

import (
	"log"
	"strings"
)

var (
	// Our custom http error logger
	httpErrorLogger *log.Logger
)

// FilteringErrorLogWriter is a custom error logger for our http servers, to filter out the copious
// 'TLS handshake error' messages we're getting.
// Very heavily based on: https://github.com/golang/go/issues/26918#issuecomment-974257205
type FilteringErrorLogWriter struct{}

func (*FilteringErrorLogWriter) Write(msg []byte) (int, error) {
	z := string(msg)
	if !(strings.HasPrefix(z, "http: TLS handshake error") && strings.HasSuffix(z, ": EOF\n")) {
		err := httpErrorLogger.Output(0, z)
		if err != nil {
			log.Println(err)
		}
	}
	return len(msg), nil
}

// HttpErrorLog filters out the copious 'TLS handshake error' messages we're getting
func HttpErrorLog() *log.Logger {
	httpErrorLogger = log.New(&FilteringErrorLogWriter{}, "", log.LstdFlags)
	return httpErrorLogger
}
