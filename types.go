package main

import "github.com/jackc/pgx/pgtype"

// TomlConfig is the structure for holding the application configuration
type TomlConfig struct {
	Paths  PathInfo
	Pg     PGInfo
	Server ServerInfo
	TLS    TLSInfo
}
type PathInfo struct {
	BaseDir string // Location of the git source
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
	Debug   bool
	Port    int
	SSLPort int
}

type TLSInfo struct {
	CertFile string // Full path of the TLS certificate file
	KeyFile  string // Full path of the TLS private key file
}

// dbEntry is used for storing the new database entries
type dbEntry struct {
	ipv4      pgtype.Text
	ipv6      pgtype.Text
	ipstrange pgtype.Text
	port      pgtype.Int4
}

// RecordDownloads are used to determine where downloads are recorded
type RecordDownloads int

const (
	RECORD_IN_PG RecordDownloads = iota
	RECORD_IN_SQLITE
	RECORD_NOWHERE
)
