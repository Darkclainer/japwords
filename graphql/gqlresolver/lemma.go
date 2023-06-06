package gqlresolver

import (
	"github.com/Darkclainer/japwords/graphql/gqlmodel"
	"github.com/Darkclainer/japwords/pkg/lemma"
)

func convertLemma(src *lemma.Lemma) *gqlmodel.Lemma {
	if src == nil {
		return nil
	}
	var forms []*gqlmodel.Word
	for _, srcForm := range src.Forms {
		forms = append(forms, convertWord(&srcForm))
	}
	var senses []*gqlmodel.Sense
	for _, srcSense := range src.Senses {
		senses = append(senses, &gqlmodel.Sense{
			Definition:   srcSense.Definition,
			PartOfSpeech: srcSense.PartOfSpeech,
			Tags:         srcSense.Tags,
		})
	}
	var audios []*gqlmodel.Audio
	for source, audio := range src.Audio {
		audios = append(audios, &gqlmodel.Audio{
			Type:   source,
			Source: audio,
		})
	}
	return &gqlmodel.Lemma{
		Slug:   convertWord(&src.Slug),
		Tags:   src.Tags,
		Forms:  forms,
		Senses: senses,
		Audio:  audios,
	}
}

func convertWord(src *lemma.Word) *gqlmodel.Word {
	var furigana []*gqlmodel.Furigana
	for _, srcFurigana := range src.Furigana {
		furigana = append(furigana, &gqlmodel.Furigana{
			Kanji:    srcFurigana.Kanji,
			Hiragana: srcFurigana.Hiragana,
		})
	}
	dst := &gqlmodel.Word{
		Word:     src.Word,
		Hiragana: src.Hiragana,
		Furigana: furigana,
		Pitch:    convertPitch(src.PitchShapes()),
	}
	return dst
}

func convertPitch(shaped []lemma.PitchShape) []*gqlmodel.Pitch {
	if len(shaped) == 0 {
		return nil
	}
	result := make([]*gqlmodel.Pitch, len(shaped))
	for i, shape := range shaped {
		newDirections := make([]gqlmodel.PitchType, len(shape.Directions))
		for j, direction := range shape.Directions {
			var newDirection gqlmodel.PitchType
			switch direction {
			case lemma.AccentDirectionUp:
				newDirection = gqlmodel.PitchTypeUp
			case lemma.AccentDirectionRight:
				newDirection = gqlmodel.PitchTypeRight
			case lemma.AccentDirectionDown:
				newDirection = gqlmodel.PitchTypeDown
			case lemma.AccentDirectionLeft:
				newDirection = gqlmodel.PitchTypeLeft
			}
			newDirections[j] = newDirection
		}
		result[i] = &gqlmodel.Pitch{
			Hiragana: shape.Hiragana,
			Pitch:    newDirections,
		}
	}
	return result
}
