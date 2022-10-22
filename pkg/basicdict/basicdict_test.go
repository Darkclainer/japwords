package basicdict

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BasicDict_Query(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	t.Run("check body and request", func(t *testing.T) {
		testFetcher := NewTestFetcher(func(req *http.Request) (*http.Response, error) {
			buffer := io.NopCloser(bytes.NewReader([]byte("hello world")))
			return &http.Response{
				StatusCode: 200,
				Body:       buffer,
				Request:    req,
			}, nil
		})
		bd := New(testFetcher)
		html, err := bd.Query(ctx, "hello")
		require.NoError(t, err)
		assert.Equal(t, []byte("hello world"), html)
	})
	t.Run("do error", func(t *testing.T) {
		testFetcher := NewTestFetcher(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("myerror")
		})
		bd := New(testFetcher)
		_, err := bd.Query(ctx, "hello")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "myerror")
		assert.Contains(t, err.Error(), "request failed")
	})
	t.Run("status not ok", func(t *testing.T) {
		testFetcher := NewTestFetcher(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewReader(nil)),
				Request:    req,
			}, nil
		})
		bd := New(testFetcher)
		_, err := bd.Query(ctx, "hello")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "response status")
	})
	t.Run("body error", func(t *testing.T) {
		testFetcher := NewTestFetcher(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(&TestFaultReader{}),
				Request:    req,
			}, nil
		})
		bd := New(testFetcher)
		_, err := bd.Query(ctx, "hello")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "faultreadererror")
	})
}

type TestFetcher struct {
	handler func(*http.Request) (*http.Response, error)
}

func NewTestFetcher(handler func(*http.Request) (*http.Response, error)) *TestFetcher {
	return &TestFetcher{
		handler: handler,
	}
}

func (tf *TestFetcher) Do(req *http.Request) (*http.Response, error) {
	return tf.handler(req)
}

type TestFaultReader struct{}

func (*TestFaultReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("faultreadererror")
}
