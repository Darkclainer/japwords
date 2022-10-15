package jisho

type Lemma struct {
	Slug  Word
	JLPT  []string
	Forms []Word
	Sense WordSense
	// Audio is array of links to audio files.
	// Key is format
	Audio map[string]string
}

type Word struct {
	Word     string
	Reading  string
	Furigana Furigana
}

type Furigana []FuriganaChar

type FuriganaChar struct {
	Kanji    string
	Hiragana string
}

type WordSense struct {
	// Definition is slice of synonymous definitions in english
	Definition   []string
	PartOfSpeech []string
	Tags         []string
}
