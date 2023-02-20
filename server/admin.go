package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/flosch/pongo2/v4"
	"github.com/gorilla/mux"
	"github.com/nathan-osman/go-ezform"
	"github.com/nathan-osman/go-ezform/fields"
	"github.com/nathan-osman/go-ezform/validators"
	"github.com/reformed-harmony/pearup/db"
)

const dateFormat = "2006-01-02 15:04:05 MST"

func (s *Server) adminIndex(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "admin/index.html", pongo2.Context{
		"title": "Admin",
	})
}

func (s *Server) adminPearups(w http.ResponseWriter, r *http.Request) {
	var pearups []*db.Pearup
	if err := s.conn.
		Order("end_date DESC").
		Find(&pearups).Error; err != nil {
		s.renderError(w, r, err.Error())
		return
	}
	s.render(w, r, "admin/pearups/index.html", pongo2.Context{
		"title":   "Pear-Ups",
		"pearups": pearups,
	})
}

func (s *Server) adminViewPearup(w http.ResponseWriter, r *http.Request) {
	err := s.conn.Transaction(func(conn *db.Conn) error {
		var (
			id            = mux.Vars(r)["id"]
			p             = &db.Pearup{}
			registrations []*db.Registration
			matches       []*db.Match
		)
		if err := conn.Find(p, id).Error; err != nil {
			return err
		}
		if r.Method == http.MethodPost {
			var (
				user1ID, _ = strconv.ParseInt(r.Form.Get("User1ID"), 10, 64)
				user2ID, _ = strconv.ParseInt(r.Form.Get("User2ID"), 10, 64)
				m          = &db.Match{
					PearupID: p.ID,
					User1ID:  user1ID,
					User2ID:  user2ID,
				}
			)
			if err := conn.Save(m).Error; err != nil {
				return err
			}
			s.addAlert(w, r, flashInfo, "match was successfully added")
			http.Redirect(w, r, fmt.Sprintf("/admin/pearups/%d", p.ID), http.StatusFound)
			return nil
		}
		if err := conn.
			Joins("LEFT JOIN users ON registrations.user_id = users.id").
			Order("users.name").
			Preload("User").
			Find(&registrations, "pearup_id = ?", id).Error; err != nil {
			return err
		}
		var (
			numMen   = 0
			numWomen = 0
		)
		for _, r := range registrations {
			if r.User.IsMale {
				numMen++
			}
			if r.User.IsFemale {
				numWomen++
			}
		}
		if err := conn.
			Preload("User1").
			Preload("User2").
			Find(&matches, "pearup_id = ?", id).Error; err != nil {
			return err
		}
		s.render(w, r, "admin/pearups/view.html", pongo2.Context{
			"title":         p.Name,
			"pearup":        p,
			"registrations": registrations,
			"matches":       matches,
			"numMen":        numMen,
			"numWomen":      numWomen,
		})
		return nil
	})
	if err != nil {
		s.renderError(w, r, err.Error())
	}
}

func (s *Server) adminViewPearupRemoveMatch(w http.ResponseWriter, r *http.Request) {
	if err := s.conn.Transaction(func(conn *db.Conn) error {
		var (
			id      = mux.Vars(r)["id"]
			matchID = mux.Vars(r)["match_id"]
			m       = &db.Match{}
		)
		if err := conn.
			Preload("User1").
			Preload("User2").
			Where("pearup_id = ?", id).
			Find(m, matchID).Error; err != nil {
			return err
		}
		if r.Method == http.MethodPost {
			if err := conn.Delete(m).Error; err != nil {
				return err
			}
			s.addAlert(w, r, flashInfo, "match was successfully deleted")
			http.Redirect(w, r, fmt.Sprintf("/admin/pearups/%d", m.PearupID), http.StatusFound)
			return nil
		}
		s.render(w, r, "delete.html", pongo2.Context{
			"title": fmt.Sprintf(
				"Delete Match \"%s and %s\"",
				m.User1.Name,
				m.User2.Name,
			),
		})
		return nil
	}); err != nil {
		s.renderError(w, r, err.Error())
	}
}

type adminEditPearupForm struct {
	Name       *fields.String
	EndDate    *fields.String
	IsComplete *fields.Boolean
	IsPublic   *fields.Boolean
	CanRequest *fields.Boolean
}

