package multidict

type Options struct {
	Workers int

	LemmaDict LemmaDict
	PitchDict PitchDict
}
