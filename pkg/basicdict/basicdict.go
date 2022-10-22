package basicdict

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type BasicDict struct {
	client Fetcher
}

type Fetcher interface {
	Do(*http.Request) (*http.Response, error)
}

func New(f Fetcher) *BasicDict {
	return &BasicDict{
		client: f,
	}
}

func (bd *BasicDict) Query(ctx context.Context, query string) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		query,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("request construction failed: %w", err)
	}
	resp, err := bd.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status %d != 200", resp.StatusCode)
	}
	var body bytes.Buffer
	if _, err := io.Copy(&body, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body.Bytes(), nil
}
