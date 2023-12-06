package gqlresolver

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Darkclainer/japwords/pkg/lemma"
)

func Test_expandLemmas(t *testing.T) {
	word1 := lemma.Word{
		Word:     "foo",
		Hiragana: "foo",
		Furigana: []lemma.FuriganaChar{
			{},
		},
		PitchShapes: []lemma.PitchShape{
			{},
		},
	}
	word2 := lemma.Word{
		Word:     "bar",
		Hiragana: "bar",
		Furigana: []lemma.FuriganaChar{
			{}, {},
		},
		PitchShapes: []lemma.PitchShape{
			{}, {},
		},
	}
	lemmas := []*lemma.Lemma{
		{
			Slug:  word1,
			Tags:  []string{"first"},
			Forms: []lemma.Word{word1},
			Senses: []lemma.WordSense{
				{
					Definition:   []string{"a", "b"},
					PartOfSpeech: []string{"pos"},
					Tags:         []string{"sensetag"},
				},
				{
					Definition:   []string{"c", "d"},
					PartOfSpeech: []string{"sop"},
				},
			},
			Audio: []lemma.Audio{
				{
					Type:   "hello",
					Source: "world",
				},
			},
		},
		{
			Slug:  word2,
			Tags:  []string{"second"},
			Forms: []lemma.Word{word1, word2},
			Senses: []lemma.WordSense{
				{
					Definition:   []string{"second"},
					PartOfSpeech: []string{"pos2"},
				},
			},
		},
	}
	expected := []*lemma.ProjectedLemma{
		{
			Slug:          word1,
			Tags:          []string{"first"},
			Forms:         []lemma.Word{word1},
			Definitions:   []string{"a", "b"},
			PartsOfSpeech: []string{"pos"},
			SenseTags:     []string{"sensetag"},
			Audio: []lemma.Audio{
				{
					Type:   "hello",
					Source: "world",
				},
			},
		},
		{
			Slug:          word1,
			Tags:          []string{"first"},
			Forms:         []lemma.Word{word1},
			Definitions:   []string{"c", "d"},
			PartsOfSpeech: []string{"sop"},
			Audio: []lemma.Audio{
				{
					Type:   "hello",
					Source: "world",
				},
			},
		},
		{
			Slug:          word2,
			Tags:          []string{"second"},
			Forms:         []lemma.Word{word1, word2},
			Definitions:   []string{"second"},
			PartsOfSpeech: []string{"pos2"},
		},
	}
	actual := expandLemmas(lemmas)
	assert.Equal(t, expected, actual)
}
