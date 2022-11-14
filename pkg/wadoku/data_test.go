package wadoku

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"japwords/pkg/basicdict/basicdicttest"
	"japwords/pkg/lemma"
)

// TestParseFiles checks how web page is parsed.
//
// NOTE: wadoku has non determenistic results for query.
// For example if you try to search `七` then order of items that
// read as `なな` and as `しち` is not determined.
// Not only order of elements can be different, but some result can appear or disappear.
// This behaviour pose a problem when you must redownload web pages.
func TestParseFiles(t *testing.T) {
	testCases := map[string]struct {
		// Count is a number of lemmas that parser should return
		Count int
		// Lemmas are specific lemmas that should be found in results
		// and should be equal to lemmas in results. They are searched
		// by slug
		Lemmas []*lemma.WadokuLemma
	}{
		"犬": {
			Lemmas: []*lemma.WadokuLemma{
				{
					Slug:     "犬",
					Hiragana: "いぬ",
					Pitches: []lemma.Pitch{
						{
							Position: 3,
							IsHigh:   false,
						},
						{
							Position: 6,
							IsHigh:   true,
						},
						{
							Position: 6,
							IsHigh:   false,
						},
					},
				},
				{
					Slug:     "犬走り",
					Hiragana: "いぬばしり",
					Pitches: []lemma.Pitch{
						{
							Position: 3,
							IsHigh:   false,
						},
						{
							Position: 15,
							IsHigh:   true,
						},
					},
				},
			},
		},
	}
	wadokuDict := New(nil, "")
	cacheDict := basicdicttest.New(t, "testdata", wadokuDict.queryURL)
	for query := range testCases {
		tc := testCases[query]
		t.Run(query, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
			defer cancel()
			html := cacheDict.GetCachedHTML(t, ctx, query)
			lemmas, err := parseHTMLBytes(html)
			require.NoError(t, err)
			// store all result lemmas in map where key is slug.
			// I store them in reversed order to deal with duplicates (if there are some)
			results := map[string]*lemma.WadokuLemma{}
			for i := len(lemmas) - 1; i >= 0; i-- {
				results[lemmas[i].Slug] = lemmas[i]
			}
			for _, expectedLemma := range tc.Lemmas {
				resultLemma, ok := results[expectedLemma.Slug]
				if !ok {
					t.Errorf("no lemma with slug %q found in results", expectedLemma.Slug)
					continue
				}
				assert.Equal(t, expectedLemma, resultLemma)
			}
		})
	}
}
