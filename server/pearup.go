package server

import (
	"fmt"
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/gorilla/mux"
	"github.com/reformed-harmony/pearup/db"
)

func (s *Server) pearups(w http.ResponseWriter, r *http.Request) {
	var pearups []*db.Pearup
	if err := s.conn.
		Order("end_date DESC").
		Where("is_public = ?", true).
		Find(&pearups).Error; err != nil {
		s.renderError(w, r, err.Error())
		return
	}
	s.render(w, r, "pearups/index.html", pongo2.Context{
		"title":   "Pear-Ups",
		"pearups": pearups,
	})
}

func (s *Server) viewPearup(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		var (
			id           = mux.Vars(r)["id"]
			user         = r.Context().Value(contextUser).(*db.User)
			pearup       = &db.Pearup{}
			registration = &db.Registration{}
			matches      []*db.Match
		)
		if err := s.conn.
			Where("is_public = ?", true).
			Find(pearup, id).Error; err != nil {
			return err
		}
		if db := s.conn.
			Where("pearup_id = ? AND user_id = ?", pearup.ID, user.ID).
			First(registration); db.Error != nil && !db.RecordNotFound() {
			return db.Error
		}
		if db := s.conn.
			Preload("User1").
			Preload("User2").
			Where("pearup_id = ? AND (user1_id = ? OR user2_id = ?)", pearup.ID, user.ID, user.ID).
			Find(&matches); db.Error != nil && !db.RecordNotFound() {
			return db.Error
		}
		s.render(w, r, "pearups/view.html", pongo2.Context{
			"title":        pearup.Name,
			"pearup":       pearup,
			"registration": registration,
			"matches":      matches,
		})
		return nil
	}()
	if err != nil {
		s.renderError(w, r, err.Error())
	}
}

func (s *Server) registerForPearup(w http.ResponseWriter, r *http.Request) {
	err := s.conn.Transaction(func(conn *db.Conn) error {
		var (
			id           = mux.Vars(r)["id"]
			user         = r.Context().Value(contextUser).(*db.User)
			pearup       = &db.Pearup{}
			registration = &db.Registration{}
		)
		if err := s.conn.
			Where("is_public = ?", true).
			Find(pearup, id).Error; err != nil {
			return err
		}
		pearupURL := fmt.Sprintf("/pearups/%d", pearup.ID)
		if len(user.Link) == 0 || user.IsMale == user.IsFemale {
			s.alertError(w, r, pearupURL, "you must complete your profile before registering")
			return nil
		}

		// If the form contains "unregister" than we are *deleting* the registration
		if _, ok := r.Form["unregister"]; ok {
			if err := s.conn.
				Where("pearup_id = ? AND user_id = ?", pearup.ID, user.ID).
				Find(registration).Error; err != nil {
				s.alertError(w, r, pearupURL, "you have not yet registered")
				return nil
			}
			if err := s.conn.Delete(registration).Error; err != nil {
				return err
			}
			s.addAlert(w, r, flashInfo, "you have been un-registered for the pear-up")
		} else {
			if !s.conn.
				Where("pearup_id = ? AND user_id = ?", pearup.ID, user.ID).
				First(registration).RecordNotFound() {
				s.alertError(w, r, pearupURL, "you have already registered")
				return nil
			}
			registration.PearupID = pearup.ID
			registration.UserID = user.ID
			if err := s.conn.Save(registration).Error; err != nil {
				return err
			}
			s.addAlert(w, r, flashInfo, "you have registered for the pear-up")
		}

		http.Redirect(w, r, fmt.Sprintf("/pearups/%d", pearup.ID), http.StatusFound)
		return nil
	})
	if err != nil {
		s.renderError(w, r, err.Error())
	}
}
