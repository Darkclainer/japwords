package config

import (
	"github.com/huandu/go-clone/generic"
)

type UserConfig struct {
	Addr       string     `yaml:"addr" koanf:"addr"`
	Anki       Anki       `yaml:"anki" koanf:"anki"`
	Dictionary Dictionary `yaml:"dictionary" koanf:"dictionary"`
}

type Anki struct {
	// Addr is host:port of AnkiConnect address.
	// Protocol must not be specified, anyway https seems to be redudant.
	Addr string `yaml:"addr" koanf:"addr"`
	// APIKey is secret that can be enabled on AnkiConnect.
	// Can be any string. Empty string means that no secret will be used.
	APIKey string `yaml:"api-key" koanf:"api-key"`

	// Deck is the name of deck in Anki where new cards will be added.
	// Deck can be any non-empty string, that doesn't contain `"` (Anki removes it),
	// that doesn't start or end with spaces.
	// I have found that Anki can convert some names with `::`, for example name
	// `hello::::world` will have alias `hello::blank::world`, but this should not
	// break anything.
	Deck string `yaml:"deck" koanf:"deck"`
	// NoteType is the name of note type that will be used for creation new notes in Anki.
	// NoteType can be any non-empty string, that doesn't contain `"` (Anki removes it) and
	// doesn't start or end with spaces.
	NoteType string `yaml:"note-type" koanf:"note-type"`

	// FieldMapping specifies how note fields will be filled in new notes.
	//
	// Key is the name of field in Anki. Can be any non-empty string that doesn't contain
	// symbols `:`, `"`, `{` or `}` and doesn't start or end with spaces.
	//
	// Value should be valid go text/template. For more details see pkg/anki/template.go.
	FieldMapping map[string]string `yaml:"fields" koanf:"fields"`
}

type Dictionary struct {
	Workers   int               `yaml:"workers" koanf:"workers"`
	UserAgent string            `yaml:"user-agent" koanf:"user-agent"`
	Headers   map[string]string `yaml:"headers" koanf:"headers"`
	Jisho     Jisho             `yaml:"jisho" koanf:"jisho"`
	Wadoku    Wadoku            `yaml:"wadoku" koanf:"wadoku"`
}

type Jisho struct {
	URL string `yaml:"url" koanf:"url"`
}

type Wadoku struct {
	URL string `yaml:"url" koanf:"url"`
}

func DefaultUserConfig() *UserConfig {
	return &UserConfig{
		Addr: "",
		Anki: Anki{
			Addr:         "",
			APIKey:       "",
			Deck:         "",
			NoteType:     "",
			FieldMapping: map[string]string{},
		},
		Dictionary: Dictionary{
			Workers:   0,
			UserAgent: "",
			Headers:   map[string]string{},
			Jisho: Jisho{
				URL: "",
			},
			Wadoku: Wadoku{
				URL: "",
			},
		},
	}
}

func (uc *UserConfig) Clone() *UserConfig {
	// don't want to write tests, so I will use third-party package
	return clone.Clone(uc)
}
