package lemma

type slugReading struct {
	Slug    string
	Reading string
}

// Enrich add pitch infromation from pitches to lemmas. Modifies lemmas, if you need: make copy.
func Enrich(lemmas []*Lemma, pitchedLemmas []*PitchedLemma) {
	wordMap := map[slugReading][]*Word{}
	for _, lemma := range lemmas {
		addLemmaWords(wordMap, lemma)
	}
	for i := len(pitchedLemmas) - 1; i >= 0; i-- {
		pitched := pitchedLemmas[i]
		key := slugReading{
			Slug:    pitched.Slug,
			Reading: pitched.Hiragana,
		}
		words := wordMap[key]
		for _, word := range words {
			word.PitchShapes = pitched.PitchShapes
		}
	}
}

func addLemmaWords(dst map[slugReading][]*Word, lemma *Lemma) {
	addWord(dst, &lemma.Slug)
	for i := range lemma.Forms {
		addWord(dst, &lemma.Forms[i])
	}
}

func addWord(dst map[slugReading][]*Word, word *Word) {
	key := slugReading{
		Slug:    word.Word,
		Reading: word.Hiragana,
	}
	dst[key] = append(dst[key], word)
}
