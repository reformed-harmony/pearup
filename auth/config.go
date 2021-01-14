package auth

import (
	"net/http"
)

// Config stores configuration information for a specific provider.
type Config struct {
	ClientID       string
	ClientSecret   string
	RedirectURL    string
	SuccessHandler http.Handler
	ErrorHandler   http.Handler
}
