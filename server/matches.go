package server

import (
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/reformed-harmony/pearup/db"
)

func (s *Server) matches(w http.ResponseWriter, r *http.Request) {
	var (
		user    = r.Context().Value(contextUser).(*db.User)
		matches []*db.Match
	)
	if err := s.conn.
		Preload("Pearup").
		Preload("User1").
		Preload("User2").
		Joins("LEFT JOIN pearups ON matches.pearup_id = pearups.id").
		Order("pearups.end_date DESC").
		Where("matches.user1_id = ? OR matches.user2_id = ?", user.ID, user.ID).
		Find(&matches).Error; err != nil {
		s.renderError(w, r, err.Error())
		return
	}
	s.render(w, r, "matches.html", pongo2.Context{
		"title":   "Matches",
		"matches": matches,
	})
}
