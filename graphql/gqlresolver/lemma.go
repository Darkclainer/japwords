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
		Pitch:    convertPitch(src.Hiragana, src.Pitches),
	}
	return dst
}

func convertPitch(hiragana string, pitches []lemma.Pitch) []*gqlmodel.Pitch {
	var gqlPitches []*gqlmodel.Pitch
	lastPosition := 0
	for i := 0; i < len(pitches)-1; i++ {
		gqlPitches = append(gqlPitches, convertAdjancentPitch(hiragana, lastPosition, pitches[i], pitches[i+1]))
		lastPosition = pitches[i].Position
	}
	if len(pitches) == 0 {
		return nil
	}
	gqlPitches = append(gqlPitches, convertLastPitch(hiragana, lastPosition, pitches[len(pitches)-1]))
	return gqlPitches
}

func convertAdjancentPitch(hiragana string, last int, left, right lemma.Pitch) *gqlmodel.Pitch {
	pitch := []gqlmodel.PitchType{convertBasePitch(left.IsHigh)}
	if left.IsHigh != right.IsHigh {
		pitch = append(pitch, gqlmodel.PitchTypeRight)
	}
	return &gqlmodel.Pitch{
		Hiragana: hiragana[last:left.Position],
		Pitch:    pitch,
	}
}

func convertLastPitch(hiragana string, last int, pitch lemma.Pitch) *gqlmodel.Pitch {
	return &gqlmodel.Pitch{
		Hiragana: hiragana[last:pitch.Position],
		Pitch:    []gqlmodel.PitchType{convertBasePitch(pitch.IsHigh)},
	}
}

func convertBasePitch(isHigh bool) gqlmodel.PitchType {
	if isHigh {
		return gqlmodel.PitchTypeUp
	}
	return gqlmodel.PitchTypeDown
}
