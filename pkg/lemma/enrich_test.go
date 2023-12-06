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
					PitchShapes: []PitchShape{
						{
							Hiragana: "world",
							Directions: []AccentDirection{
								AccentDirectionDown,
							},
						},
					},
				},
			},
			Expected: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
						PitchShapes: []PitchShape{
							{
								Hiragana: "world",
								Directions: []AccentDirection{
									AccentDirectionDown,
								},
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
					PitchShapes: []PitchShape{
						{
							Hiragana: "world",
							Directions: []AccentDirection{
								AccentDirectionDown,
							},
						},
					},
				},
			},
			Expected: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
						PitchShapes: []PitchShape{
							{
								Hiragana: "world",
								Directions: []AccentDirection{
									AccentDirectionDown,
								},
							},
						},
					},
				},
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
						PitchShapes: []PitchShape{
							{
								Hiragana: "world",
								Directions: []AccentDirection{
									AccentDirectionDown,
								},
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
					PitchShapes: []PitchShape{
						{
							Hiragana: "world",
							Directions: []AccentDirection{
								AccentDirectionDown,
							},
						},
					},
				},
				{
					Slug:     "hello",
					Hiragana: "world",
					PitchShapes: []PitchShape{
						{
							Hiragana: "world",
							Directions: []AccentDirection{
								AccentDirectionUp,
							},
						},
					},
				},
			},
			Expected: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
						PitchShapes: []PitchShape{
							{
								Hiragana: "world",
								Directions: []AccentDirection{
									AccentDirectionDown,
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "readings from form",
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
					PitchShapes: []PitchShape{
						{
							Hiragana: "world",
							Directions: []AccentDirection{
								AccentDirectionDown,
							},
						},
					},
				},
				{
					Slug:     "greating",
					Hiragana: "world",
					PitchShapes: []PitchShape{
						{
							Hiragana: "world",
							Directions: []AccentDirection{
								AccentDirectionUp,
							},
						},
					},
				},
			},
			Expected: []*Lemma{
				{
					Slug: Word{
						Word:     "hello",
						Hiragana: "world",
						PitchShapes: []PitchShape{
							{
								Hiragana: "world",
								Directions: []AccentDirection{
									AccentDirectionDown,
								},
							},
						},
					},
					Forms: []Word{
						{
							Word:     "greating",
							Hiragana: "world",
							PitchShapes: []PitchShape{
								{
									Hiragana: "world",
									Directions: []AccentDirection{
										AccentDirectionUp,
									},
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
