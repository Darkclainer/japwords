package anki

import (
	"html/template"
	"net/http"
	"net/http/httptest"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
	"github.com/Darkclainer/japwords/pkg/config"
)

type Config struct {
	Addr   string
	APIKey string

	Deck     string
	NoteType string

	Mapping map[string]*Template
}

type Template struct {
	Src  string
	Tmpl *template.Template
}

func (c *Config) Equal(o any) bool {
	oc, ok := o.(*Config)
	if !ok {
		return false
	}
	simpleEqual := c.Addr == oc.Addr &&
		c.APIKey == oc.APIKey &&
		c.Deck == oc.Deck &&
		c.NoteType == oc.NoteType
	if !simpleEqual {
		return false
	}
	if len(c.Mapping) != len(oc.Mapping) {
		return false
	}
	for field, tmpl := range c.Mapping {
		otherTmpl, ok := oc.Mapping[field]
		if !ok || otherTmpl.Src != tmpl.Src {
			return false
		}
	}
	return true
}

func (c *Config) options() *ankiconnect.Options {
	return &ankiconnect.Options{
		URL:    c.Addr,
		APIKey: c.APIKey,
	}
}

// ConfigReloader allows to change anki part of user config.
type ConfigReloader struct {
	anki   *Anki
	config *Config
	ucm    *config.Manager
}

// Config is implementation if config.Cosumer interface
func (cr *ConfigReloader) Config(uc *config.UserConfig) (*Config, error) {
	return &Config{
		Addr:     uc.Anki.Addr,
		APIKey:   uc.Anki.APIKey,
		Deck:     uc.Anki.Deck,
		NoteType: uc.Anki.NoteType,
	}, nil
}

// Reload is implementation if config.Reloader interface
func (cr *ConfigReloader) Reload(o any) error {
	oc, ok := o.(*Config)
	if !ok {
		panic("unreachable")
	}
	return cr.anki.ReloadConfig(oc)
}

func (cr *ConfigReloader) UpdateConnection(addr string, apikey string) error {
	return cr.ucm.UpdateConfig(func(uc *config.UserConfig) error {
		uc.Anki.Addr = addr
		uc.Anki.APIKey = apikey
		return nil
	})
}

func (cr *ConfigReloader) UpdateDeck(name string) error {
	return cr.ucm.UpdateConfig(func(uc *config.UserConfig) error {
		uc.Anki.Deck = name
		return nil
	})
}
