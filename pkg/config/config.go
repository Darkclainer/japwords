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
	Addr   string `koanf:"addr"`
	APIKey string `koanf:"api-key"`

	Deck     string `koanf:"deck"`
	NoteType string `koanf:"note-type"`

	FieldMapping map[string]string `koanf:"fields"`
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
		Anki: Anki{},
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
