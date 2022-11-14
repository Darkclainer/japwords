// lemma is like model package. It contains structure for jisho and wadoku
// packages and also can compose results from them.
package lemma

type Lemma struct {
	Slug   Word
	Tags   []string
	Forms  []Word
	Senses []WordSense
	// Audio is array of links to audio files.
	// Key is format
	Audio map[string]string
}

type Word struct {
	Word     string
	Hiragana string
	Furigana Furigana
	// Pitches are encoded japanese pitch accent.
	// Every element of Pitches describe what accent (high or low)
	// should be used from previous element of Pitches up to and including
	// current specified position. Last elemnt can have virtual position
	// past length of Hiragana string, in that case it indicates
	// pitch accent of particle that follow word.
	//
	// For example consider world: 紙 「かみ」
	// The pitch accent of this word should be encoded next way:
	//
	// []Pitch{
	// 	{ Position: 3, IsHigh: false },
	// 	{ Position: 6, IsHigh: true  }
	// 	{ Position: 9, IsHigh: false },
	// }
	//
	// That means that 紙が will read as:
	// か (low) み (high) が (low)
	//
	// Note: Tokyo dialect can be encoded in much easier way, but I left
	// flexibility in case dictionary contain some entries that doesn't
	// follow Tokyo dialect rules.
	Pitches []Pitch
}

type Pitch struct {
	Position int
	IsHigh   bool
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
