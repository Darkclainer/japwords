package anki

import "github.com/Darkclainer/japwords/pkg/lemma"

// TODO: Need different pitch representation?
// Lemma is more specific variant of lemma.Lemma.
// This structure include only one meaning.
type Lemma struct {
	Slug          Word
	SenseIndex    int
	Tags          []string
	Forms         []Word
	Definitions   []string
	PartsOfSpeech []string
	SenseTags     []string
	Audio         map[string]string
}

// Word is simplified version of lemma.Word,
// difference is in Pitch representation.
type Word struct {
	Word     string
	Hiragana string
	Furigana lemma.Furigana
	Pitches  []lemma.PitchShape
}
