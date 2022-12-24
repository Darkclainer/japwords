package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Anki_request_assertRequest(t *testing.T) {
	testCases := []struct {
		Name            string
		Action          string
		APIKey          string
		Params          any
		ExpectedRequest *fullRequest
	}{
		{
			Name:   "empty action",
			Action: "hello",
			ExpectedRequest: &fullRequest{
				Action: "hello",
			},
		},
		{
			Name:   "scalar value param",
			Action: "string",
			Params: "myparam",
			ExpectedRequest: &fullRequest{
				Action: "string",
				Params: "myparam",
			},
		},
		{
			Name:   "scalar struct param",
			Action: "string",
			Params: struct {
				Value string
			}{
				Value: "myparam",
			},
			ExpectedRequest: &fullRequest{
				Action: "string",
				Params: map[string]any{
					"Value": "myparam",
				},
			},
		},
		{
			Name:   "apikey",
			Action: "test",
			APIKey: "mykey",
			ExpectedRequest: &fullRequest{
				Action: "test",
				Key:    "mykey",
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			server := httptest.NewServer(sequentialHandler(
				handlerAssertRequest(t, tc.ExpectedRequest),
				handlerRespondJSON(t, &fullResponse{
					Result: nil,
				}),
			))
			defer server.Close()
			a, err := New(&Options{
				URL:    server.URL,
				APIKey: tc.APIKey,
			})
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			err = a.request(ctx, tc.Action, tc.Params, nil)
			assert.NoError(t, err)
		})
	}
}

func Test_Anki_request_ConnectionError(t *testing.T) {
	a, err := New(&Options{
		URL: "http://127.0.0.1:0",
	})
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = a.request(ctx, "action", nil, nil)
	var connectionError *ConnectionError
	assert.ErrorAs(t, err, &connectionError)
}

func Test_Anki_request_Errors(t *testing.T) {
	testCases := []struct {
		Name           string
		ResponseStatus int
		ResponseBody   string
		ErrorAssert    assert.ErrorAssertionFunc
	}{
		{
			Name:           "bad status/bad body",
			ResponseStatus: http.StatusBadRequest,
			ResponseBody:   `{"`,
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				var expectedError *UnexpectedResponseError
				if !assert.ErrorAs(t, err, &expectedError) {
					return false
				}
				if !assert.Equal(t, http.StatusBadRequest, expectedError.Status) {
					return false
				}
				return assert.Contains(t, expectedError.Err.Error(), "unable to decode response body")
			},
		},
		{
			Name:           "bad status/empty error",
			ResponseStatus: http.StatusBadRequest,
			ResponseBody:   `{}`,
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				var expectedError *UnexpectedResponseError
				if !assert.ErrorAs(t, err, &expectedError) {
					return false
				}
				if !assert.Equal(t, http.StatusBadRequest, expectedError.Status) {
					return false
				}
				return assert.Contains(t, expectedError.Err.Error(), "unexpected status code")
			},
		},
		{
			Name:           "bad status/server error",
			ResponseStatus: http.StatusBadRequest,
			ResponseBody:   `{"error": "myerror"}`,
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				var expectedError *ServerError
				if !assert.ErrorAs(t, err, &expectedError) {
					return false
				}
				if !assert.Nil(t, expectedError.Err) {
					return false
				}
				return assert.Equal(t, "myerror", expectedError.Message)
			},
		},
		{
			Name:           "bad status/permission denied",
			ResponseStatus: http.StatusBadRequest,
			ResponseBody:   `{"error": "` + AnkiMessagePermissionDenied + `"}`,
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				var expectedError *ServerError
				if !assert.ErrorAs(t, err, &expectedError) {
					return false
				}
				if !assert.Equal(t, "", expectedError.Message) {
					return false
				}
				return assert.ErrorIs(t, err, ErrPermissionDenied)
			},
		},
		{
			Name:           "ok status/bad body",
			ResponseStatus: http.StatusOK,
			ResponseBody:   `{"`,
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				var expectedError *UnexpectedResponseError
				if !assert.ErrorAs(t, err, &expectedError) {
					return false
				}
				if !assert.Equal(t, http.StatusOK, expectedError.Status) {
					return false
				}
				return assert.Contains(t, expectedError.Err.Error(), "unable to decode response body")
			},
		},
		{
			Name:           "ok status/server error",
			ResponseStatus: http.StatusOK,
			ResponseBody:   `{"error": "myerror"}`,
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				var expectedError *ServerError
				if !assert.ErrorAs(t, err, &expectedError) {
					return false
				}
				if !assert.Nil(t, expectedError.Err) {
					return false
				}
				return assert.Equal(t, "myerror", expectedError.Message)
			},
		},
		{ // this is how anki-connect actually responds in my version
			Name:           "ok status/permission denied",
			ResponseStatus: http.StatusOK,
			ResponseBody:   `{"error": "` + AnkiMessagePermissionDenied + `"}`,
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				var expectedError *ServerError
				if !assert.ErrorAs(t, err, &expectedError) {
					return false
				}
				if !assert.Equal(t, "", expectedError.Message) {
					return false
				}
				return assert.ErrorIs(t, err, ErrPermissionDenied)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			server := httptest.NewServer(sequentialHandler(
				handlerRespondStatusBody(t, tc.ResponseStatus, tc.ResponseBody),
			))
			defer server.Close()
			a, err := New(&Options{
				URL: server.URL,
			})
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			err = a.request(ctx, "action", nil, nil)
			tc.ErrorAssert(t, err)
		})
	}
}

func sequentialHandler(handlers ...http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			handler.ServeHTTP(w, r)
		}
	})
}

func handlerAssertRequest(t *testing.T, expected *fullRequest) http.Handler {
	return http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json; charset=utf-8", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json; charset=utf-8", r.Header.Get("Accept"))
		buffer, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		var request fullRequest
		decoder := json.NewDecoder(bytes.NewReader(buffer))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&request)
		assert.NoError(t, err)
		assert.Equal(t, apiVersion, request.Version)
		request.Version = 0
		assert.Equal(t, expected, &request)
		r.Body = io.NopCloser(bytes.NewReader(buffer))
	})
}

func handlerRespondJSON(t *testing.T, resp *fullResponse) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		buffer, err := json.Marshal(resp)
		require.NoError(t, err)
		_, err = w.Write(buffer)
		require.NoError(t, err)
	})
}

func handlerRespondStatusBody(t *testing.T, status int, body string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	})
}

func prepareMockServer(t *testing.T, handlers ...http.Handler) (context.Context, *Anki) {
	server := httptest.NewServer(sequentialHandler(handlers...))
	t.Cleanup(server.Close)
	a, err := New(&Options{
		URL: server.URL,
	})
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	t.Cleanup(cancel)
	return ctx, a
}
