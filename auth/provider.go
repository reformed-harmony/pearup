package auth

import (
	"context"
	"net/http"

	"github.com/dghubble/gologin/v2"
	"golang.org/x/oauth2"
)

const (

	// ProviderKey provides access to the Provider interface upon successful auth.
	ProviderKey = "provider"

	// ErrorKey provides access to the error message when auth fails.
	ErrorKey = "error"
)

// Provider implements authentication for a specific provider.
type Provider interface {

	// Name returns a unique name for the provider.
	Name() string

	// LoginHandler returns an HTTP handler for beginning auth.
	LoginHandler() http.Handler

	// CallbackHandler returns an HTTP handler for completing auth.
	CallbackHandler() http.Handler

	// User returns information about the user.
	User(context.Context) (*User, error)
}

type providerData struct {
	config         *oauth2.Config
	successHandler http.Handler
	errorHandler   http.Handler
}

func (p *providerData) init(provider Provider, cfg *Config, endpoint oauth2.Endpoint, scopes []string) {
	p.config = &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     endpoint,
		Scopes:       scopes,
	}
	p.successHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), ProviderKey, provider))
		cfg.SuccessHandler.ServeHTTP(w, r)
	})
	p.errorHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), ErrorKey, gologin.ErrorFromContext(r.Context())))
		cfg.ErrorHandler.ServeHTTP(w, r)
	})
}
