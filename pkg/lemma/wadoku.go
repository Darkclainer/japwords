package lemma

// WadokuLemma is data that we extract from wadoku dictionary.
type WadokuLemma struct {
	Slug     string
	Hiragana string
	Pitches  []Pitch
}
