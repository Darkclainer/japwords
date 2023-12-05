package gqlresolver

import "github.com/Darkclainer/japwords/pkg/lemma"

func expandLemmas(lemmas []*lemma.Lemma) []*lemma.ProjectedLemma {
	var projectedLemmas []*lemma.ProjectedLemma
	for _, l := range lemmas {
		for i, wordSense := range l.Senses {
			projectedLemmas = append(projectedLemmas, &lemma.ProjectedLemma{
				Slug:          expandWord(l.Slug),
				SenseIndex:    i,
				Tags:          l.Tags,
				Forms:         expandWords(l.Forms),
				Definitions:   wordSense.Definition,
				PartsOfSpeech: wordSense.PartOfSpeech,
				SenseTags:     wordSense.Tags,
				Audio:         l.Audio,
			})
		}
	}
	return projectedLemmas
}

func expandWord(word lemma.Word) lemma.ProjectedWord {
	return lemma.ProjectedWord{
		Word:     word.Word,
		Hiragana: word.Hiragana,
		Furigana: word.Furigana,
		Pitches:  word.PitchShapes(),
	}
}

func expandWords(words []lemma.Word) []lemma.ProjectedWord {
	projectedWords := make([]lemma.ProjectedWord, len(words))
	for i := range words {
		projectedWords[i] = expandWord(words[i])
	}
	return projectedWords
}
