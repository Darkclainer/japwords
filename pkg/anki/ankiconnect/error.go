package ankiconnect

import (
	"errors"
	"fmt"
)

// invalid api key for anki-connect
const ankiErrorInvalidAPIKey = "valid api key must be provided"

// user didn't choose profile, so collections is unavailable
const ankiErrorCollectionUnavailable = "collection is not available"

// ErrInvalidAPIKey error indicating that Anki-Connect protected by api key,
// and nor or invalid api key was provided
var ErrInvalidAPIKey = errors.New("invalid API Key was provided")

// ErrCollectionUnavailable error indicating that action can not be processed,
// because anki is partially loaded (particular example is when user didn't login in profile).
var ErrCollectionUnavailable = errors.New("collection is not available")

// ServerError represent errors that we quite sure related to anki.
// Eiher it's message from field "error" from response or can be wrapper
// around ErrPermissionDenied
type ServerError struct {
	Message string
	Err     error
}

func (e *ServerError) Unwrap() error {
	return e.Err
}

func (e *ServerError) Error() string {
	msg := e.Message
	if e.Err != nil {
		msg = e.Err.Error()
	}
	return fmt.Sprintf("anki-connect responed with error: %s", msg)
}

// UnexpectedResponseError represent error that possibly ocured because we connected to some server
// but probably not to anki-connect. Either we received unxpected status code or we were unable
// to decode body.
type UnexpectedResponseError struct {
	Status int
	Err    error
}

func (e *UnexpectedResponseError) Unwrap() error {
	return e.Err
}

func (e *UnexpectedResponseError) Error() string {
	return fmt.Sprintf("server gave unexpected response (with status %d): %s", e.Status, e.Err)
}

// ConnectionError represents error occured when we tried to connect to server.
type ConnectionError struct {
	Err error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("can not connect to anki-connect server: %s", e.Err)
}

func (e *ConnectionError) Unwrap() error {
	return e.Err
}

func newServerError(message string) error {
	switch message {
	case ankiErrorInvalidAPIKey:
		return newInvalidAPIKey()
	case ankiErrorCollectionUnavailable:
		return newCollectionUnavailable()
	}
	return &ServerError{
		Message: message,
	}
}

func newInvalidAPIKey() error {
	return &ServerError{
		Err: ErrInvalidAPIKey,
	}
}

func newCollectionUnavailable() error {
	return &ServerError{
		Err: ErrCollectionUnavailable,
	}
}

func newUnexpectedStatusError(status int) error {
	return &UnexpectedResponseError{
		Status: status,
		Err:    errors.New("unexpected status code"),
	}
}

func newUnableDecodedError(status int, err error) error {
	return &UnexpectedResponseError{
		Status: status,
		Err:    fmt.Errorf("unable to decode response body: %w", err),
	}
}
