package anki

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Darkclainer/japwords/pkg/lemma"
)

func mustConvertMapping(t *testing.T, mapping map[string]string) TemplateMapping {
	templateMapping, mappingErrs := convertMapping(mapping)
	for _, mappingErr := range mappingErrs {
		t.Fatalf("mustConvertMapping failed on template with key %s: %s", mappingErr.Key, mappingErr.Msg)
	}
	return templateMapping
}

func Test_RenderRawTemplate(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		result, err := RenderRawTemplate("{{.Slug.Word}}-test", &lemma.ProjectedLemma{
			Slug: lemma.Word{
				Word: "hello",
			},
		})
		require.NoError(t, err)
		assert.Equal(t, "hello-test", result)
	})
	t.Run("Error", func(t *testing.T) {
		_, err := RenderRawTemplate("{{.NotExists}}-test", &lemma.ProjectedLemma{
			Slug: lemma.Word{
				Word: "hello",
			},
		})
		require.Error(t, err)
	})
}

func Test_convertMapping(t *testing.T) {
	testLemma := &lemma.ProjectedLemma{
		Slug: lemma.Word{
			Word:     "hello",
			Hiragana: "world",
		},
	}
	testCases := []struct {
		Name            string
		SrcMapping      map[string]string
		ExpectedMapping TemplateMapping
		RenderExpected  map[string]string
		ErrorAssert     assert.ValueAssertionFunc
	}{
		{
			Name:            "empty",
			ExpectedMapping: map[string]*Template{},
			RenderExpected:  map[string]string{},
			ErrorAssert:     assert.Empty,
		},
		{
			Name: "one key",
			SrcMapping: map[string]string{
				"key": `{{.Slug.Word}}`,
			},
			ExpectedMapping: map[string]*Template{
				"key": {
					Src: `{{.Slug.Word}}`,
				},
			},
			RenderExpected: map[string]string{
				"key": "hello",
			},
			ErrorAssert: assert.Empty,
		},
		{
			Name: "sprig function",
			SrcMapping: map[string]string{
				"key": `{{upper .Slug.Word}}`,
			},
			ExpectedMapping: map[string]*Template{
				"key": {
					Src: `{{upper .Slug.Word}}`,
				},
			},
			RenderExpected: map[string]string{
				"key": "HELLO",
			},
			ErrorAssert: assert.Empty,
		},
		{
			Name: "two key",
			SrcMapping: map[string]string{
				"key1": `{{.Slug.Word}}`,
				"key2": `{{.Slug.Hiragana}}`,
			},
			ExpectedMapping: map[string]*Template{
				"key1": {
					Src: `{{.Slug.Word}}`,
				},
				"key2": {
					Src: `{{.Slug.Hiragana}}`,
				},
			},
			RenderExpected: map[string]string{
				"key1": "hello",
				"key2": "world",
			},
			ErrorAssert: assert.Empty,
		},
		{
			Name: "error",
			SrcMapping: map[string]string{
				"key": `{{.NotExist}}`,
			},
			ErrorAssert: assert.NotEmpty,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			actualMapping, err := convertMapping(tc.SrcMapping)
			tc.ErrorAssert(t, err)
			if err != nil {
				return
			}
			assert.True(t, tc.ExpectedMapping.Equal(actualMapping))
			renderActual := map[string]string{}
			for key, tmpl := range actualMapping {
				var buffer bytes.Buffer
				err := tmpl.Tmpl.Execute(&buffer, testLemma)
				if err != nil {
					t.Errorf("executing mapping with key %s failed: %s", key, err)
					continue
				}
				renderActual[key] = buffer.String()
			}
			assert.Equal(t, tc.RenderExpected, renderActual)
		})
	}
}

