package gqlresolver

import "github.com/Darkclainer/japwords/pkg/lemma"

func expandLemmas(lemmas []*lemma.Lemma) []*lemma.ProjectedLemma {
	var projectedLemmas []*lemma.ProjectedLemma
	for _, l := range lemmas {
		for _, wordSense := range l.Senses {
			projectedLemmas = append(projectedLemmas, &lemma.ProjectedLemma{
				Slug:          l.Slug,
				Tags:          l.Tags,
				Forms:         l.Forms,
				Definitions:   wordSense.Definition,
				PartsOfSpeech: wordSense.PartOfSpeech,
				SenseTags:     wordSense.Tags,
				Audio:         l.Audio,
			})
		}
	}
	return projectedLemmas
}
