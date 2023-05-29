package anki

import "github.com/Darkclainer/japwords/pkg/lemma"

// TODO: Need different pitch representation?
// Lemma is more specific variant of lemma.Lemma.
// This structure include only one meaning.
type Lemma struct {
	Slug          lemma.Word
	SenseIndex    int
	Tags          []string
	Forms         []lemma.Word
	Definitions   []string
	PartsOfSpeech []string
	SenseTags     []string
	Audio         map[string]string
}
