package lemma

// PitchedLemma is data that we extract from wadoku dictionary.
type PitchedLemma struct {
	Slug     string
	Hiragana string
	Pitches  []Pitch
}
