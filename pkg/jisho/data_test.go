package jisho

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Darkclainer/japwords/pkg/basicdict/basicdicttest"
	"github.com/Darkclainer/japwords/pkg/lemma"
)

// TestParseFiles checks how real pages are parsed.
func TestParseFiles(t *testing.T) {
	testCases := map[string]struct {
		// Count is a number of lemmas that parser should return
		Count int
		// Lemmas are specific lemmas (by index) that parser should return
		Lemmas map[int]*lemma.Lemma
	}{
		"犬": {
			Count: 20,
			Lemmas: map[int]*lemma.Lemma{
				0: {
					Slug: lemma.Word{
						Word:     "犬",
						Hiragana: "いぬ",
						Furigana: []lemma.FuriganaChar{
							{
								Kanji:    "犬",
								Hiragana: "いぬ",
							},
						},
					},
					Tags: []string{
						"Common word",
						"JLPT N5",
						"Wanikani level 2",
					},
					Forms: []lemma.Word{
						{
							Word:     "狗",
							Hiragana: "いぬ",
						},
						{
							Word: "イヌ",
						},
					},
					Senses: []lemma.WordSense{
						{
							Definition:   []string{"dog (Canis (lupus) familiaris)"},
							PartOfSpeech: []string{"Noun"},
						},
						{
							Definition: []string{
								"squealer",
								"rat",
								"snitch",
								"informer",
								"informant",
								"spy",
							},
							PartOfSpeech: []string{"Noun"},
							Tags: []string{
								"Derogatory",
								"Usually written using kana alone",
							},
						},
						{
							Definition: []string{
								"loser",
								"asshole",
							},
							PartOfSpeech: []string{"Noun"},
							Tags:         []string{"Derogatory"},
						},
						{
							Definition: []string{
								"counterfeit",
								"inferior",
								"useless",
								"wasteful",
							},
							PartOfSpeech: []string{"Noun, used as a prefix"},
						},
					},
					Audio: []lemma.Audio{
						{
							Type:   "audio/mpeg",
							Source: "https://d1vjc5dkcd3yh2.cloudfront.net/audio/10ce3f5eb7b4a9a03c4dafce2af60e28.mp3",
						},
						{
							Type:   "audio/ogg",
							Source: "https://d1vjc5dkcd3yh2.cloudfront.net/audio_ogg/10ce3f5eb7b4a9a03c4dafce2af60e28.ogg",
						},
					},
				},
				3: {
					Slug: lemma.Word{
						Word:     "犬も食わない",
						Hiragana: "いぬもくわない",
						Furigana: []lemma.FuriganaChar{
							{
								Kanji:    "犬",
								Hiragana: "いぬ",
							},
							{
								Hiragana: "も",
							},
							{
								Kanji:    "食",
								Hiragana: "く",
							},
							{
								Hiragana: "わ",
							},
							{
								Hiragana: "な",
							},
							{
								Hiragana: "い",
							},
						},
					},
					Forms: []lemma.Word{
						{
							Word:     "犬も喰わない",
							Hiragana: "いぬもくわない",
						},
						{
							Word:     "犬もくわない",
							Hiragana: "いぬもくわない",
						},
					},
					Senses: []lemma.WordSense{
						{
							Definition: []string{
								"avoided by everybody",
								"disliked by everybody",
								"not even a dog will eat (it)",
							},
							PartOfSpeech: []string{
								"Expressions (phrases, clauses, etc.)",
								"I-adjective (keiyoushi)",
							},
							Tags: []string{
								"Idiomatic expression",
							},
						},
					},
				},
			},
		},
	}
	jishoDict := New(nil, "")
	cacheDict := basicdicttest.New(t, "testdata", jishoDict.queryURL)
	for query := range testCases {
		tc := testCases[query]
		t.Run(query, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
			defer cancel()
			html := cacheDict.GetCachedHTML(t, ctx, query)
			lemmas, err := parseHTMLBytes(html)
			require.NoError(t, err)
			require.Equal(t, tc.Count, len(lemmas))
			for i, l := range tc.Lemmas {
				if i >= len(lemmas) {
					t.Fatalf("index %d is out of range (0, %d)", i, len(lemmas))
				}
				assert.Equal(t, l, lemmas[i])
			}
		})
	}
}
