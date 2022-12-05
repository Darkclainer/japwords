package basicdicttest

import (
	"context"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Darkclainer/japwords/pkg/basicdict"
	"github.com/Darkclainer/japwords/pkg/fetcher"
)

// restoreCache used to download response on query if requested files are not presented
var restoreCache = flag.Bool("restore", false, "restore html files")

// CachedDictionary is used for test, to get page and save it in testdata
type CachedDictionary struct {
	client       *basicdict.BasicDict
	cacheDir     string
	queryHandler func(string) string
}

func New(t *testing.T, cacheDir string, queryHandler func(string) string) *CachedDictionary {
	fetcherClient, err := fetcher.New(&fetcher.Config{})
	require.NoError(t, err)
	client := basicdict.New(fetcherClient)
	return &CachedDictionary{
		client:       client,
		cacheDir:     cacheDir,
		queryHandler: queryHandler,
	}
}

func (cd *CachedDictionary) GetCachedHTML(t *testing.T, ctx context.Context, query string) []byte {
	t.Helper()
	path := cd.getHTMLName(query)
	html, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) && *restoreCache {
		return cd.restoreHTML(t, ctx, query)
	}
	if err != nil {
		t.Fatalf("file for query %q not found, use -restore flag to create it", query)
	}
	return html
}

func (cd *CachedDictionary) restoreHTML(t *testing.T, ctx context.Context, query string) []byte {
	html, err := cd.client.Query(ctx, cd.queryHandler(query))
	require.NoError(t, err)
	err = os.WriteFile(cd.getHTMLName(query), html, 0o540)
	require.NoError(t, err)
	return html
}

func (cd *CachedDictionary) getHTMLName(query string) string {
	return filepath.Join(cd.cacheDir, query+".html")
}
