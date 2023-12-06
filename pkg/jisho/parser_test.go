package jisho

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Darkclainer/japwords/pkg/htmltest"
	"github.com/Darkclainer/japwords/pkg/lemma"
)

func Test_parseHTML(t *testing.T) {
	testCases := []struct {
		Name        string
		HTML        string
		Expected    []*lemma.Lemma
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "simple",
			HTML: `
	<html><body>
		<div id="page_container"><div id="main_results"></div></div>
	</body></html>`,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "one lemma",
			HTML: `
	<html><body>
		<div id="page_container"><div id="main_results">
			<div id="primary"><div class="exact_block">
				<div class="concept_light">
					<div class="concept_light-wrapper">
						<div class="concept_light-readings">
							<div class="concept_light-representation">
								<span class="text">hello</span>
							</div>
						</div>
					</div>
				</div>
			</div></div>
		</div></div>
	</body></html>`,
			Expected: []*lemma.Lemma{
				{
					Slug: lemma.Word{
						Word: "hello",
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "no main_results",
			HTML: `
	<html><body>
		<div id="page_container"></div>
	</body></html>`,
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "main_results")
			},
		},
		{
			Name: "no page_container",
			HTML: `
	<html><body>
		<div id="main_results"></div>
	</body></html>`,
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "main_results")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			lemmas, err := parseHTML(strings.NewReader(tc.HTML))
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, lemmas)
		})
	}
}

