package lemma

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Word_PitchShapes(t *testing.T) {
	testCases := []struct {
		Name     string
		Hiragana string
		Pitches  []Pitch
		Expected []PitchShape
	}{
		{
			Name:     "no pitches",
			Hiragana: "hello",
			Expected: []PitchShape{
				{
					Hiragana: "hello",
				},
			},
		},
		{
			Name:     "up only",
			Hiragana: "hello",
			Pitches: []Pitch{
				{
					Position: 5,
					IsHigh:   true,
				},
			},
			Expected: []PitchShape{
				{
					Hiragana: "hello",
					Directions: []AccentDirection{
						AccentDirectionUp,
					},
				},
			},
		},
		{
			Name:     "down only",
			Hiragana: "hello",
			Pitches: []Pitch{
				{
					Position: 5,
					IsHigh:   false,
				},
			},
			Expected: []PitchShape{
				{
					Hiragana: "hello",
					Directions: []AccentDirection{
						AccentDirectionDown,
					},
				},
			},
		},
		{
			Name:     "up down",
			Hiragana: "hello",
			Pitches: []Pitch{
				{
					Position: 1,
					IsHigh:   true,
				},
				{
					Position: 5,
					IsHigh:   false,
				},
			},
			Expected: []PitchShape{
				{
					Hiragana: "h",
					Directions: []AccentDirection{
						AccentDirectionUp,
					},
				},
				{
					Hiragana: "ello",
					Directions: []AccentDirection{
						AccentDirectionDown,
						AccentDirectionLeft,
					},
				},
			},
		},
		{
			Name:     "down up",
			Hiragana: "hello",
			Pitches: []Pitch{
				{
					Position: 1,
					IsHigh:   false,
				},
				{
					Position: 5,
					IsHigh:   true,
				},
			},
			Expected: []PitchShape{
				{
					Hiragana: "h",
					Directions: []AccentDirection{
						AccentDirectionDown,
					},
				},
				{
					Hiragana: "ello",
					Directions: []AccentDirection{
						AccentDirectionUp,
						AccentDirectionLeft,
					},
				},
			},
		},
		{
			Name:     "up down tail",
			Hiragana: "hello",
			Pitches: []Pitch{
				{
					Position: 5,
					IsHigh:   true,
				},
				{
					Position: 5,
					IsHigh:   false,
				},
			},
			Expected: []PitchShape{
				{
					Hiragana: "hello",
					Directions: []AccentDirection{
						AccentDirectionUp,
						AccentDirectionRight,
					},
				},
			},
		},
		{
			Name:     "down up head",
			Hiragana: "hello",
			Pitches: []Pitch{
				{
					Position: 0,
					IsHigh:   false,
				},
				{
					Position: 5,
					IsHigh:   true,
				},
			},
			Expected: []PitchShape{
				{
					Hiragana: "hello",
					Directions: []AccentDirection{
						AccentDirectionUp,
						AccentDirectionLeft,
					},
				},
			},
		},
		{
			Name:     "incomplete pitch data",
			Hiragana: "hello",
			Pitches: []Pitch{
				{
					Position: 1,
					IsHigh:   true,
				},
			},
			Expected: []PitchShape{
				{
					Hiragana: "h",
					Directions: []AccentDirection{
						AccentDirectionUp,
					},
				},
				{
					Hiragana: "ello",
				},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			word := Word{
				Hiragana: tc.Hiragana,
				Pitches:  tc.Pitches,
			}
			result := word.PitchShapes()
			assert.Equal(t, tc.Expected, result)
		})
	}
}
