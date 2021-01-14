package server

import (
	"encoding/gob"
	"net/http"
)

const (
	flashInfo   = "info"
	flashDanger = "danger"
)

type alert struct {
	Type string
	Body string
}

func init() {
	gob.Register(&alert{})
}

func (s *Server) addAlert(w http.ResponseWriter, r *http.Request, alertType, body string) {
	session, _ := s.store.Get(r, sessionName)
	defer session.Save(r, w)
	session.AddFlash(&alert{
		Type: alertType,
		Body: body,
	})
}

func (s *Server) getAlerts(w http.ResponseWriter, r *http.Request) []interface{} {
	session, _ := s.store.Get(r, sessionName)
	defer session.Save(r, w)
	return session.Flashes()
}

func (s *Server) alertError(w http.ResponseWriter, r *http.Request, url, body string) {
	s.addAlert(w, r, flashDanger, body)
	http.Redirect(w, r, url, http.StatusFound)
}
