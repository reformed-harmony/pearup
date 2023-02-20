package assets

import (
	"embed"
	"io/fs"
)

var (
	//go:embed static
	Static embed.FS

	//go:embed templates/*
	templates    embed.FS
	Templates, _ = fs.Sub(templates, "templates")
)
