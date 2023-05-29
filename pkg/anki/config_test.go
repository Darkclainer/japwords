package anki

import (
	"errors"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/config/configtest"
)

func Test_Config_Equal(t *testing.T) {
	testCases := []struct {
		Name     string
		First    *Config
		Second   *Config
		Expected bool
	}{
		{
			Name:     "empty",
			Expected: true,
		},
		{
			Name:     "first nil",
			Second:   &Config{},
			Expected: false,
		},
		{
			Name:     "second nil",
			First:    &Config{},
			Expected: false,
		},
		{
			Name:     "default values",
			First:    &Config{},
			Second:   &Config{},
			Expected: true,
		},
		{
			Name: "neq Addr",
			First: &Config{
				Addr: "a",
			},
			Second: &Config{
				Addr: "b",
			},
			Expected: false,
		},
		{
			Name: "neq APIKey",
			First: &Config{
				APIKey: "a",
			},
			Second: &Config{
				APIKey: "b",
			},
			Expected: false,
		},
		{
			Name: "neq Deck",
			First: &Config{
				Deck: "a",
			},
			Second: &Config{
				Deck: "b",
			},
			Expected: false,
		},
		{
			Name: "neq NoteType",
			First: &Config{
				NoteType: "a",
			},
			Second: &Config{
				NoteType: "b",
			},
			Expected: false,
		},
		{
			Name: "mapping eq nonempty",
			First: &Config{
				Mapping: map[string]*Template{
					"hello": {
						Src: "a",
					},
				},
			},
			Second: &Config{
				Mapping: map[string]*Template{
					"hello": {
						Src: "a",
					},
				},
			},
			Expected: true,
		},
		{
			Name: "mapping different len",
			First: &Config{
				Mapping: map[string]*Template{
					"hello": {
						Src: "a",
					},
				},
			},
			Second:   &Config{},
			Expected: false,
		},
		{
			Name: "mapping different key",
			First: &Config{
				Mapping: map[string]*Template{
					"hello": {
						Src: "a",
					},
				},
			},
			Second: &Config{
				Mapping: map[string]*Template{
					"world": {
						Src: "a",
					},
				},
			},
			Expected: false,
		},
		{
			Name: "mapping different value",
			First: &Config{
				Mapping: map[string]*Template{
					"hello": {
						Src: "a",
					},
				},
			},
			Second: &Config{
				Mapping: map[string]*Template{
					"hello": {
						Src: "b",
					},
				},
			},
			Expected: false,
		},
		{
			Name: "mapping template itself doesnt matter",
			First: &Config{
				Mapping: map[string]*Template{
					"hello": {
						Src:  "a",
						Tmpl: template.New("hello"),
					},
				},
			},
			Second: &Config{
				Mapping: map[string]*Template{
					"hello": {
						Src:  "a",
						Tmpl: template.New("world"),
					},
				},
			},
			Expected: true,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, tc.First.Equal(tc.Second))
		})
	}
}

func Test_ConfigReloader_New(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		configManager := configtest.New(t, &config.UserConfig{
			Anki: config.Anki{
				Addr:   "testaddr",
				APIKey: "testapikey",
			},
		})
		var factoryCalled bool
		anki := NewAnki(func(o *ankiconnect.Options) (AnkiClient, error) {
			assert.Equal(t, "testaddr", o.URL)
			assert.Equal(t, "testapikey", o.APIKey)
			factoryCalled = true
			return nil, nil
		})
		configReloader, err := NewConfigReloader(anki, configManager)
		assert.NotNil(t, configReloader)
		// we can check that way, that ConfigReloader implement Reloader interface
		assert.NotNil(t, configReloader.updateConfigFn)
		assert.NoError(t, err)
		assert.True(t, factoryCalled)
	})
	t.Run("Error", func(t *testing.T) {
		configManager := configtest.New(t, &config.UserConfig{
			Anki: config.Anki{
				Addr:   "testaddr",
				APIKey: "testapikey",
			},
		})
		anki := NewAnki(func(o *ankiconnect.Options) (AnkiClient, error) {
			assert.Equal(t, "testaddr", o.URL)
			assert.Equal(t, "testapikey", o.APIKey)
			return nil, errors.New("myerror")
		})
		configReloader, err := NewConfigReloader(anki, configManager)
		assert.Nil(t, configReloader)
		// we can check that way, that ConfigReloader implement Reloader interface
		assert.ErrorContains(t, err, "myerror")
	})
}
