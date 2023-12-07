package anki

import (
	"errors"
	"fmt"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
	"github.com/Darkclainer/japwords/pkg/config"
)

type Config struct {
	Addr   string
	APIKey string

	Deck     string
	NoteType string

	AudioField         string
	AudioPreferredType string

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
		c.NoteType == oc.NoteType &&
		c.AudioField == oc.AudioField &&
		c.AudioPreferredType == oc.AudioPreferredType
	if !scalarEq {
		return false
	}
	return c.Mapping.Equal(oc.Mapping)
}

func (c *Config) options() *ankiconnect.Options {
	return &ankiconnect.Options{
		URL:    ankiAddrToURL(c.Addr),
		APIKey: c.APIKey,
	}
}

// ankiAddrToURL convert addr to url. I prefer to use addr, because
// ankiconnect anyway used with http.
func ankiAddrToURL(addr string) string {
	return "http://" + addr
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

// Config is implementation of config.Consumer interface.
// It returns errors, but they are not supposed to be examined.
func (cr *ConfigReloader) Config(uc *config.UserConfig) (config.Part, error) {
	conf := uc.Anki
	var errs []error
	err := validateAddr(conf.Addr)
	if err != nil {
		errs = append(errs, fmt.Errorf("anki config Addr validation failed: %w", err))
	}
	err = validateDeckName(conf.Deck)
	if err != nil {
		errs = append(errs, fmt.Errorf("anki config Deck validation failed: %w", err))
	}
	err = validateNoteType(conf.NoteType)
	if err != nil {
		errs = append(errs, fmt.Errorf("anki config NoteType validation failed: %w", err))
	}
	mappingErrs := validateMappingKeys(conf.FieldMapping)
	for _, mappingErr := range mappingErrs {
		errs = append(errs, fmt.Errorf("anki config Mapping validation failed: %w", mappingErr))
	}
	mapping, mappingErrs := convertMapping(uc.Anki.FieldMapping)
	for _, mappingErr := range mappingErrs {
		errs = append(errs, mappingErr)
	}
	// if empty, it is disabled, so we don't need to check
	if conf.Audio.Field != "" {
		err = validateFieldName(conf.Audio.Field)
		if err != nil {
			errs = append(errs, fmt.Errorf("anki config Audio.Field validation failed: %w", err))
		}
	}

	return &Config{
		Addr:               uc.Anki.Addr,
		APIKey:             uc.Anki.APIKey,
		Deck:               uc.Anki.Deck,
		NoteType:           uc.Anki.NoteType,
		AudioField:         uc.Anki.Audio.Field,
		AudioPreferredType: uc.Anki.Audio.PreferredType,
		Mapping:            mapping,
	}, errors.Join(errs...)
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
	if err := validateAddr(addr); err != nil {
		return &ValidationError{Msg: err.Error()}
	}
	return cr.updateConfigFn(func(uc *config.UserConfig) error {
		uc.Anki.Addr = addr
		uc.Anki.APIKey = apikey
		return nil
	})
}

func (cr *ConfigReloader) UpdateDeck(name string) error {
	if err := validateDeckName(name); err != nil {
		return &ValidationError{Msg: err.Error()}
	}
	return cr.updateConfigFn(func(uc *config.UserConfig) error {
		uc.Anki.Deck = name
		return nil
	})
}

func (cr *ConfigReloader) UpdateNoteType(name string) error {
	if err := validateNoteType(name); err != nil {
		return &ValidationError{Msg: err.Error()}
	}
	return cr.updateConfigFn(func(uc *config.UserConfig) error {
		uc.Anki.NoteType = name
		return nil
	})
}

func (cr *ConfigReloader) UpdateMapping(mapping map[string]string) error {
	keyErrs := validateMappingKeys(mapping)
	_, valueErrs := convertMapping(mapping)
	if len(keyErrs) != 0 || len(valueErrs) != 0 {
		return &MappingValidationErrors{
			KeyErrors:   keyErrs,
			ValueErrors: valueErrs,
		}
	}
	return cr.updateConfigFn(func(uc *config.UserConfig) error {
		uc.Anki.FieldMapping = mapping
		return nil
	})
}

func (cr *ConfigReloader) UpdateAudio(field, preferredType string) error {
	if field != "" {
		if err := validateFieldName(field); err != nil {
			return &ValidationError{Msg: err.Error()}
		}
	}
	return cr.updateConfigFn(func(uc *config.UserConfig) error {
		uc.Anki.Audio.Field = field
		uc.Anki.Audio.PreferredType = preferredType
		return nil
	})
}
