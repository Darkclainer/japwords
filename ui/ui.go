package ui

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
)

//go:embed public
var public embed.FS

// fileSystem is little wrapper that provides fallback for index.html
type fileSystem struct {
	http.FileSystem
}

func (mfs *fileSystem) Open(name string) (http.File, error) {
	f, err := mfs.FileSystem.Open(name)
	if errors.Is(err, fs.ErrNotExist) {
		return mfs.FileSystem.Open("index.html")
	}
	return f, err
}

func Handler(path string) http.Handler {
	publicfs, err := fs.Sub(public, "public")
	if err != nil {
		panic(err)
	}

	return http.StripPrefix(
		path,
		http.FileServer(
			&fileSystem{
				FileSystem: http.FS(publicfs),
			},
		),
	)
}
