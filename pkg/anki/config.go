package anki

import (
	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
	"github.com/Darkclainer/japwords/pkg/config"
)

type Config struct {
	Addr   string
	APIKey string

	Deck     string
	NoteType string

	Mapping TemplateMapping
}

func (c *Config) Equal(o any) bool {
	oc, ok := o.(*Config)
	if !ok {
		return false
	}
	// both nil or same object
	if c == oc {
		return true
	}
	if c == nil || oc == nil {
		return false
	}
	scalarEq := c.Addr == oc.Addr &&
		c.APIKey == oc.APIKey &&
		c.Deck == oc.Deck &&
		c.NoteType == oc.NoteType
	if !scalarEq {
		return false
	}
	return c.Mapping.Equal(oc.Mapping)
}

func (c *Config) options() *ankiconnect.Options {
	return &ankiconnect.Options{
		URL:    c.Addr,
		APIKey: c.APIKey,
	}
}

// ConfigReloader allows to change anki part of user config.
type ConfigReloader struct {
	anki           *Anki
	updateConfigFn config.UpdateConfigFunc
}

func NewConfigReloader(anki *Anki, configManager *config.Manager) (*ConfigReloader, error) {
	reloader := &ConfigReloader{
		anki: anki,
	}
	_, updateFn, err := configManager.Register(reloader)
	if err != nil {
		return nil, err
	}
	reloader.updateConfigFn = updateFn
	return reloader, nil
}

// Config is implementation if config.Cosumer interface
func (cr *ConfigReloader) Config(uc *config.UserConfig) (config.Part, error) {
	// TODO: we can do check here
	return &Config{
		Addr:     uc.Anki.Addr,
		APIKey:   uc.Anki.APIKey,
		Deck:     uc.Anki.Deck,
		NoteType: uc.Anki.NoteType,
	}, nil
}

// Reload is implementation if config.Reloader interface
func (cr *ConfigReloader) Reload(o config.Part) error {
	oc, ok := o.(*Config)
	if !ok {
		panic("unreachable")
	}
	return cr.anki.ReloadConfig(oc)
}

func (cr *ConfigReloader) UpdateConnection(addr string, apikey string) error {
	// TODO: check for addr?
	return cr.updateConfigFn(func(uc *config.UserConfig) error {
		uc.Anki.Addr = addr
		uc.Anki.APIKey = apikey
		return nil
	})
}

func (cr *ConfigReloader) UpdateDeck(name string) error {
	return cr.updateConfigFn(func(uc *config.UserConfig) error {
		uc.Anki.Deck = name
		return nil
	})
}

func (cr *ConfigReloader) UpdateNoteType(name string) error {
	return cr.updateConfigFn(func(uc *config.UserConfig) error {
		uc.Anki.NoteType = name
		return nil
	})
}

func (cr *ConfigReloader) UpdateMapping(mapping map[string]string) error {
	return cr.updateConfigFn(func(uc *config.UserConfig) error {
		uc.Anki.FieldMapping = mapping
		return nil
	})
}
