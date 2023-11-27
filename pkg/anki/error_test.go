package anki

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
)

func Test_convertAnkiError(t *testing.T) {
	testCases := []struct {
		Name        string
		Err         error
		IsPermanent bool
		AssertError assert.ErrorAssertionFunc
	}{
		{
			Name: "ErrCollectionUnavailable",
			Err: &ankiconnect.ServerError{
				Err: ankiconnect.ErrCollectionUnavailable,
			},
			IsPermanent: true,
			AssertError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, ErrCollectionUnavailable)
			},
		},
		{
			Name: "ErrInvalidAPIKey",
			Err: &ankiconnect.ServerError{
				Err: ankiconnect.ErrInvalidAPIKey,
			},
			IsPermanent: true,
			AssertError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, ErrInvalidAPIKey)
			},
		},
		{
			Name: "ErrUnknownServerError",
			Err: &ankiconnect.ServerError{
				Err: errors.New("uknown error"),
			},
			IsPermanent: false,
			AssertError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, ErrUnknownServerError) && assert.ErrorContains(tt, err, "unknown error")
			},
		},
		{
			Name: "anki UnexpectedResponse",
			Err: &ankiconnect.UnexpectedResponseError{
				Err: errors.New("some unexpected error"),
			},
			IsPermanent: true,
			AssertError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				var connErr *ConnectionError
				return assert.ErrorAs(tt, err, &connErr)
			},
		},
		{
			Name: "anki ConnectionError",
			Err: &ankiconnect.ConnectionError{
				Err: errors.New("some connection error"),
			},
			IsPermanent: true,
			AssertError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				var connErr *ConnectionError
				return assert.ErrorAs(tt, err, &connErr)
			},
		},
		{
			Name:        "not anki error",
			Err:         errors.New("not anki"),
			IsPermanent: false,
			AssertError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "not anki")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			actualErr, actualIsAnkiError := convertAnkiError(tc.Err)
			assert.Equal(t, tc.IsPermanent, actualIsAnkiError)
			tc.AssertError(t, actualErr)
		})
	}
}
