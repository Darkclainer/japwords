package lemma

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Enrich(t *testing.T) {
	testCases := []struct {
		Name         string
		Lemmas       []*Lemma
		PitchedLemma []*PitchedLemma
		Expected     []*Lemma
	}{
		{
			Name: "empty",
		},
		{
			Name: "without readings",
			Lemmas: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
					},
				},
			},
			Expected: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
					},
				},
			},
		},
		{
			Name: "single match",
			Lemmas: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
					},
				},
			},
			PitchedLemma: []*PitchedLemma{
				{
					Slug:     "hello",
					Hiragana: "world",
					Pitches: []Pitch{
						{
							Position: 1,
							IsHigh:   false,
						},
					},
				},
			},
			Expected: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 1,
								IsHigh:   false,
							},
						},
					},
				},
			},
		},
		{
			Name: "duplicated lemma",
			Lemmas: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
					},
				},
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
					},
				},
			},
			PitchedLemma: []*PitchedLemma{
				{
					Slug:     "hello",
					Hiragana: "world",
					Pitches: []Pitch{
						{
							Position: 1,
							IsHigh:   false,
						},
					},
				},
			},
			Expected: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 1,
								IsHigh:   false,
							},
						},
					},
				},
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 1,
								IsHigh:   false,
							},
						},
					},
				},
			},
		},
		{
			Name: "duplicated reading priority",
			Lemmas: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
					},
				},
			},
			PitchedLemma: []*PitchedLemma{
				{
					Slug:     "hello",
					Hiragana: "world",
					Pitches: []Pitch{
						{
							Position: 1,
							IsHigh:   false,
						},
					},
				},
				{
					Slug:     "hello",
					Hiragana: "world",
					Pitches: []Pitch{
						{
							Position: 1,
							IsHigh:   true,
						},
					},
				},
			},
			Expected: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 1,
								IsHigh:   false,
							},
						},
					},
				},
			},
		},
		{
			Name: "readings form form",
			Lemmas: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
					},
					Forms: []Word{
						{
							Word:     "greating",
							Hiragana: "world",
						},
					},
				},
			},
			PitchedLemma: []*PitchedLemma{
				{
					Slug:     "hello",
					Hiragana: "world",
					Pitches: []Pitch{
						{
							Position: 1,
						},
					},
				},
				{
					Slug:     "greating",
					Hiragana: "world",
					Pitches: []Pitch{
						{
							Position: 2,
						},
					},
				},
			},
			Expected: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 1,
							},
						},
					},
					Forms: []Word{
						{
							Word:     "greating",
							Hiragana: "world",
							Pitches: []Pitch{
								{
									Position: 2,
								},
							},
						},
					},
				},
			},
		},
		// TODO: check that duplicated lemmas both filled, check case duplicated reading
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			Enrich(tc.Lemmas, tc.PitchedLemma)
			assert.Equal(t, tc.Expected, tc.Lemmas)
		})
	}
}
