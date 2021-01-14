//go:generate go run assets/generate.go

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/reformed-harmony/pearup/db"
	"github.com/reformed-harmony/pearup/matcher"
	"github.com/reformed-harmony/pearup/server"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "pearup"
	app.Usage = "run the pearup application"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db-host",
			Value:  "postgres",
			EnvVar: "DB_HOST",
			Usage:  "PostgreSQL database host",
		},
		cli.IntFlag{
			Name:   "db-port",
			Value:  5432,
			EnvVar: "DB_PORT",
			Usage:  "PostgreSQL database port",
		},
		cli.StringFlag{
			Name:   "db-name",
			Value:  "postgres",
			EnvVar: "DB_NAME",
			Usage:  "PostgreSQL database name",
		},
		cli.StringFlag{
			Name:   "db-user",
			Value:  "postgres",
			EnvVar: "DB_NAME",
			Usage:  "PostgreSQL database user",
		},
		cli.StringFlag{
			Name:   "db-password",
			Value:  "postgres",
			EnvVar: "DB_PASSWORD",
			Usage:  "PostgreSQL database password",
		},
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "DEBUG",
			Usage:  "enable debug logging",
		},
		cli.StringFlag{
			Name:   "fb-client-id",
			EnvVar: "FB_CLIENT_ID",
			Usage:  "Facebook client ID",
		},
		cli.StringFlag{
			Name:   "fb-client-secret",
			EnvVar: "FB_CLIENT_SECRET",
			Usage:  "Facebook client secret",
		},
		cli.StringFlag{
			Name:   "google-client-id",
			EnvVar: "GOOGLE_CLIENT_ID",
			Usage:  "Google client ID",
		},
		cli.StringFlag{
			Name:   "google-client-secret",
			EnvVar: "GOOGLE_CLIENT_SECRET",
			Usage:  "Google client secret",
		},
		cli.StringFlag{
			Name:   "media-dir",
			Value:  "media",
			EnvVar: "MEDIA_DIR",
			Usage:  "storage directory for media",
		},
		cli.StringFlag{
			Name:   "server-addr",
			Value:  ":8000",
			EnvVar: "SERVER_ADDR",
			Usage:  "server address",
		},
		cli.StringFlag{
			Name:   "server-host",
			EnvVar: "SERVER_HOST",
			Usage:  "server host",
		},
		cli.StringFlag{
			Name:   "server-secret-key",
			EnvVar: "SERVER_SECRET_KEY",
			Usage:  "secret key for sessions",
		},
		cli.StringFlag{
			Name:   "site-theme",
			Value:  "flatly",
			EnvVar: "SITE_THEME",
			Usage:  "website CSS theme",
		},
		cli.StringFlag{
			Name:   "site-title",
			Value:  "Pearup",
			EnvVar: "SITE_TITLE",
			Usage:  "website title",
		},
		cli.StringFlag{
			Name:   "site-description",
			Value:  "This website provides an easy way to meet new people through pear-ups.",
			EnvVar: "SITE_DESCRIPTION",
			Usage:  "website description",
		},
	}
	app.Action = func(c *cli.Context) error {

		// Enable debug logging if requested
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}

		// Connect to the database
		conn, err := db.New(&db.Config{
			Host:     c.String("db-host"),
			Port:     c.Int("db-port"),
			Name:     c.String("db-name"),
			User:     c.String("db-user"),
			Password: c.String("db-password"),
		})
		if err != nil {
			return err
		}
		defer conn.Close()

		// Perform all database migrations
		if err = conn.Migrate(); err != nil {
			return err
		}

		// Create the matcher
		m := matcher.New(&matcher.Config{
			Conn: conn,
		})
		defer m.Close()

		// Start the server
		s, err := server.New(&server.Config{
			Addr:                 c.String("server-addr"),
			Debug:                c.Bool("debug"),
			Host:                 c.String("server-host"),
			FacebookClientID:     c.String("fb-client-id"),
			FacebookClientSecret: c.String("fb-client-secret"),
			GoogleClientID:       c.String("google-client-id"),
			GoogleClientSecret:   c.String("google-client-secret"),
			MediaDir:             c.String("media-dir"),
			SecretKey:            c.String("server-secret-key"),
			SiteTheme:            c.String("site-theme"),
			SiteTitle:            c.String("site-title"),
			SiteDescription:      c.String("site-description"),
			Conn:                 conn,
			Matcher:              m,
		})
		if err != nil {
			return err
		}
		defer s.Close()

		// Wait for SIGINT or SIGTERM
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err.Error())
	}
}
