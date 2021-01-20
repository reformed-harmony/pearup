package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/reformed-harmony/pearup/assets"
	"github.com/reformed-harmony/pearup/auth"
	"github.com/reformed-harmony/pearup/db"
	"github.com/reformed-harmony/pearup/matcher"
	"github.com/sirupsen/logrus"

	// Load the naturaltime filter
	_ "github.com/flosch/pongo2-addons"
)

// Server provides the web interface.
type Server struct {
	cfg         *Config
	listener    net.Listener
	router      *mux.Router
	mediaDir    string
	conn        *db.Conn
	matcher     *matcher.Matcher
	store       *sessions.CookieStore
	templateSet *pongo2.TemplateSet
	log         *logrus.Entry
	stopped     chan bool
}

// New creates and initializes the server.
func New(cfg *Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	var (
		r = mux.NewRouter()
		s = &Server{
			cfg:         cfg,
			listener:    l,
			router:      mux.NewRouter(),
			mediaDir:    cfg.MediaDir,
			conn:        cfg.Conn,
			matcher:     cfg.Matcher,
			store:       sessions.NewCookieStore([]byte(cfg.SecretKey)),
			templateSet: pongo2.NewSet("", &vfsgenLoader{}),
			log:         logrus.WithField("context", "server"),
			stopped:     make(chan bool),
		}
		server = http.Server{
			Handler: r,
		}
	)

	// Create the Facebook and Google auth providers
	var (
		fProvider = auth.NewFacebook(&auth.Config{
			ClientID:       cfg.FacebookClientID,
			ClientSecret:   cfg.FacebookClientSecret,
			RedirectURL:    fmt.Sprintf("https://%s/login/facebook/callback", cfg.Host),
			SuccessHandler: http.HandlerFunc(s.loginSucceeded),
			ErrorHandler:   http.HandlerFunc(s.loginFailed),
		})
		gProvider = auth.NewGoogle(&auth.Config{
			ClientID:       cfg.GoogleClientID,
			ClientSecret:   cfg.GoogleClientSecret,
			RedirectURL:    fmt.Sprintf("https://%s/login/google/callback", cfg.Host),
			SuccessHandler: http.HandlerFunc(s.loginSucceeded),
			ErrorHandler:   http.HandlerFunc(s.loginFailed),
		})
	)

	// API routes
	s.router.HandleFunc("/api/users", s.requireLogin(s.apiUsers))

	// User routes
	s.router.HandleFunc("/profile", s.requireLogin(s.profile))
	s.router.HandleFunc("/options", s.requireLogin(s.options))
	s.router.HandleFunc("/matches", s.requireLogin(s.matches))

	// Pear-Up Routes
	s.router.HandleFunc("/pearups", s.requireLogin(s.pearups))
	s.router.HandleFunc("/pearups/{id:[0-9]+}", s.requireLogin(s.viewPearup))
	s.router.HandleFunc("/pearups/{id:[0-9]+}/register", s.requireLogin(s.registerForPearup)).Methods(http.MethodPost)

	// Admin routes
	s.router.HandleFunc("/admin", s.requireAdmin(s.adminIndex))
	s.router.HandleFunc("/admin/pearups", s.requireAdmin(s.adminPearups))
	s.router.HandleFunc("/admin/pearups/new", s.requireAdmin(s.adminEditPearup))
	s.router.HandleFunc("/admin/pearups/{id:[0-9]+}", s.requireAdmin(s.adminViewPearup))
	s.router.HandleFunc("/admin/pearups/{id:[0-9]+}/{match_id:[0-9]+}/delete", s.requireAdmin(s.adminViewPearupRemoveMatch))
	s.router.HandleFunc("/admin/pearups/{id:[0-9]+}/edit", s.requireAdmin(s.adminEditPearup))
	s.router.HandleFunc("/admin/pearups/{id:[0-9]+}/delete", s.requireAdmin(s.adminDeletePearup))
	s.router.HandleFunc("/admin/users", s.requireAdmin(s.adminUsers))
	s.router.HandleFunc("/admin/users/{id:[0-9]+}/edit", s.requireAdmin(s.adminEditUser))

	// Miscellaneous routes
	s.router.HandleFunc("/", s.index)
	s.router.HandleFunc("/login", s.login)
	s.router.HandleFunc("/logout", s.logout)
	s.router.HandleFunc("/privacy", s.privacy)

	// Login routes
	s.router.Handle("/login/facebook", fProvider.LoginHandler())
	s.router.Handle("/login/facebook/callback", fProvider.CallbackHandler())
	s.router.Handle("/login/google", gProvider.LoginHandler())
	s.router.Handle("/login/google/callback", gProvider.CallbackHandler())

	// Static and media files
	r.PathPrefix("/static").Handler(http.FileServer(assets.Assets))
	r.PathPrefix("/media").Handler(
		http.StripPrefix("/media", http.FileServer(http.Dir(cfg.MediaDir))),
	)

	// All page requests must go through the ServeHTTP method
	r.PathPrefix("/").Handler(s)

	go func() {
		defer close(s.stopped)
		defer s.log.Info("server stopped")
		s.log.Info("server started")
		if err := server.Serve(l); err != nil {
			s.log.Error(err)
		}
	}()
	return s, nil
}

// ServeHTTP loads the current user from the database if logged in.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			s.renderError(w, r, err.Error())
			return
		}
	}
	r = s.loadUser(r)
	if s.cfg.Debug {
		rec := NewRecorder(w)
		s.router.ServeHTTP(rec, r)
		s.log.Debugf(
			"%s - %s %s [%d %s] %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			rec.StatusCode,
			http.StatusText(rec.StatusCode),
			rec.Elapsed(),
		)
	} else {
		s.router.ServeHTTP(w, r)
	}
}

// Close shuts down the server.
func (s *Server) Close() {
	s.listener.Close()
	<-s.stopped
}
