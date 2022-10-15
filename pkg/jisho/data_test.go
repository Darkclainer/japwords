package jisho

import (
	"context"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"japwords/pkg/fetcher"
)

// restoreCache used to download response on query if requested files are not presented
var restoreCache = flag.Bool("restore", false, "restore html files")

// TestParseFiles is meant how real overall pages are parsed.
func TestParseFiles(t *testing.T) {
	testCases := map[string]struct {
		// Slugs are all Lemma's slug that parser should return
		Slugs []string
		// Lemmas are specific lemmas (by index) that parser should return
		Lemmas map[int]*Lemma
	}{
		"犬": {},
	}
	for query := range testCases {
		tc := testCases[query]
		t.Run(query, func(t *testing.T) {
			html := getCachedHTML(t, query)
			lemmas, err := parseHTML(html)
			require.NoError(t, err)
			var gotSlugs []string
			for _, l := range lemmas {
				gotSlugs = append(gotSlugs, l.Slug.Word)
			}
			require.Equal(t, tc.Slugs, gotSlugs)
			for i, l := range tc.Lemmas {
				assert.Equal(t, l, lemmas[i])
			}
		})
	}
}

func getCachedHTML(t *testing.T, query string) []byte {
	t.Helper()
	path := getHTMLName(query)
	html, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) && *restoreCache {
		return restoreHTML(t, query)
	}
	if err != nil {
		t.Fatalf("file for query %q not found, use -restore flag to create it", query)
	}
	return html
}

func restoreHTML(t *testing.T, query string) []byte {
	client, err := fetcher.New(fetcher.In{
		Config: &fetcher.Config{},
	})
	require.NoError(t, err)
	j := New(client, "")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	html, err := j.queryHTML(ctx, query)
	require.NoError(t, err)
	err = os.WriteFile(getHTMLName(query), html, 0o540)
	require.NoError(t, err)
	return html
}

func getHTMLName(query string) string {
	return filepath.Join("testdata", query+".html")
}