func (s *Server) adminEditPearup(w http.ResponseWriter, r *http.Request) {
	err := s.conn.Transaction(func(conn *db.Conn) error {
		var (
			id   = mux.Vars(r)["id"]
			p    = &db.Pearup{}
			form = &adminEditPearupForm{
				Name:       fields.NewString(validators.NonEmpty{}),
				EndDate:    fields.NewString(validators.NonEmpty{}),
				IsComplete: fields.NewBoolean(),
				IsPublic:   fields.NewBoolean(),
				CanRequest: fields.NewBoolean(),
			}
			ctx = pongo2.Context{
				"form": form,
			}
		)
		form.EndDate.SetValue(time.Now().Format(dateFormat))
		if len(id) != 0 {
			if err := conn.Find(p, id).Error; err != nil {
				return err
			}
			form.Name.SetValue(p.Name)
			form.EndDate.SetValue(p.EndDate.Format(dateFormat))
			form.IsComplete.SetValue(p.IsComplete)
			form.IsPublic.SetValue(p.IsPublic)
			form.CanRequest.SetValue(p.CanRequest)
			ctx["title"] = fmt.Sprintf("Edit %s", p.Name)
		} else {
			ctx["title"] = "New Pear-Up"
		}
		if r.Method == http.MethodPost {
			for {
				if err := ezform.Validate(r, form); err != nil {
					if err == ezform.ErrInvalid {
						break
					}
					return err
				}
				t, err := time.Parse(dateFormat, form.EndDate.Value())
				if err != nil {
					ctx["error"] = "end date is invalid"
					break
				}
				p.Name = form.Name.Value()
				p.EndDate = t
				p.IsComplete = form.IsComplete.Value()
				p.IsPublic = form.IsPublic.Value()
				p.CanRequest = form.CanRequest.Value()
				if err := s.conn.Save(p).Error; err != nil {
					return err
				}
				s.matcher.Trigger()
				http.Redirect(w, r, "/admin/pearups", http.StatusFound)
				return nil
			}
		}
		s.render(w, r, "admin/pearups/edit.html", ctx)
		return nil
	})
	if err != nil {
		s.renderError(w, r, err.Error())
	}
}

func (s *Server) adminDeletePearup(w http.ResponseWriter, r *http.Request) {
	err := s.conn.Transaction(func(conn *db.Conn) error {
		var (
			id = mux.Vars(r)["id"]
			p  = &db.Pearup{}
		)
		if err := s.conn.Find(p, id).Error; err != nil {
			return err
		}
		if r.Method == http.MethodPost {
			if err := s.conn.Delete(p).Error; err != nil {
				return err
			}
			http.Redirect(w, r, "/admin/pearups", http.StatusFound)
			return nil
		}
		s.render(w, r, "delete.html", pongo2.Context{
			"title": fmt.Sprintf("Delete %s", p.Name),
		})
		return nil
	})
	if err != nil {
		s.renderError(w, r, err.Error())
	}
}

func (s *Server) adminUsers(w http.ResponseWriter, r *http.Request) {
	var (
		users   []*db.User
		page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	)
	if err := s.conn.
		Order("name").
		Offset(page * 30).
		Limit(31).
		Find(&users).Error; err != nil {
		s.renderError(w, r, err.Error())
		return
	}
	if page == 0 {
		page = 1
	}
	hasMore := false
	if len(users) > 30 {
		users = users[:30]
		hasMore = true
	}
	s.render(w, r, "admin/users/index.html", pongo2.Context{
		"title":   "Users",
		"users":   users,
		"page":    page,
		"hasMore": hasMore,
	})
}

type adminEditUserForm struct {
	Link    *fields.String
	Gender  *fields.String
	IsAdmin *fields.Boolean
}

func (s *Server) adminEditUser(w http.ResponseWriter, r *http.Request) {
	err := s.conn.Transaction(func(conn *db.Conn) error {
		var (
			id   = mux.Vars(r)["id"]
			user = &db.User{}
			form = &adminEditUserForm{
				Link: fields.NewString(
					&linkValidator{},
				),
				Gender: fields.NewString(
					&validators.Choice{Choices: genderMap},
				),
				IsAdmin: fields.NewBoolean(),
			}
			ctx = pongo2.Context{
				"form":    form,
				"genders": genderMap,
			}
		)
		if err := conn.Find(user, id).Error; err != nil {
			return err
		}
		form.Link.SetValue(user.Link)
		if user.IsMale {
			form.Gender.SetValue(genderMale)
		}
		if user.IsFemale {
			form.Gender.SetValue(genderFemale)
		}
		if user.IsAdmin {
			form.IsAdmin.SetValue(user.IsAdmin)
		}
		ctx["title"] = fmt.Sprintf("Edit %s", user.Name)
		if r.Method == http.MethodPost {
			for {
				if err := ezform.Validate(r, form); err != nil {
					if err == ezform.ErrInvalid {
						break
					}
					return err
				}
				user.Link = form.Link.Value()
				user.IsMale = form.Gender.Value() == genderMale
				user.IsFemale = form.Gender.Value() == genderFemale
				user.IsAdmin = form.IsAdmin.Value()
				if err := s.conn.Save(user).Error; err != nil {
					return err
				}
				http.Redirect(w, r, "/admin/users", http.StatusFound)
				return nil
			}
		}
		s.render(w, r, "admin/users/edit.html", ctx)
		return nil
	})
	if err != nil {
		s.renderError(w, r, err.Error())
	}
}
