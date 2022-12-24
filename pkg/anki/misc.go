package anki

import (
	"context"
	"errors"
)

func (a *Anki) Version(ctx context.Context) (int, error) {
	var result int
	err := a.request(ctx, "version", nil, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

type RequestPermissionResponse struct {
	// Permission also returned as error from method
	Permission    string `json:"permission,omitempty"`
	RequireAPIKey bool   `json:"requireApiKey,omitempty"`
	Version       int    `json:"version,omitempty"`
}

var ErrRequestPermissionDenied = errors.New("request permission denied")

func (a *Anki) RequestPermission(ctx context.Context) (*RequestPermissionResponse, error) {
	var result RequestPermissionResponse
	err := a.request(ctx, "requestPermission", nil, &result)
	if err != nil {
		return nil, err
	}
	if result.Permission == "denied" {
		return &result, newSpecificServerError(ErrRequestPermissionDenied)
	}
	return &result, nil
}

func (a *Anki) LoadProfile(ctx context.Context, name string) error {
	request := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	var result bool
	err := a.request(ctx, "loadProfile", &request, &result)
	if err != nil {
		return err
	}
	if !result {
		return newServerError("profile load failed (probably specified profile doesn't exists)")
	}
	return nil
}
