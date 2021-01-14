package server

import (
	"net/http"
)

const (
	sessionRedirect = "redirect"
)

func (s *Server) setRedirect(w http.ResponseWriter, r *http.Request, url string) {
	session, _ := s.store.Get(r, sessionName)
	defer session.Save(r, w)
	session.Values[sessionRedirect] = url
}

func (s *Server) getRedirect(w http.ResponseWriter, r *http.Request) string {
	session, _ := s.store.Get(r, sessionName)
	defer session.Save(r, w)
	url := session.Values[sessionRedirect]
	delete(session.Values, sessionRedirect)
	if url == nil {
		return ""
	}
	return url.(string)
}
