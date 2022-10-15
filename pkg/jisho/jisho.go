package jisho

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://jisho.org/search/"

type Jisho struct {
	client  Fetcher
	baseURL string
}

type Fetcher interface {
	Do(*http.Request) (*http.Response, error)
}

func New(client Fetcher, baseURL string) *Jisho {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Jisho{
		client:  client,
		baseURL: baseURL,
	}
}

func (j *Jisho) Query(ctx context.Context, query string) ([]*Lemma, error) {
	// htmlBody, err := j.queryHTML(ctx, query)
	return nil, errors.New("unimplemented")
}

func (j *Jisho) queryHTML(ctx context.Context, query string) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		j.queryURL(query),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("request construction failed: %w", err)
	}
	resp, err := j.client.Do(req)
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

func (j *Jisho) queryURL(query string) string {
	return j.baseURL + url.PathEscape(query)
}
