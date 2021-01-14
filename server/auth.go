package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/reformed-harmony/pearup/auth"
	"github.com/reformed-harmony/pearup/db"
)

type key int

const (
	contextUser key = iota

	sessionName   = "session"
	sessionUserID = "userID"
)

func (s *Server) loadUser(r *http.Request) *http.Request {
	session, _ := s.store.Get(r, sessionName)
	if v, _ := session.Values[sessionUserID]; v != nil {
		u := &db.User{}
		if err := s.conn.First(u, v).Error; err == nil {
			r = r.WithContext(context.WithValue(r.Context(), contextUser, u))
		}
	}
	return r
}

func (s *Server) requireLogin(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(contextUser) == nil {
			s.setRedirect(w, r, r.URL.Path)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		fn(w, r)
	}
}

func (s *Server) requireAdmin(fn http.HandlerFunc) http.HandlerFunc {
	return s.requireLogin(func(w http.ResponseWriter, r *http.Request) {
		if !r.Context().Value(contextUser).(*db.User).IsAdmin {
			s.renderError(w, r, "only an admin may access this page")
			return
		}
		fn(w, r)
	})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "login.html", pongo2.Context{
		"title": "Login",
	})
}

func (s *Server) findOrCreateUser(ctx context.Context) (*db.User, bool, error) {
	provider := ctx.Value(auth.ProviderKey).(auth.Provider)
	a, err := provider.User(ctx)
	if err != nil {
		return nil, false, err
	}
	defer a.Picture.Close()
	var (
		u       = &db.User{}
		created = false
	)
	switch provider.Name() {
	case "facebook":
		u.FacebookID = a.ID
	case "google":
		u.GoogleID = a.ID
	}
	if err := s.conn.Transaction(func(conn *db.Conn) error {
		db := conn.Where(u).First(u)
		if db.Error != nil {
			if !db.RecordNotFound() {
				return db.Error
			}
			created = true
		}
		u.Name = a.Name
		u.Email = a.Email
		return conn.Save(u).Error
	}); err != nil {
		return nil, false, err
	}
	if err := s.updatePicture(u, a.Picture); err != nil {
		return nil, false, err
	}
	return u, created, nil
}

func (s *Server) loginSucceeded(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		u, created, err := s.findOrCreateUser(r.Context())
		if err != nil {
			return err
		}
		session, _ := s.store.Get(r, sessionName)
		session.Values[sessionUserID] = u.ID
		session.Save(r, w)
		var (
			url      string
			greeting string
		)
		if created {
			url = "/profile"
			greeting = fmt.Sprintf("welcome, %s", u.Name)
		} else {
			url = s.getRedirect(w, r)
			if len(url) == 0 {
				url = "/"
			}
			greeting = fmt.Sprintf("welcome back, %s", u.Name)
		}
		s.addAlert(w, r, flashInfo, greeting)
		http.Redirect(w, r, url, http.StatusFound)
		return nil
	}()
	if err != nil {
		s.renderError(w, r, err.Error())
	}
}

func (s *Server) loginFailed(w http.ResponseWriter, r *http.Request) {
	s.renderError(w, r, r.Context().Value(auth.ErrorKey).(error).Error())
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, sessionName)
	session.Values[sessionUserID] = nil
	session.Save(r, w)
	s.addAlert(w, r, flashInfo, "you have successfully been logged out")
	http.Redirect(w, r, "/", http.StatusFound)
}
