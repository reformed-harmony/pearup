package db

// Config provides the configuration for the database.
type Config struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}
