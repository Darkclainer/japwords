//go:build !embedui

package ui

import (
	"net/http"
)

func Handler(path string) http.Handler {
	return http.NotFoundHandler()
}
