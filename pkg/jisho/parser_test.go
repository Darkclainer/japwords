package jisho

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseConceptLight(t *testing.T) {
	testCases := []struct {
		Name        string
		HTML        string
		Expected    *Lemma
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "simple",
			HTML: `
		<div> 
			<div class="concept_light-wrapper">
				<div class="concept_light-readings">
					<div class="concept_light-representation">
						<span class="text">  he </span>
					</div>
				</div>
			</div>
		</div>`,

			Expected: &Lemma{
				Slug: Word{
					Word:     "he",
					Furigana: newTestFurigana("h", "", "e", ""),
					Reading:  "he",
				},
				Audio: map[string]string{},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "real",
			HTML: `
		<div>
			<div class="concept_light-wrapper  columns zero-padding">
				<div class="concept_light-readings japanese japanese_gothic" lang="ja">
					<div class="concept_light-representation">      
						<span class="furigana"><span class="kanji-2-up kanji">いぬ</span></span>
						<span class="text">犬</span>
					</div>
				</div>
			</div>
		</div>`,
			Expected: &Lemma{
				Slug: Word{
					Word:     "犬",
					Furigana: newTestFurigana("犬", "いぬ"),
					Reading:  "いぬ",
				},
				Audio: map[string]string{},
			},
			ErrorAssert: assert.NoError,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			sel := mustDocument(t, tc.HTML)
			lemma, err := parseConceptLight(sel.Selection)
			tc.ErrorAssert(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.Expected, lemma)
		})
	}
}

func Test_parseRepresentation(t *testing.T) {
	testCases := []struct {
		Name        string
		HTML        string
		Expected    Word
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "simple",
			HTML: ` <div> <span class="text">  he </span> </div> `,
			Expected: Word{
				Word:     "he",
				Furigana: newTestFurigana("h", "", "e", ""),
				Reading:  "he",
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "with furigana",
			HTML: `
			<div>
				<span class="furigana">
					<span>げん</span>
					<span>き</span>
				</span>
				<span class="text">元気</span> 
			</div> `,
			Expected: Word{
				Word:     "元気",
				Furigana: newTestFurigana("元", "げん", "気", "き"),
				Reading:  "げんき",
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "empty furigana",
			HTML: `
			<div>
				<span class="furigana">
					<span></span>
					<span></span>
				</span>
				<span class="text">元気</span> 
			</div> `,
			Expected: Word{
				Word:     "元気",
				Furigana: newTestFurigana("元", "", "気", ""),
				Reading:  "元気",
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name:        "no representation",
			HTML:        ` <div> <span class="text">  </span> </div> `,
			ErrorAssert: assert.Error,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			sel := mustDocument(t, tc.HTML)
			slug, err := parseRepresentation(sel.Selection)
			tc.ErrorAssert(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.Expected, slug)
		})
	}
}

func newTestFurigana(parts ...string) Furigana {
	if len(parts)%2 == 1 {
		panic("number of parts should be even")
	}
	var furigana Furigana
	for i := 1; i < len(parts); i += 2 {
		furigana = append(furigana, FuriganaChar{
			Kanji:    parts[i-1],
			Hiragana: parts[i],
		})
	}
	return furigana
}

func mustDocument(t *testing.T, src string) *goquery.Document {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(src))
	require.NoError(t, err)
	return doc
}
