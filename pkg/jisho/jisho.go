package jisho

import (
	"context"
	"net/url"

	"github.com/Darkclainer/japwords/pkg/lemma"
)

const defaultBaseURL = "https://jisho.org/search/"

type Jisho struct {
	client  BasicDict
	baseURL string
}

type BasicDict interface {
	Query(context.Context, string) ([]byte, error)
}

func New(client BasicDict, baseURL string) *Jisho {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Jisho{
		client:  client,
		baseURL: baseURL,
	}
}

func (j *Jisho) Query(ctx context.Context, query string) ([]*lemma.Lemma, error) {
	url := j.queryURL(query)
	htmlBody, err := j.client.Query(ctx, url)
	if err != nil {
		return nil, err
	}
	return parseHTMLBytes(htmlBody)
}

func (j *Jisho) queryURL(query string) string {
	return j.baseURL + url.PathEscape(query)
}
