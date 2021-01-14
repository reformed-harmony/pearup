package server

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/nathan-osman/go-ezform"
	"github.com/nathan-osman/go-ezform/fields"
	"github.com/nathan-osman/go-ezform/validators"
	"github.com/reformed-harmony/pearup/db"
)

const (
	genderMale   = "male"
	genderFemale = "female"
)

var (
	genderMap = map[string]string{
		genderMale:   "Male",
		genderFemale: "Female",
	}

	errInvalidURL = errors.New("invalid Facebook URL (please visit fb.com/me)")
)

// fixLink works normalizes the Facebook profile URLs
func fixLink(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return link
	}
	host := strings.ToLower(u.Host)
	if host == "facebook.com" || host == "m.facebook.com" {
		u.Host = "www.facebook.com"
	}
	u.Fragment = ""
	return u.String()
}

type linkValidator struct{}

func (l *linkValidator) Validate(v string) error {
	u, err := url.Parse(fixLink(v))
	if err != nil || u.Host != "www.facebook.com" {
		return errInvalidURL
	}
	return nil
}

type editProfileForm struct {
	Link   *fields.String
	Gender *fields.String
}

func (s *Server) profile(w http.ResponseWriter, r *http.Request) {
	err := s.conn.Transaction(func(conn *db.Conn) error {
		var (
			user = r.Context().Value(contextUser).(*db.User)
			form = &editProfileForm{
				Link: fields.NewString(
					&linkValidator{},
				),
				Gender: fields.NewString(
					&validators.Choice{Choices: genderMap},
				),
			}
			ctx = pongo2.Context{
				"title":   "Profile",
				"form":    form,
				"genders": genderMap,
			}
		)
		form.Link.SetValue(user.Link)
		if user.IsMale {
			form.Gender.SetValue(genderMale)
		}
		if user.IsFemale {
			form.Gender.SetValue(genderFemale)
		}
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
				if err := s.conn.Save(user).Error; err != nil {
					return err
				}
				url := s.getRedirect(w, r)
				if len(url) == 0 {
					url = "/"
				}
				s.addAlert(w, r, flashInfo, "your profile has been updated")
				http.Redirect(w, r, url, http.StatusFound)
				return nil
			}
		}
		s.render(w, r, "profile.html", ctx)
		return nil
	})
	if err != nil {
		s.renderError(w, r, err.Error())
	}
}
