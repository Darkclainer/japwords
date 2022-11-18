package lemma

type slugReading struct {
	Slug    string
	Reading string
}

// Enrich add infom from reading to lemmas. Modifies lemmas, if you need: make copy.
func Enrich(lemmas []*Lemma, readings []*WadokuLemma) {
	wordMap := map[slugReading][]*Word{}
	for _, lemma := range lemmas {
		addLemmaWords(wordMap, lemma)
	}
	for i := len(readings) - 1; i >= 0; i-- {
		reading := readings[i]
		key := slugReading{
			Slug:    reading.Slug,
			Reading: reading.Hiragana,
		}
		words := wordMap[key]
		for _, word := range words {
			word.Pitches = reading.Pitches
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
