package fetcher

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomHeader(t *testing.T) {
	serverAddr := newTestServer(t, http.HandlerFunc(
		func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "test-user-agent", r.Header.Get("User-Agent"))
		},
	))

	fetcher, err := New(In{
		Config: &Config{
			Headers: map[string]string{
				"User-Agent": "test-user-agent",
			},
		},
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, "http://"+serverAddr, nil)
	require.NoError(t, err)

	_, err = fetcher.Do(req)
	require.NoError(t, err)
}

func newTestServer(t *testing.T, handler http.Handler) string {
	listener := newTestListener(t)
	var server http.Server
	server.Addr = listener.Addr().String()
	server.Handler = handler

	waitClose := make(chan struct{})
	go func() {
		defer func() {
			waitClose <- struct{}{}
		}()
		if err := server.Serve(listener); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			t.Errorf("http server closed with error: %s", err)
		}
	}()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err := server.Shutdown(ctx)
		require.NoError(t, err)
		<-waitClose
	})
	return server.Addr
}

func newTestListener(t *testing.T) net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			t.Fatalf("failed to listen on port: %s", err)
		}
	}
	return l
}
