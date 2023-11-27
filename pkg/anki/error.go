package anki

import (
	"errors"
	"fmt"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
)

var (
	// redefine this errors from client, because we probably want to expose them as API
	ErrForbiddenOrigin       = errors.New("anki-connect forbid request from client origin")
	ErrInvalidAPIKey         = errors.New("anki-connect rejected request because api key is invalid")
	ErrCollectionUnavailable = errors.New("anki-connect is not ready for specified action")

	ErrNoteTypeNotExists     = errors.New("selected note type doesn't exists")
	ErrDeckAlreadyExists     = errors.New("deck with the same name already exists")
	ErrNoteTypeAlreadyExists = errors.New("note type with the same name already exists")

	// ErrUnknownServerError unrecognized error from anki-connect, but probably should
	ErrUnknownServerError = errors.New("anki-connect returned unknown error")
)

type ConnectionError struct {
	Msg string
}

func (err *ConnectionError) Error() string {
	return fmt.Sprintf("unable to connect to anki-connect: %s", err.Msg)
}

// convertAnkiError tries to convert error to known error, and returns true if error considered to be "permanent"
func convertAnkiError(err error) (error, bool) {
	var serverError *ankiconnect.ServerError
	if errors.As(err, &serverError) {
		switch serverError.Err {
		case ankiconnect.ErrCollectionUnavailable:
			return ErrCollectionUnavailable, true
		case ankiconnect.ErrInvalidAPIKey:
			return ErrInvalidAPIKey, true
		default:
			// this case should not happen, but Anki-Connect do not provide
			// documentation for all errors, so it's signal to that additional
			// errors should be checked
			// It's not permanent, because it is expected to come from actions like
			// create something.
			return fmt.Errorf("%w: %s", ErrUnknownServerError, err.Error()), false
		}
	}
	var connectionError *ankiconnect.ConnectionError
	if errors.As(err, &connectionError) {
		return &ConnectionError{
			Msg: connectionError.Error(),
		}, true
	}
	var unexpectedResponseError *ankiconnect.UnexpectedResponseError
	if errors.As(err, &unexpectedResponseError) {
		return &ConnectionError{
			Msg: fmt.Sprintf("got unexpected response from server: %s", unexpectedResponseError.Error()),
		}, true
	}
	// this is not anki-client related error
	return err, false
}
