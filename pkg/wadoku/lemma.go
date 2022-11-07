package wadoku

type Lemma struct {
	Slug    string
	Reading Reading
}

type Reading struct {
	Hiragana string
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
	// 	{ Position: 6, IsHigh: true  },
	// 	{ Position: 6, IsHigh: false },
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
