package anki

import (
	"errors"
	"fmt"
)

const AnkiMessagePermissionDenied = "valid api key must be provided"

var ErrPermissionDenied = errors.New("permsission denied")

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
	if message == AnkiMessagePermissionDenied {
		return newPermissionDeniedError()
	}
	return &ServerError{
		Message: message,
	}
}

func newPermissionDeniedError() error {
	return &ServerError{
		Err: ErrPermissionDenied,
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
