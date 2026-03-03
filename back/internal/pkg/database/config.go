package database

const Postgres = "postgres"

type DbConfig struct {
	Driver               string
	Host                 string
	User                 string
	Password             string
	Db                   string
	Port                 string
	SslMode              bool
	MaxOpenConnections   int
	MaxIdleConnections   int
	Logging              bool
	DisableAutomaticPing bool
}
