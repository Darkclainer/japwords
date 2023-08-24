package anki

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
)

func Test_Anki_New(t *testing.T) {
	t.Run("DefaultOK", func(t *testing.T) {
		anki := NewAnki(DefaultClientConstructor)
		assert.Nil(t, anki.client)
	})
}

func TestAnki_ReloadConfig(t *testing.T) {
	t.Run("Options", func(t *testing.T) {
		config := &Config{
			Addr:   "myaddr",
			APIKey: "mykey",
		}
		constructor := func(opts *ankiconnect.Options) (AnkiClient, error) {
			assert.Equal(t, "http://myaddr", opts.URL)
			assert.Equal(t, "mykey", opts.APIKey)
			return nil, nil
		}
		anki := NewAnki(constructor)
		err := anki.ReloadConfig(config)
		require.NoError(t, err)
		assert.Equal(t, config, anki.config)
	})
	t.Run("OK", func(t *testing.T) {
		config := Config{
			Addr:   "myaddr",
			APIKey: "myfirst",
		}
		counter := 0
		constructor := func(opts *ankiconnect.Options) (AnkiClient, error) {
			counter++
			assert.Equal(t, "http://myaddr", opts.URL)
			if counter == 1 {
				assert.Equal(t, "myfirst", opts.APIKey)
			} else {
				assert.Equal(t, "mysecond", opts.APIKey)
			}
			return nil, nil
		}
		anki := NewAnki(constructor)
		configCopy := config
		err := anki.ReloadConfig(&configCopy)
		require.NoError(t, err)
		config.APIKey = "mysecond"
		err = anki.ReloadConfig(&config)
		assert.NoError(t, err)
		assert.Equal(t, 2, counter)
	})
	t.Run("Error", func(t *testing.T) {
		config := Config{
			Addr:   "myaddr",
			APIKey: "myfirst",
		}
		counter := 0
		constructor := func(opts *ankiconnect.Options) (AnkiClient, error) {
			counter++
			if counter == 1 {
				return DefaultClientConstructor(opts)
			} else {
				return nil, errors.New("testerr")
			}
		}
		anki := NewAnki(constructor)
		err := anki.ReloadConfig(&config)
		require.NoError(t, err)
		client := anki.client
		err = anki.ReloadConfig(&config)
		assert.Error(t, err)
		assert.Same(t, client, anki.client)
	})
}

func Test_Anki_FullStateCheck_OK(t *testing.T) {
	testCases := []struct {
		Name            string
		Config          *Config
		Permissions     *ankiconnect.RequestPermissionResponse
		DeckNames       []string
		ModelNames      []string
		ModelFieldNames []string
		Expected        *StateResult
	}{
		{
			Name:   "permission denied",
			Config: &Config{},
			Permissions: &ankiconnect.RequestPermissionResponse{
				Permission:    ankiconnect.PermissionDenied,
				RequireAPIKey: true,
				Version:       5,
			},
			Expected: &StateResult{
				Connected:         true,
				Version:           5,
				PermissionGranted: false,
				APIKeyRequired:    true,
			},
		},
		{
			Name: "deck not exists",
			Config: &Config{
				Deck: "testdeck",
			},
			Permissions: &ankiconnect.RequestPermissionResponse{
				Permission: ankiconnect.PermissionGranted,
			},
			Expected: &StateResult{
				Connected:         true,
				PermissionGranted: true,
			},
		},
		{
			Name: "deck exists",
			Config: &Config{
				Deck: "testdeck",
			},
			Permissions: &ankiconnect.RequestPermissionResponse{
				Permission: ankiconnect.PermissionGranted,
			},
			DeckNames: []string{"mydeck", "testdeck"},
			Expected: &StateResult{
				Connected:         true,
				PermissionGranted: true,
				DeckExists:        true,
			},
		},
		{
			Name: "note type exists",
			Config: &Config{
				NoteType: "testnote",
			},
			Permissions: &ankiconnect.RequestPermissionResponse{
				Permission: ankiconnect.PermissionGranted,
			},
			ModelNames: []string{"mynote", "testnote"},
			Expected: &StateResult{
				Connected:         true,
				PermissionGranted: true,
				NoteTypeExists:    true,
			},
		},
		{
			Name: "missing fields",
			Config: &Config{
				NoteType: "testnote",
				Mapping: TemplateMapping{
					"key1": nil,
					"key2": nil,
					"key3": nil,
				},
			},
			Permissions: &ankiconnect.RequestPermissionResponse{
				Permission: ankiconnect.PermissionGranted,
			},
			ModelNames:      []string{"testnote"},
			ModelFieldNames: []string{"key2"},
			Expected: &StateResult{
				Connected:         true,
				PermissionGranted: true,
				NoteTypeExists:    true,
				NoteMissingFields: []string{"key1", "key3"},
			},
		},
		{
			Name: "no missing fields",
			Config: &Config{
				NoteType: "testnote",
				Mapping: TemplateMapping{
					"key1": nil,
					"key2": nil,
					"key3": nil,
				},
			},
			Permissions: &ankiconnect.RequestPermissionResponse{
				Permission: ankiconnect.PermissionGranted,
			},
			ModelNames:      []string{"testnote"},
			ModelFieldNames: []string{"key1", "key2", "key3", "key4"},
			Expected: &StateResult{
				Connected:         true,
				PermissionGranted: true,
				NoteTypeExists:    true,
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			var client *MockAnkiClient
			anki := NewAnki(func(_ *ankiconnect.Options) (AnkiClient, error) {
				client = NewMockAnkiClient(t)
				client.On("RequestPermission", mock.Anything).
					Return(tc.Permissions, nil).
					Once()
				client.On("DeckNames", mock.Anything).
					Return(tc.DeckNames, nil).
					Maybe()
				client.On("ModelNames", mock.Anything).
					Return(tc.ModelNames, nil).
					Maybe()
				client.On("ModelFieldNames", mock.Anything, tc.Config.NoteType).
					Return(tc.ModelFieldNames, nil).
					Maybe()
				return client, nil
			})
			err := anki.ReloadConfig(tc.Config)
			require.NoError(t, err)
			actual, err := anki.FullStateCheck(context.Background())
			assert.NoError(t, err)
			slices.Sort(actual.NoteMissingFields)
			assert.Equal(t, tc.Expected, actual)
			client.AssertExpectations(t)
		})
	}
}