func Test_initTemplate(t *testing.T) {
	testCases := []struct {
		Name        string
		Tmpl        string
		Lemma       lemma.ProjectedLemma
		Expected    string
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name:        "empty",
			ErrorAssert: assert.NoError,
		},
		{
			Name:        "check template return error",
			Tmpl:        "{{.notexists}}",
			ErrorAssert: assert.Error,
		},
		{
			Name: "sprig functions imported",
			Lemma: lemma.ProjectedLemma{
				Slug: lemma.Word{
					Word: "hello",
				},
			},
			Tmpl:        "{{upper .Slug.Word}}",
			Expected:    "HELLO",
			ErrorAssert: assert.NoError,
		},
		{
			Name: "renderFurigana",
			Lemma: lemma.ProjectedLemma{
				Slug: lemma.Word{
					Furigana: lemma.Furigana{
						{
							Kanji:    "he",
							Hiragana: "12",
						},
						{
							Kanji:    "llo",
							Hiragana: "345",
						},
					},
				},
			},
			Tmpl:        "{{renderFurigana .Slug}}",
			Expected:    "he[12]llo[345]",
			ErrorAssert: assert.NoError,
		},
		{
			Name: "renderPitch",
			Lemma: lemma.ProjectedLemma{
				Slug: lemma.Word{
					Hiragana: "hello",
					PitchShapes: []lemma.PitchShape{
						{
							Hiragana: "h",
							Directions: []lemma.AccentDirection{
								lemma.AccentDirectionDown,
							},
						},
						{
							Hiragana: "ello",
							Directions: []lemma.AccentDirection{
								lemma.AccentDirectionUp,
								lemma.AccentDirectionLeft,
							},
						},
					},
				},
			},
			Tmpl:        `{{renderPitch .Slug "span" "u" "r" "d" "l"}}`,
			Expected:    `<span class="d">h</span><span class="u l">ello</span>`,
			ErrorAssert: assert.NoError,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			tmpl := template.New("testtemplate")
			err := initTemplate(tmpl, tc.Tmpl)
			tc.ErrorAssert(t, err)
			if err != nil {
				return
			}
			var buffer bytes.Buffer
			err = tmpl.Execute(&buffer, &tc.Lemma)
			require.NoError(t, err)
			assert.Equal(t, tc.Expected, buffer.String())
		})
	}
}

// Test_checkTemplate tests that unexisting fields returns errors
func Test_checkTemplate(t *testing.T) {
	testCases := []struct {
		Name        string
		Tmpl        string
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name:        "empty",
			Tmpl:        "",
			ErrorAssert: assert.NoError,
		},
		{
			Name:        "Slug.Word",
			Tmpl:        "{{.Slug.Word}}",
			ErrorAssert: assert.NoError,
		},
		{
			Name:        "Not existing field",
			Tmpl:        `{{.hello}}`,
			ErrorAssert: assert.Error,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			tmpl := template.Must(template.New("testtemplate").Parse(tc.Tmpl))
			tmpl.Option("missingkey=error")
			err := checkTemplate(tmpl)
			tc.ErrorAssert(t, err)
		})
	}
}

