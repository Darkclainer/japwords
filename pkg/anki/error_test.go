package anki

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ServerError_Unwrap(t *testing.T) {
	t.Run("not nil", func(t *testing.T) {
		wrappedErr := errors.New("some error")
		var serverErr error = &ServerError{
			Err: wrappedErr,
		}
		assert.ErrorIs(t, serverErr, wrappedErr)
	})
	t.Run("nil", func(t *testing.T) {
		var serverErr error = &ServerError{
			Err: nil,
		}
		assert.Nil(t, errors.Unwrap(serverErr))
	})
}

func Test_ServerError_Error(t *testing.T) {
	const standardMessage = "anki-connect responed with error: "
	testCases := []struct {
		Name     string
		Err      *ServerError
		Expected string
	}{
		{
			Name:     "empty",
			Err:      &ServerError{},
			Expected: standardMessage,
		},
		{
			Name: "msg only",
			Err: &ServerError{
				Message: "my error",
			},
			Expected: standardMessage + "my error",
		},
		{
			Name: "error only",
			Err: &ServerError{
				Err: errors.New("wrapped error"),
			},
			Expected: standardMessage + "wrapped error",
		},
		{
			Name: "msg and error",
			Err: &ServerError{
				Message: "my error",
				Err:     errors.New("wrapped error"),
			},
			Expected: standardMessage + "wrapped error",
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			var err error = tc.Err
			assert.Equal(t, tc.Expected, err.Error())
		})
	}
}

func Test_UnexpectedResponseError_Unwrap(t *testing.T) {
	t.Run("not nil", func(t *testing.T) {
		wrappedErr := errors.New("some error")
		var unexpectedResponseErr error = &UnexpectedResponseError{
			Err: wrappedErr,
		}
		assert.ErrorIs(t, unexpectedResponseErr, wrappedErr)
	})
	t.Run("nil", func(t *testing.T) {
		var unexpectedResponseErr error = &UnexpectedResponseError{
			Err: nil,
		}
		assert.Nil(t, errors.Unwrap(unexpectedResponseErr))
	})
}

func Test_ConnectionError_Unwrap(t *testing.T) {
	t.Run("not nil", func(t *testing.T) {
		wrappedErr := errors.New("some error")
		var connectionErr error = &ConnectionError{
			Err: wrappedErr,
		}
		assert.ErrorIs(t, connectionErr, wrappedErr)
	})
	t.Run("nil", func(t *testing.T) {
		var connectionErr error = &ConnectionError{
			Err: nil,
		}
		assert.Nil(t, errors.Unwrap(connectionErr))
	})
}
