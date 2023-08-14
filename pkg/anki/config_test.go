package anki

import (
	"errors"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	ankiConfig := config.Anki{
		Addr:     "testaddr:3030",
		APIKey:   "testapikey",
		Deck:     "testdeck",
		NoteType: "testnote",
	}
	t.Run("OK", func(t *testing.T) {
		configManager := configtest.New(t, &config.UserConfig{
			Anki: ankiConfig,
		})
		var factoryCalled bool
		anki := NewAnki(func(o *ankiconnect.Options) (AnkiClient, error) {
			assert.Equal(t, ankiConfig.Addr, o.URL)
			assert.Equal(t, ankiConfig.APIKey, o.APIKey)
			factoryCalled = true
			return nil, nil
		})
		configReloader, err := NewConfigReloader(anki, configManager)
		require.NoError(t, err)
		assert.NotNil(t, configReloader)
		// we can check that way, that ConfigReloader implement Reloader interface
		assert.NotNil(t, configReloader.updateConfigFn)
		assert.True(t, factoryCalled)
	})
	t.Run("Error", func(t *testing.T) {
		configManager := configtest.New(t, &config.UserConfig{
			Anki: ankiConfig,
		})
		anki := NewAnki(func(o *ankiconnect.Options) (AnkiClient, error) {
			assert.Equal(t, ankiConfig.Addr, o.URL)
			assert.Equal(t, ankiConfig.APIKey, o.APIKey)
			return nil, errors.New("myerror")
		})
		configReloader, err := NewConfigReloader(anki, configManager)
		assert.Nil(t, configReloader)
		// we can check that way, that ConfigReloader implement Reloader interface
		assert.ErrorContains(t, err, "myerror")
	})
}

func Test_ConfigReloader_Config(t *testing.T) {
	testCases := []struct {
		Name        string
		UserConfig  *config.UserConfig
		Expected    *Config
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "Ok",
			UserConfig: &config.UserConfig{
				Anki: config.Anki{
					Addr:     "testaddr:3030",
					APIKey:   "testapikey",
					Deck:     "testdeck",
					NoteType: "testnote",
					FieldMapping: map[string]string{
						"mykey": "mymapping",
					},
				},
			},
			Expected: &Config{
				Addr:     "testaddr:3030",
				APIKey:   "testapikey",
				Deck:     "testdeck",
				NoteType: "testnote",
				Mapping: TemplateMapping{
					"mykey": &Template{
						Src: "mymapping",
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "invalid addr",
			UserConfig: &config.UserConfig{
				Anki: config.Anki{
					Addr:     "testaddr3030",
					Deck:     "testdeck",
					NoteType: "testnote",
				},
			},
			ErrorAssert: assert.Error,
		},
		{
			Name: "invalid deck",
			UserConfig: &config.UserConfig{
				Anki: config.Anki{
					Addr:     "testaddr:3030",
					Deck:     "test\"deck",
					NoteType: "testnote",
				},
			},
			ErrorAssert: assert.Error,
		},
		{
			Name: "invalid note type",
			UserConfig: &config.UserConfig{
				Anki: config.Anki{
					Addr:     "testaddr:3030",
					Deck:     "testdeck",
					NoteType: "test\"note",
				},
			},
			ErrorAssert: assert.Error,
		},
		{
			Name: "invalid mapping field key",
			UserConfig: &config.UserConfig{
				Anki: config.Anki{
					Addr:     "testaddr:3030",
					Deck:     "testdeck",
					NoteType: "testnote",
					FieldMapping: map[string]string{
						"hello:there": "hi",
					},
				},
			},
			ErrorAssert: assert.Error,
		},
		{
			Name: "invalid mapping field value",
			UserConfig: &config.UserConfig{
				Anki: config.Anki{
					Addr:     "testaddr:3030",
					Deck:     "testdeck",
					NoteType: "testnote",
					FieldMapping: map[string]string{
						"hellothere": "{{hi",
					},
				},
			},
			ErrorAssert: assert.Error,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			// we don't need anything for test, this method could be function
			reloader := &ConfigReloader{}
			actual, err := reloader.Config(tc.UserConfig)
			tc.ErrorAssert(t, err)
			if err != nil {
				return
			}
			if !tc.Expected.Equal(actual) {
				// for nice message we need to nil all templates, because I don't want to compare them with assert
				actual := actual.(*Config)
				if actual.Mapping != nil {
					for _, template := range actual.Mapping {
						template.Tmpl = nil
					}
				}
				assert.Equal(t, tc.Expected, actual)
			}
		})
	}
}

func Test_ConfigReloader_Reload(t *testing.T) {
	ankiConfig := config.Anki{
		Addr:     "testaddr:3030",
		APIKey:   "testapikey",
		Deck:     "testdeck",
		NoteType: "testnote",
	}
	configManager := configtest.New(t, &config.UserConfig{
		Anki: ankiConfig,
	})
	var factoryCalled int
	anki := NewAnki(func(_ *ankiconnect.Options) (AnkiClient, error) {
		factoryCalled++
		return nil, nil
	})
	configReloader, err := NewConfigReloader(anki, configManager)
	require.NoError(t, err)
	// for Reload validation is not called
	conf := &Config{
		Addr:     "myaddr",
		APIKey:   "myapikey",
		Deck:     "mydeck",
		NoteType: "mynotetype",
	}
	err = configReloader.Reload(conf)
	require.NoError(t, err)
	assert.Equal(t, 2, factoryCalled)
	assert.Equal(t, conf, anki.wrapper.config)
}

func Test_ConfigReloader_UpdateConnection(t *testing.T) {
}

func Test_ConfigReloader_UpdateDeck(t *testing.T) {
}

func Test_ConfigReloader_UpdateNoteType(t *testing.T) {
}

func Test_ConfigReloader_UpdateMapping(t *testing.T) {
}
