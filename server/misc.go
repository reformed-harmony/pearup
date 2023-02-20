package server

import (
	"net/http"

	"github.com/flosch/pongo2/v4"
)

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "index.html", pongo2.Context{})
}

func (s *Server) privacy(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "privacy.html", pongo2.Context{
		"title": "Privacy Policy",
	})
}
