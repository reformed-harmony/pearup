package server

import (
	"net/http"
	"strconv"

	"github.com/flosch/pongo2/v4"
)

func (s *Server) render(w http.ResponseWriter, r *http.Request, templateName string, ctx pongo2.Context) {
	err := func() error {
		t, err := s.templateSet.FromFile(templateName)
		if err != nil {
			return err
		}
		ctx["request"] = r
		ctx["site_theme"] = s.cfg.SiteTheme
		ctx["site_title"] = s.cfg.SiteTitle
		ctx["site_description"] = s.cfg.SiteDescription
		ctx["alerts"] = s.getAlerts(w, r)
		ctx["user"] = r.Context().Value(contextUser)
		b, err := t.ExecuteBytes(ctx)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return nil
	}()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) renderError(w http.ResponseWriter, r *http.Request, err string) {
	s.render(w, r, "error.html", pongo2.Context{
		"title": "Error",
		"error": err,
	})
}
