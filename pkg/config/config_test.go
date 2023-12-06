package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Darkclainer/japwords/pkg/anki"
	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/lemma"
)

// Test_Anki_ConfigReloader_DefaultConfig tests that default configuration for anki is valid
func Test_Anki_ConfigReloader_DefaultConfig(t *testing.T) {
	reloader := &anki.ConfigReloader{}
	_, err := reloader.Config(config.DefaultUserConfig())
	require.NoError(t, err)
}

// Test_DefaultUserConfig_AnkiFieldMapping tests that default mapping is what we want
func Test_DefaultUserConfig_AnkiFieldMapping(t *testing.T) {
	mapping := config.DefaultUserConfig().Anki.FieldMapping
	testCases := []struct {
		Name     string
		Lemma    lemma.ProjectedLemma
		Expected map[string]string
	}{
		{
			Name: "Sort full",
			Lemma: lemma.ProjectedLemma{
				Slug: lemma.Word{
					Word: "Hello",
				},
				Definitions: []string{"world hello", "there"},
			},
			Expected: map[string]string{
				"Sort": "Hello-world_hell",
			},
		},
		{
			Name: "Sort min len",
			Lemma: lemma.ProjectedLemma{
				Slug: lemma.Word{
					Word: "Hello",
				},
				Definitions: []string{"world"},
			},
			Expected: map[string]string{
				"Sort": "Hello-world",
			},
		},
		{
			Name: "Sort no definitions",
			Lemma: lemma.ProjectedLemma{
				Slug: lemma.Word{
					Word: "Hello",
				},
			},
			Expected: map[string]string{
				"Sort": "Hello-",
			},
		},
		{
			Name: "Kanji",
			Lemma: lemma.ProjectedLemma{
				Slug: lemma.Word{
					Word: "Hello",
				},
			},
			Expected: map[string]string{
				"Kanji": "Hello",
			},
		},
		{
			Name: "Furigana",
			Lemma: lemma.ProjectedLemma{
				Slug: lemma.Word{
					Furigana: lemma.Furigana{
						{
							Kanji:    "he",
							Hiragana: "12",
						},
						{
							Kanji: "l",
						},
						{
							Kanji:    "lo",
							Hiragana: "3",
						},
					},
				},
			},
			Expected: map[string]string{
				"Furigana": "he[12]l lo[3]",
			},
		},
		{
			Name: "Kana",
			Lemma: lemma.ProjectedLemma{
				Slug: lemma.Word{
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
								lemma.AccentDirectionRight,
							},
						},
					},
				},
			},
			Expected: map[string]string{
				"Kana": `<span class="border-u">he</span><span class="border-d border-l border-r">llo</span>`,
			},
		},
		{
			Name: "PoS one",
			Lemma: lemma.ProjectedLemma{
				PartsOfSpeech: []string{
					"first",
				},
			},
			Expected: map[string]string{
				"PoS": `<span class="pos">first</span>`,
			},
		},
		{
			Name: "PoS two",
			Lemma: lemma.ProjectedLemma{
				PartsOfSpeech: []string{
					"first", "second, something",
				},
			},
			Expected: map[string]string{
				"PoS": `<span class="pos">first</span> <span class="pos">second, something</span>`,
			},
		},
		{
			Name: "English",
			Lemma: lemma.ProjectedLemma{
				Definitions: []string{
					"one", "two", "three",
				},
			},
			Expected: map[string]string{
				"English": `<span>one</span> <span>two</span> <span>three</span>`,
			},
		},
		{
			Name: "SenseTags",
			Lemma: lemma.ProjectedLemma{
				SenseTags: []string{
					"one", "two", "three",
				},
			},
			Expected: map[string]string{
				"SenseTags": `<span class="sensetag">one</span> <span class="sensetag">two</span> <span class="sensetag">three</span>`,
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			actual := map[string]string{}
			for name, template := range mapping {
				result, err := anki.RenderRawTemplate(template, &tc.Lemma)
				require.NoError(t, err)
				// filter out annoying values
				if _, ok := tc.Expected[name]; ok {
					actual[name] = result
				}
			}
			assert.Equal(t, tc.Expected, actual)
		})
	}
}
