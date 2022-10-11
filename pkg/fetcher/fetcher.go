package fetcher

import (
	"net/http"
	"time"

	"go.uber.org/fx"
)

// Fetcher is a wrapper for http.Client that provides means
// for configuration all http requests, for example providing
// custom user agent.
type Fetcher struct {
	client       http.Client
	extraHeaders map[string]string
}

type In struct {
	fx.In

	Config *Config
}

func New(in In) (*Fetcher, error) {
	fetcher := Fetcher{
		client: http.Client{
			Timeout: time.Second * 30,
		},
		extraHeaders: in.Config.Headers,
	}
	return &fetcher, nil
}

func (f *Fetcher) Do(req *http.Request) (*http.Response, error) {
	for k, v := range f.extraHeaders {
		req.Header.Set(k, v)
	}
	return f.client.Do(req)
}
