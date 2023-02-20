package server

import (
	"net/http"
	"strconv"

	"github.com/flosch/pongo2/v4"
	"github.com/reformed-harmony/pearup/db"
)

const (
	optionsActionAddRequest      = "add_request"
	optionsActionAddExclusion    = "add_exclusion"
	optionsActionDeleteRequest   = "delete_request"
	optionsActionDeleteExclusion = "delete_exclusion"
)

func (s *Server) options(w http.ResponseWriter, r *http.Request) {
	if err := func() error {
		var (
			user       = r.Context().Value(contextUser).(*db.User)
			requests   []*db.Request
			exclusions []*db.Exclusion
			ctx        = pongo2.Context{
				"title": "Options",
			}
		)
		if r.Method == http.MethodPost {
			var (
				userID, _    = strconv.ParseInt(r.Form.Get("UserID"), 10, 64)
				err          error
				flashMessage string
			)
			switch r.Form.Get("Action") {
			case optionsActionAddRequest:
				err = s.conn.Save(&db.Request{
					UserID:          user.ID,
					RequestedUserID: userID,
				}).Error
				flashMessage = "added request"
			case optionsActionAddExclusion:
				err = s.conn.Save(&db.Exclusion{
					UserID:         user.ID,
					ExcludedUserID: userID,
				}).Error
				flashMessage = "added exclusion"
			case optionsActionDeleteRequest:
				err = s.conn.Delete(&db.Request{
					UserID:          user.ID,
					RequestedUserID: userID,
				}).Error
				flashMessage = "deleted request"
			case optionsActionDeleteExclusion:
				err = s.conn.Delete(&db.Exclusion{
					UserID:         user.ID,
					ExcludedUserID: userID,
				}).Error
				flashMessage = "deleted exclusion"
			}
			if err != nil {
				ctx["error"] = err
			} else {
				s.addAlert(w, r, flashInfo, flashMessage)
				http.Redirect(w, r, "/options", http.StatusFound)
				return nil
			}
		}
		if err := s.conn.
			Preload("RequestedUser").
			Where("user_id = ?", user.ID).
			Find(&requests).Error; err != nil {
			return err
		}
		ctx["requests"] = requests
		if err := s.conn.
			Preload("ExcludedUser").
			Where("user_id = ?", user.ID).
			Find(&exclusions).Error; err != nil {
			return err
		}
		ctx["exclusions"] = exclusions
		s.render(w, r, "options.html", ctx)
		return nil
	}(); err != nil {
		s.renderError(w, r, err.Error())
	}
}
