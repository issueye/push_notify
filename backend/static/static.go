package static

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var staticEmbed embed.FS

// GetStaticFS returns the static assets filesystem
func GetStaticFS() fs.FS {
	f, err := fs.Sub(staticEmbed, "dist")
	if err != nil {
		panic(err)
	}
	return f
}
