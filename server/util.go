package server

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/reformed-harmony/pearup/db"
)

func (s *Server) updatePicture(u *db.User, r io.ReadCloser) error {
	u.Picture = fmt.Sprintf("%d.jpg", u.ID)
	filename := path.Join(s.mediaDir, u.Picture)
	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}
	return s.conn.Save(u).Error
}
