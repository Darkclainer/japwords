package wadoku

import (
	"context"
	"net/url"

	"github.com/Darkclainer/japwords/pkg/lemma"
)

const defaultBaseURL = "https://www.wadoku.de/search/"

type Wadoku struct {
	client  BasicDict
	baseURL string
}
type BasicDict interface {
	Query(context.Context, string) ([]byte, error)
}

func New(client BasicDict, baseURL string) *Wadoku {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Wadoku{
		client:  client,
		baseURL: baseURL,
	}
}

func (w *Wadoku) Query(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
	url := w.queryURL(query)
	htmlBody, err := w.client.Query(ctx, url)
	if err != nil {
		return nil, err
	}
	return parseHTMLBytes(htmlBody)
}

func (j *Wadoku) queryURL(query string) string {
	return j.baseURL + url.PathEscape(query)
}
