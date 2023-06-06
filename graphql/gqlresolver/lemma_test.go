package gqlresolver

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Darkclainer/japwords/graphql/gqlmodel"
	"github.com/Darkclainer/japwords/pkg/lemma"
)

func Test_convertLemma(t *testing.T) {
	testCases := []struct {
		Name     string
		Lemma    *lemma.Lemma
		Expected *gqlmodel.Lemma
	}{
		{
			Name: "Empty",
		},
		{
			Name: "everything",
			Lemma: &lemma.Lemma{
				Slug: lemma.Word{
					Word:     "hello",
					Hiragana: "world",
					Furigana: []lemma.FuriganaChar{
						{
							Kanji:    "a",
							Hiragana: "b",
						},
						{
							Kanji:    "c",
							Hiragana: "d",
						},
					},
					Pitches: []lemma.Pitch{
						{
							Position: 1,
							IsHigh:   false,
						},
						{
							Position: 5,
							IsHigh:   true,
						},
					},
				},
				Tags: []string{
					"taga", "tagb",
				},
				Forms: []lemma.Word{
					{
						Word:     "greetings",
						Hiragana: "world",
						Furigana: []lemma.FuriganaChar{
							{
								Kanji:    "d",
								Hiragana: "e",
							},
						},
						Pitches: []lemma.Pitch{
							{
								Position: 2,
								IsHigh:   true,
							},
							{
								Position: 5,
								IsHigh:   false,
							},
						},
					},
					{
						Word:     "simple",
						Hiragana: "word",
					},
				},
				Senses: []lemma.WordSense{
					{
						Definition: []string{
							"def1",
							"def2",
							"def3",
						},
						PartOfSpeech: []string{
							"pos1", "pos2",
						},
						Tags: []string{
							"tag1", "tag2",
						},
					},
				},
				Audio: map[string]string{
					"k1": "v1",
					"k2": "v2",
				},
			},
			Expected: &gqlmodel.Lemma{
				Slug: &gqlmodel.Word{
					Word:     "hello",
					Hiragana: "world",
					Furigana: []*gqlmodel.Furigana{
						{
							Kanji:    "a",
							Hiragana: "b",
						},
						{
							Kanji:    "c",
							Hiragana: "d",
						},
					},
					Pitch: []*gqlmodel.Pitch{
						{
							Hiragana: "w",
							Pitch: []gqlmodel.PitchType{
								gqlmodel.PitchTypeDown,
							},
						},
						{
							Hiragana: "orld",
							Pitch: []gqlmodel.PitchType{
								gqlmodel.PitchTypeUp,
								gqlmodel.PitchTypeLeft,
							},
						},
					},
				},
				Tags: []string{
					"taga", "tagb",
				},
				Forms: []*gqlmodel.Word{
					{
						Word:     "greetings",
						Hiragana: "world",
						Furigana: []*gqlmodel.Furigana{
							{
								Kanji:    "d",
								Hiragana: "e",
							},
						},
						Pitch: []*gqlmodel.Pitch{
							{
								Hiragana: "wo",
								Pitch: []gqlmodel.PitchType{
									gqlmodel.PitchTypeUp,
								},
							},
							{
								Hiragana: "rld",
								Pitch: []gqlmodel.PitchType{
									gqlmodel.PitchTypeDown,
									gqlmodel.PitchTypeLeft,
								},
							},
						},
					},
					{
						Word:     "simple",
						Hiragana: "word",
						Pitch: []*gqlmodel.Pitch{
							{
								Hiragana: "word",
								Pitch:    []gqlmodel.PitchType{},
							},
						},
					},
				},
				Senses: []*gqlmodel.Sense{
					{
						Definition: []string{
							"def1",
							"def2",
							"def3",
						},
						PartOfSpeech: []string{
							"pos1", "pos2",
						},
						Tags: []string{
							"tag1", "tag2",
						},
					},
				},
				Audio: []*gqlmodel.Audio{
					{
						Type:   "k1",
						Source: "v1",
					},
					{
						Type:   "k2",
						Source: "v2",
					},
				},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			result := convertLemma(tc.Lemma)
			if result != nil {
				// as we convert audio from map, elements are in undefined order
				sort.Slice(result.Audio, func(i, j int) bool {
					return result.Audio[i].Type < result.Audio[j].Type
				})
			}
			assert.Equal(t, tc.Expected, result)
		})
	}
}
