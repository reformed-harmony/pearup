package server

import (
	"github.com/reformed-harmony/pearup/db"
	"github.com/reformed-harmony/pearup/matcher"
)

// Config stores the configuration for the embedded web server.
type Config struct {
	Addr                 string
	Debug                bool
	Host                 string
	FacebookClientID     string
	FacebookClientSecret string
	GoogleClientID       string
	GoogleClientSecret   string
	MediaDir             string
	SecretKey            string
	SiteTheme            string
	SiteTitle            string
	SiteDescription      string
	Conn                 *db.Conn
	Matcher              *matcher.Matcher
}
