package config

type UserConfig struct {
	Addr       string
	Anki       Anki
	Dictionary Dictionary
}

type Anki struct{}

type Dictionary struct {
	Workers   int
	UserAgent string
	Headers   map[string]string
	Jisho     Jisho
	Wadoku    Wadoku
}

type Jisho struct {
	URL string
}

type Wadoku struct {
	URL string
}