func Test_renderFurigana(t *testing.T) {
	testCases := []struct {
		Name     string
		Furigana lemma.Furigana
		Expected string
	}{
		{
			Name: "empty",
		},
		{
			Name: "kanji single",
			Furigana: lemma.Furigana{
				{
					Kanji: "hel",
				},
			},
			Expected: "hel",
		},
		{
			Name: "kanji many",
			Furigana: lemma.Furigana{
				{
					Kanji: "hel",
				},
				{
					Kanji: "lo",
				},
			},
			Expected: "hello",
		},
		{
			Name: "hiraga single",
			Furigana: lemma.Furigana{
				{
					Hiragana: "hel",
				},
			},
			Expected: "hel",
		},
		{
			Name: "hiraga many",
			Furigana: lemma.Furigana{
				{
					Hiragana: "hel",
				},
				{
					Hiragana: "lo",
				},
			},
			Expected: "hello",
		},
		{
			Name: "hiraga kanji",
			Furigana: lemma.Furigana{
				{
					Hiragana: "he",
				},
				{
					Kanji: "l",
				},
				{
					Hiragana: "lo",
				},
			},
			Expected: "hello",
		},
		{
			Name: "furigana",
			Furigana: lemma.Furigana{
				{
					Hiragana: "hel",
					Kanji:    "12",
				},
				{
					Hiragana: "lo",
					Kanji:    "34",
				},
			},
			Expected: "12[hel]34[lo]",
		},
		{
			Name: "mixed",
			Furigana: lemma.Furigana{
				{
					Hiragana: "hel",
					Kanji:    "12",
				},
				{
					Hiragana: "lo",
					Kanji:    "34",
				},
				{
					Hiragana: "w",
				},
				{
					Kanji: "o",
				},
				{
					Hiragana: "rld",
					Kanji:    "56",
				},
			},
			Expected: "12[hel]34[lo]wo 56[rld]",
		},
		{
			Name: "group delimiting",
			Furigana: lemma.Furigana{
				{
					Kanji:    "12",
					Hiragana: "hel",
				},
				{
					Hiragana: "w",
				},
				{
					Kanji:    "34",
					Hiragana: "lo",
				},
				{
					Hiragana: "w",
				},
				{
					Kanji:    "56",
					Hiragana: "rld",
				},
				{
					Kanji:    "789",
					Hiragana: "foo",
				},
			},
			Expected: "12[hel]w 34[lo]w 56[rld]789[foo]",
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(t.Name(), func(t *testing.T) {
			word := lemma.Word{
				Furigana: tc.Furigana,
			}
			actual := renderFuriganaTemplate(&word)
			assert.Equal(t, tc.Expected, actual)
		})

	}
}

func Test_renderPitch(t *testing.T) {
	directionClasses := []string{"u", "r", "d", "l"}
	testCases := []struct {
		Name             string
		Word             lemma.Word
		Tag              string
		DirectionClasses []string
		Expected         string
	}{
		{
			Name:             "empty",
			Tag:              "span",
			DirectionClasses: directionClasses,
		},
		{
			Name: "no pitch",
			Word: lemma.Word{
				Hiragana: "hello",
				PitchShapes: []lemma.PitchShape{
					{
						Hiragana: "hello",
					},
				},
			},
			DirectionClasses: directionClasses,
			Tag:              "span",
			Expected:         `<span class="">hello</span>`,
		},
		{
			Name: "simple up",
			Word: lemma.Word{
				Hiragana: "hello",
				PitchShapes: []lemma.PitchShape{
					{
						Hiragana: "hello",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionUp,
						},
					},
				},
			},
			Tag:              "span",
			Expected:         `<span class="u">hello</span>`,
			DirectionClasses: directionClasses,
		},
		{
			Name: "tail up down",
			Word: lemma.Word{
				Hiragana: "hello",
				PitchShapes: []lemma.PitchShape{
					{
						Hiragana: "hello",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionUp,
							lemma.AccentDirectionRight,
						},
					},
				},
			},
			Tag:              "span",
			Expected:         `<span class="u r">hello</span>`,
			DirectionClasses: directionClasses,
		},
		{
			Name: "head down up",
			Word: lemma.Word{
				Hiragana: "hello",
				PitchShapes: []lemma.PitchShape{
					{
						Hiragana: "hello",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionUp,
							lemma.AccentDirectionLeft,
						},
					},
				},
			},
			Tag:              "span",
			Expected:         `<span class="u l">hello</span>`,
			DirectionClasses: directionClasses,
		},
		{
			Name: "middle down up",
			Word: lemma.Word{
				Hiragana: "hello",
				PitchShapes: []lemma.PitchShape{
					{
						Hiragana: "he",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionDown,
						},
					},
					{
						Hiragana: "llo",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionUp,
							lemma.AccentDirectionLeft,
						},
					},
				},
			},
			Tag:              "span",
			Expected:         `<span class="d">he</span><span class="u l">llo</span>`,
			DirectionClasses: directionClasses,
		},
		{
			Name: "middle up down",
			Word: lemma.Word{
				Hiragana: "hello",
				PitchShapes: []lemma.PitchShape{
					{
						Hiragana: "he",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionUp,
						},
					},
					{
						Hiragana: "llo",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionDown,
							lemma.AccentDirectionLeft,
						},
					},
				},
			},
			Tag:              "span",
			Expected:         `<span class="u">he</span><span class="d l">llo</span>`,
			DirectionClasses: directionClasses,
		},
		{
			Name: "incomplet",
			Word: lemma.Word{
				Hiragana: "hello",
				PitchShapes: []lemma.PitchShape{
					{
						Hiragana: "he",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionUp,
						},
					},
					{
						Hiragana: "llo",
					},
				},
			},
			Tag:              "span",
			Expected:         `<span class="u">he</span><span class="">llo</span>`,
			DirectionClasses: directionClasses,
		},
		{
			Name: "all directions",
			Word: lemma.Word{
				Hiragana: "hello",
				PitchShapes: []lemma.PitchShape{
					{
						Hiragana: "he",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionDown,
						},
					},
					{
						Hiragana: "llo",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionUp,
							lemma.AccentDirectionLeft,
							lemma.AccentDirectionRight,
						},
					},
				},
			},
			Tag:              "span",
			Expected:         `<span class="d">he</span><span class="u l r">llo</span>`,
			DirectionClasses: directionClasses,
		},
		{
			Name: "different tag",
			Word: lemma.Word{
				Hiragana: "hello",
				PitchShapes: []lemma.PitchShape{
					{
						Hiragana: "he",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionDown,
						},
					},
					{
						Hiragana: "llo",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionUp,
							lemma.AccentDirectionLeft,
							lemma.AccentDirectionRight,
						},
					},
				},
			},
			Tag:              "div",
			Expected:         `<div class="d">he</div><div class="u l r">llo</div>`,
			DirectionClasses: directionClasses,
		},
		{
			Name: "different classes",
			Word: lemma.Word{
				Hiragana: "hello",
				PitchShapes: []lemma.PitchShape{
					{
						Hiragana: "he",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionDown,
						},
					},
					{
						Hiragana: "llo",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionUp,
							lemma.AccentDirectionLeft,
							lemma.AccentDirectionRight,
						},
					},
				},
			},
			Tag:              "div",
			Expected:         `<div class="down">he</div><div class="up left right">llo</div>`,
			DirectionClasses: []string{"up", "right", "down", "left"},
		},
		{
			Name: "more classes than needed",
			Word: lemma.Word{
				Hiragana: "hello",
				PitchShapes: []lemma.PitchShape{
					{
						Hiragana: "he",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionDown,
						},
					},
					{
						Hiragana: "llo",
						Directions: []lemma.AccentDirection{
							lemma.AccentDirectionUp,
							lemma.AccentDirectionLeft,
							lemma.AccentDirectionRight,
						},
					},
				},
			},
			Tag:              "div",
			Expected:         `<div class="down">he</div><div class="up left right">llo</div>`,
			DirectionClasses: []string{"up", "right", "down", "left", "extra"},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			actual, err := renderPitch(&tc.Word, tc.Tag, tc.DirectionClasses)
			require.NoError(t, err)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_renderPitch_error(t *testing.T) {
	directionClasses := []string{"u", "r", "d", "l"}
	testCases := []struct {
		Name             string
		Tag              string
		DirectionClasses []string
		ErrorAssert      assert.ErrorAssertionFunc
	}{
		{
			Name:             "no tag",
			DirectionClasses: directionClasses,
			Tag:              "",
			ErrorAssert: func(tt assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(tt, err, "tag")
			},
		},
		{
			Name:             "not enough classes",
			Tag:              "span",
			DirectionClasses: directionClasses[:3],
			ErrorAssert: func(tt assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(tt, err, "direction classes")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			_, err := renderPitch(&lemma.Word{}, tc.Tag, tc.DirectionClasses)
			tc.ErrorAssert(t, err)
		})
	}
}
