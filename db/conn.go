package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	// Enable runtime support for Postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Conn maintains a connection to the database.
type Conn struct {
	*gorm.DB
	log *logrus.Entry
}

// New creates a new connection to the database and initializes it.
func New(cfg *Config) (*Conn, error) {
	g, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			cfg.Host,
			cfg.Port,
			cfg.Name,
			cfg.User,
			cfg.Password,
		),
	)
	if err != nil {
		return nil, err
	}
	c := &Conn{
		DB:  g,
		log: logrus.WithField("context", "db"),
	}
	return c, nil
}

// Migrate performs all database migrations.
func (c *Conn) Migrate() error {
	c.log.Info("performing migrations...")
	return c.AutoMigrate(
		&User{},
		&Pearup{},
		&Registration{},
		&Match{},
		&Request{},
		&Exclusion{},
	).Error
}

// Transaction wraps the provided callback function in a database transaction.
func (c *Conn) Transaction(fn func(*Conn) error) error {
	var (
		tx   = c.Begin()
		conn = &Conn{DB: c.DB}
	)
	conn.DB = tx
	if err := fn(conn); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