func Test_parseMainResults(t *testing.T) {
	testCases := []struct {
		Name        string
		HTML        string
		Expected    []*lemma.Lemma
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name:        "simple",
			HTML:        `<div id="root"> </div>`,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "one lemma",
			HTML: `
	<div id="root">
	<div id="primary"><div class="exact_block">
		<div class="concept_light">
			<div class="concept_light-wrapper">
				<div class="concept_light-readings">
					<div class="concept_light-representation">
						<span class="text">hello</span>
					</div>
				</div>
			</div>
		</div>
	</div></div>
	</div>`,
			Expected: []*lemma.Lemma{
				{
					Slug: lemma.Word{
						Word: "hello",
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "two lemmas",
			HTML: `
	<div id="root">
	<div id="primary"><div class="exact_block">
		<div class="concept_light">
			<div class="concept_light-wrapper">
				<div class="concept_light-readings">
					<div class="concept_light-representation">
						<span class="text">hello</span>
					</div>
				</div>
			</div>
		</div>
		<div class="concept_light">
			<div class="concept_light-wrapper">
				<div class="concept_light-readings">
					<div class="concept_light-representation">
						<span class="text">world</span>
					</div>
				</div>
			</div>
		</div>
	</div></div>
	</div>`,
			Expected: []*lemma.Lemma{
				{
					Slug: lemma.Word{
						Word: "hello",
					},
				},
				{
					Slug: lemma.Word{
						Word: "world",
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "lemmas with error",
			HTML: `
	<div id="root">
	<div id="primary"><div class="exact_block">
		<div class="concept_light">
			<div class="concept_light-wrapper">
				<div class="concept_light-readings">
					<div class="concept_light-representation">
						<span class="text">hello</span>
					</div>
				</div>
			</div>
		</div>
		<div class="concept_light">
			<div class="concept_light-wrapper">
				<div class="concept_light-readings">
					<div class="concept_light-representation">
					</div>
				</div>
			</div>
		</div>
	</div></div>
	</div>`,
			Expected: []*lemma.Lemma{
				{
					Slug: lemma.Word{
						Word: "hello",
					},
				},
			},
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				var batchErr *LemmaBatchError
				if !assert.ErrorAs(t, err, &batchErr) {
					return false
				}
				if !assert.Len(t, batchErr.Errs, 1) {
					return false
				}
				var lemmaErr *LemmaError
				if !assert.ErrorAs(t, batchErr.Errs[0], &lemmaErr) {
					return false
				}
				return assert.Equal(t, 1, lemmaErr.ID)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			sel := htmltest.MustRootSelection(t, tc.HTML)
			lemmas, err := parseMainResults(sel)
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, lemmas)
		})
	}
}

func Test_parseConceptLight(t *testing.T) {
	testCases := []struct {
		Name        string
		HTML        string
		Expected    *lemma.Lemma
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "simple",
			HTML: `
		<div id="root"> 
			<div class="concept_light-wrapper">
				<div class="concept_light-readings">
					<div class="concept_light-representation">
						<span class="text">  he </span>
					</div>
				</div>
			</div>
		</div>`,

			Expected: &lemma.Lemma{
				Slug: lemma.Word{
					Word: "he",
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "real light wrapper",
			HTML: `
		<div id="root">
			<div class="concept_light-wrapper  columns zero-padding">
				<div class="concept_light-readings japanese japanese_gothic" lang="ja">
					<div class="concept_light-representation">      
						<span class="furigana"><span class="kanji-2-up kanji">いぬ</span></span>
						<span class="text">犬</span>
					</div>
				</div>
			</div>
		</div>`,
			Expected: &lemma.Lemma{
				Slug: lemma.Word{
					Word:     "犬",
					Furigana: newTestFurigana("犬", "いぬ"),
					Hiragana: "いぬ",
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "full",
			HTML: `
		<div id="root">
			<div class="concept_light-wrapper  columns zero-padding">
				<div class="concept_light-readings japanese japanese_gothic" lang="ja">
					<div class="concept_light-representation">
						<span class="furigana">
							<span class="kanji-2-up kanji">いぬ</span>
						</span>
						<span class="text">犬</span>
					</div>
				</div>
				<div class="concept_light-status">
					<span class="concept_light-tag concept_light-common success label">Common word</span>
					<span class="concept_light-tag label">
						<a href="http://wanikani.com/">Wanikani level 2</a>
					</span>
					<audio id="audio_犬:いぬ" preload="none">
						<source src="audio1" type="audio/mpeg">
						<source src="audio2" type="audio/ogg">
					</audio>
				</div>
			</div>
			<div class="concept_light-meanings medium-9 columns">
    				<div class="meanings-wrapper">
					<div class="meaning-tags">Noun</div>
					<div class="meaning-wrapper">
						<div class="meaning-definition zero-padding">
							<span class="meaning-definition-section_divider">1. </span>
							<span class="meaning-meaning">dog</span>
						</div>
					</div>
					<div class="meaning-tags">Noun</div>
					<div class="meaning-wrapper">
						<div class="meaning-definition zero-padding">
							<span class="meaning-definition-section_divider">2. </span>
							<span class="meaning-meaning">squealer; rat</span>
							<span class="supplemental_info">
								<span class="sense-tag tag-tag">Derogatory</span>,
								<span class="sense-tag tag-tag">Usually</span>, 
							</span>
						</div>
					</div>
					<div class="meaning-tags">Other forms</div>
					<div class="meaning-wrapper">
						<div class="meaning-definition zero-padding">
							<span class="meaning-meaning">
								<span class="break-unit">狗 【いぬ】</span>、
								<span class="break-unit">イヌ</span>
							</span>
						</div>
					</div>
				</div>
			</div>
		</div>`,
			Expected: &lemma.Lemma{
				Slug: lemma.Word{Word: "犬", Furigana: newTestFurigana("犬", "いぬ"), Hiragana: "いぬ"},
				Tags: []string{
					"Common word",
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
						Definition:   []string{"dog"},
						PartOfSpeech: []string{"Noun"},
					},
					{
						Definition:   []string{"squealer", "rat"},
						PartOfSpeech: []string{"Noun"},
						Tags:         []string{"Derogatory", "Usually"},
					},
				},
				Audio: []lemma.Audio{
					{
						Type:   "audio/mpeg",
						Source: "audio1",
					},
					{
						Type:   "audio/ogg",
						Source: "audio2",
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "long representation",
			HTML: `
		<div id="root">
			<div class="concept_light-wrapper concept_light-long_representation  columns zero-padding">
				<div class="concept_light-readings japanese japanese_gothic" lang="ja">
					<div class="concept_light-representation">
						<span class="furigana">
							<span class="kanji-2-up kanji">いぬ</span>
						</span>
						<span class="text">犬</span>
					</div>
				</div>
			</div>
			<div class="concept_light-status">
				<span class="concept_light-tag concept_light-common success label">Common word</span>
				<span class="concept_light-tag label">
					<a href="http://wanikani.com/">Wanikani level 2</a>
				</span>
				<audio id="audio_犬:いぬ" preload="none">
					<source src="audio1" type="audio/mpeg">
				</audio>
			</div>
			<div class="concept_light-meanings medium-9 columns">
    				<div class="meanings-wrapper">
					<div class="meaning-tags">Noun</div>
					<div class="meaning-wrapper">
						<div class="meaning-definition zero-padding">
							<span class="meaning-definition-section_divider">1. </span>
							<span class="meaning-meaning">dog</span>
						</div>
					</div>
				</div>
			</div>
		</div>`,
			Expected: &lemma.Lemma{
				Slug: lemma.Word{Word: "犬", Furigana: newTestFurigana("犬", "いぬ"), Hiragana: "いぬ"},
				Tags: []string{
					"Common word",
					"Wanikani level 2",
				},
				Senses: []lemma.WordSense{
					{
						Definition:   []string{"dog"},
						PartOfSpeech: []string{"Noun"},
					},
				},
				Audio: []lemma.Audio{
					{
						Type:   "audio/mpeg",
						Source: "audio1",
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "no slug",
			HTML: `
		<div id="root"> 
			<div class="concept_light-wrapper">
				<div class="concept_light-readings">
					<div class="concept_light-representation">
					</div>
				</div>
			</div>
		</div>`,

			ErrorAssert: assert.Error,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			sel := htmltest.MustRootSelection(t, tc.HTML)
			lemma, err := parseConceptLight(sel)
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
		Expected    lemma.Word
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "simple",
			HTML: ` <div id="root"> <span class="text">  he </span> </div> `,
			Expected: lemma.Word{
				Word: "he",
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "with furigana",
			HTML: `
			<div id="root">
				<span class="furigana">
					<span>げん</span>
					<span>き</span>
				</span>
				<span class="text">元気</span> 
			</div> `,
			Expected: lemma.Word{
				Word:     "元気",
				Furigana: newTestFurigana("元", "げん", "気", "き"),
				Hiragana: "げんき",
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "empty furigana",
			HTML: `
			<div id="root">
				<span class="furigana">
					<span></span>
					<span></span>
				</span>
				<span class="text">元気</span> 
			</div> `,
			Expected: lemma.Word{
				Word: "元気",
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "with furigana and hiragana",
			HTML: `
			<div id="root">
				<span class="furigana">
					<span>いぬ</span>
					<span></span>
					<span>ちる</span>
				</span>
				<span class="text">犬<span>と</span>猿</span> 
			</div> `,
			Expected: lemma.Word{
				Word:     "犬と猿",
				Furigana: newTestFurigana("犬", "いぬ", "", "と", "猿", "ちる"),
				Hiragana: "いぬとちる",
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name:        "no representation",
			HTML:        ` <div id="root"> <span class="text">  </span> </div> `,
			ErrorAssert: assert.Error,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			sel := htmltest.MustRootSelection(t, tc.HTML)
			slug, err := parseRepresentation(sel)
			tc.ErrorAssert(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.Expected, slug)
		})
	}
}

func Test_parseStatus(t *testing.T) {
	testCases := []struct {
		Name  string
		HTML  string
		Audio []lemma.Audio
		Tags  []string
	}{
		{
			Name: "nothing",
			HTML: ` <div id="root"></div> `,
		},
		{
			Name: "tags",
			HTML: `
		<div id="root">
			<span class="concept_light-tag concept_light-common success label">Common word</span>
			<span class="concept_light-tag label">JLPT N3</span>
			<span class="concept_light-tag label"><a href="http://wanikani.com/">Wanikani level 22</a></span> 
		</div>`,
			Tags: []string{
				"Common word",
				"JLPT N3",
				"Wanikani level 22",
			},
		},
		{
			Name: "known audio",
			HTML: `
		<div id="root">
			<audio>
				<source type="audio/mpeg" src="https://example.com/file.mp3">
				<source type="audio/ogg" src="https://example.com/file.ogg">
			</audio>
		
		</div>`,
			Audio: []lemma.Audio{
				{Type: "audio/mpeg", Source: "https://example.com/file.mp3"},
				{Type: "audio/ogg", Source: "https://example.com/file.ogg"},
			},
		},
		{
			Name: "unknown audio",
			HTML: `
		<div id="root">
			<audio>
				<source src="https://example.com/file.mp3">
			</audio>
		
		</div>`,
			Audio: []lemma.Audio{
				{Type: "unknown", Source: "https://example.com/file.mp3"},
			},
		},
		{
			Name: "audio protocol indepenent",
			HTML: `
		<div id="root">
			<audio>
				<source src="//example.com/file.mp3">
			</audio>
		
		</div>`,
			Audio: []lemma.Audio{
				{Type: "unknown", Source: "https://example.com/file.mp3"},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			sel := htmltest.MustRootSelection(t, tc.HTML)
			audio, tags := parseStatus(sel)
			assert.Equal(t, tc.Audio, audio)
			assert.Equal(t, tc.Tags, tags)
		})
	}
}

func Test_parseMeanings(t *testing.T) {
	testCases := []struct {
		Name   string
		Src    string
		Senses []lemma.WordSense
		Forms  []lemma.Word
	}{
		{
			Name: "empty",
			Src:  `<div id="root"></div>`,
		},
		{
			Name: "one definition without tag",
			Src: `
		<div id="root">
			<div class="meaning-wrapper">
				<div class="meaning-definition">
					<span class="meaning-meaning">hello</span>
				</div>
			</div>
		</div>`,
			Senses: []lemma.WordSense{
				{
					Definition: []string{"hello"},
				},
			},
		},
		{
			Name: "tag without definition",
			Src: `
		<div id="root">
			<div class="meaning-tags">Noun</div>
			<div class="meaning-wrapper">
			</div>
		</div>`,
		},
		{
			Name: "two definitions",
			Src: `
		<div id="root">
				<div class="meaning-tags">Noun</div>
				<div class="meaning-wrapper">
					<div class="meaning-definition">
						<span class="meaning-meaning">foo</span>
					</div>
				</div>
				<div class="meaning-tags">Adj</div>
				<div class="meaning-wrapper">
					<div class="meaning-definition">
						<span class="meaning-meaning">bar</span>
					</div>
				</div>
		</div>`,
			Senses: []lemma.WordSense{
				{
					Definition:   []string{"foo"},
					PartOfSpeech: []string{"Noun"},
				},
				{
					Definition:   []string{"bar"},
					PartOfSpeech: []string{"Adj"},
				},
			},
		},
		{
			Name: "two tags with definition",
			Src: `
		<div id="root">
			<div class="meaning-tags">Adj</div>
			<div class="meaning-tags">Noun</div>
			<div class="meaning-wrapper">
				<div class="meaning-definition">
					<span class="meaning-meaning">hello</span>
				</div>
			</div>
		</div>`,
			Senses: []lemma.WordSense{
				{
					Definition:   []string{"hello"},
					PartOfSpeech: []string{"Noun"},
				},
			},
		},
		{
			Name: "one definition with POS",
			Src: `
		<div id="root">
			<div class="meaning-tags">Noun</div>
			<div class="meaning-wrapper">
				<div class="meaning-definition">
					<span class="meaning-definition-section_divider">1. </span>
					<span class="meaning-meaning">squealer; rat</span>
					<span class="supplemental_info">
						<span class="sense-tag tag-tag">Derogatory</span>, 
						<span class="sense-tag tag-tag">Something</span>, 
					</span>
				</div>
			</div>
		</div>`,
			Senses: []lemma.WordSense{
				{
					Definition:   []string{"squealer", "rat"},
					PartOfSpeech: []string{"Noun"},
					Tags:         []string{"Derogatory", "Something"},
				},
			},
		},
		{
			Name: "one definition with other form",
			Src: `
		<div id="root">
			<div class="meaning-tags">Noun</div>
			<div class="meaning-wrapper">
				<div class="meaning-definition">
					<span class="meaning-meaning">hello</span>
				</div>
			</div>
			<div class="meaning-tags">Other forms</div>
			<div class="meaning-wrapper">
				<div class="meaning-definition">
					<span class="meaning-meaning">
						<span class="break-unit">狗 【いぬ】</span>
					</span>
				</div>
			</div>
		</div>`,
			Senses: []lemma.WordSense{
				{
					Definition:   []string{"hello"},
					PartOfSpeech: []string{"Noun"},
				},
			},
			Forms: []lemma.Word{
				{
					Word:     "狗",
					Hiragana: "いぬ",
				},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			senses, forms := parseMeanings(
				htmltest.MustRootSelection(t, tc.Src),
			)
			assert.Equal(t, tc.Senses, senses, "senses")
			assert.Equal(t, tc.Forms, forms, "forms")
		})
	}
}

func Test_parseMeaningDefinition(t *testing.T) {
	testCases := []struct {
		Name        string
		Src         string
		Definitions []string
		Tags        []string
	}{
		{
			Name: "empty",
			Src:  `<div id="root"></div>`,
		},
		{
			Name: "single definition",
			Src: `
		<div id="root">
			<span class="meaning-meaning">Hello world!</span>
		</div>`,
			Definitions: []string{"Hello world!"},
		},
		{
			Name: "several definitions",
			Src: `
		<div id="root">
			<span class="meaning-meaning">Hello; world!</span>
		</div>`,
			Definitions: []string{"Hello", "world!"},
		},
		{
			Name: "single tag",
			Src: `
		<div id="root">
			<span class="supplemental_info">
				<span class="tag-tag">hello</span>
			</span>
		</div>`,
			Tags: []string{"hello"},
		},
		{
			Name: "several tags",
			Src: `
		<div id="root">
			<span class="supplemental_info">
				<span class="tag-tag">hello</span>
				<span class="tag-tag">world</span>
			</span>
		</div>`,
			Tags: []string{"hello", "world"},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			definitions, tags := parseMeaningDefinition(
				htmltest.MustRootSelection(t, tc.Src),
			)
			assert.Equal(t, tc.Definitions, definitions, "definitions")
			assert.Equal(t, tc.Tags, tags, "tags")
		})
	}
}

func Test_parseDefinitionOtherForms(t *testing.T) {
	testCases := []struct {
		Name  string
		Src   string
		Forms []lemma.Word
	}{
		{
			Name: "empty",
			Src:  `<div id="root"></div>`,
		},
		{
			Name: "one form",
			Src: `
		<div id="root">
			<span class="meaning-meaning">
				<span class="break-unit">狗 【いぬ】</span>
			</span>
		</div>`,
			Forms: []lemma.Word{
				{
					Word:     "狗",
					Hiragana: "いぬ",
				},
			},
		},
		{
			Name: "two forms",
			Src: `
		<div id="root">
			<span class="meaning-meaning">
				<span class="break-unit">狗 【いぬ】</span>
				<span class="break-unit">イヌ</span>
			</span>
		</div>`,
			Forms: []lemma.Word{
				{
					Word:     "狗",
					Hiragana: "いぬ",
				},
				{
					Word: "イヌ",
				},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			forms := parseDefinitionOtherForms(
				htmltest.MustRootSelection(t, tc.Src),
			)
			assert.Equal(t, tc.Forms, forms, "forms")
		})
	}
}

func Test_splitPartOfSpeech(t *testing.T) {
	testCases := []struct {
		Name     string
		Src      string
		Expected []string
	}{
		{
			Name: "empty",
		},
		{
			Name:     "single word",
			Src:      "Noun",
			Expected: []string{"Noun"},
		},
		{
			Name:     "single POS with comma",
			Src:      "Noun, test",
			Expected: []string{"Noun, test"},
		},
		{
			Name:     "two POS",
			Src:      "Noun, Ad",
			Expected: []string{"Noun", "Ad"},
		},
		{
			Name:     "three POS",
			Src:      "Noun, Adjective, Something with some",
			Expected: []string{"Noun", "Adjective", "Something with some"},
		},
		{
			Name:     "end edge case",
			Src:      "Noun, A",
			Expected: []string{"Noun", "A"},
		},
		{
			Name:     "one letter",
			Src:      "N",
			Expected: []string{"N"},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			got := splitPartOfSpeech(tc.Src)
			assert.Equal(t, tc.Expected, got)
		})
	}
}

func Test_parseOtherForm(t *testing.T) {
	testCases := []struct {
		Name  string
		Src   string
		Word  lemma.Word
		Exist bool
	}{
		{
			Name: "empty",
		},
		{
			Name: "kanji without reading",
			Src:  "hello",
			Word: lemma.Word{
				Word: "hello",
			},
			Exist: true,
		},
		{
			Name: "kanji with reading",
			Src:  "hello 【world】",
			Word: lemma.Word{
				Word:     "hello",
				Hiragana: "world",
			},
			Exist: true,
		},
		{
			Name: "kanji with reading without whitespace",
			Src:  "hello【world】",
			Word: lemma.Word{
				Word:     "hello",
				Hiragana: "world",
			},
			Exist: true,
		},
		{
			Name: "kanji with unfinished reading",
			Src:  "hello【world",
			Word: lemma.Word{
				Word: "hello",
			},
			Exist: true,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			word, exist := parseOtherForm(tc.Src)
			assert.Equal(t, tc.Exist, exist, "exist")
			assert.Equal(t, tc.Word, word, "word")
		})
	}
}

func newTestFurigana(parts ...string) lemma.Furigana {
	if len(parts)%2 == 1 {
		panic("number of parts should be even")
	}
	var furigana lemma.Furigana
	for i := 1; i < len(parts); i += 2 {
		furigana = append(furigana, lemma.FuriganaChar{
			Kanji:    parts[i-1],
			Hiragana: parts[i],
		})
	}
	return furigana
}
