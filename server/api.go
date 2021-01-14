package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/reformed-harmony/pearup/db"
)

func writeJSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func writeJSONError(w http.ResponseWriter, err error) {
	writeJSON(w, struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	})
}

func (s *Server) apiUsers(w http.ResponseWriter, r *http.Request) {
	var users []*db.User
	if err := s.conn.
		Where("POSITION(LOWER(?) in LOWER(name)) > 0", r.URL.Query().Get("q")).
		Limit(8).
		Find(&users).Error; err != nil {
		writeJSONError(w, err)
		return
	}
	var result []interface{}
	for _, u := range users {
		result = append(result, map[string]interface{}{
			"value":   u.ID,
			"text":    u.Name,
			"picture": u.Picture,
		})
	}
	writeJSON(w, result)
}
