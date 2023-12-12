package anki

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
	"github.com/Darkclainer/japwords/pkg/lemma"
)

func Test_Anki_New(t *testing.T) {
	t.Run("DefaultOK", func(t *testing.T) {
		anki := NewAnki(DefaultStatefullClientConstructor)
		assert.Nil(t, anki.client)
	})
}

func Test_Anki_ReloadConfig(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		config := Config{
			Addr:   "myaddr",
			APIKey: "myfirst",
		}
		counter := 0
		var clients []*MockStatefullClient
		constructor := func(config *Config) (StatefullClient, error) {
			counter++
			assert.Equal(t, "myaddr", config.Addr)
			if counter == 1 {
				assert.Equal(t, "myfirst", config.APIKey)
			} else {
				assert.Equal(t, "mysecond", config.APIKey)
			}
			client := NewMockStatefullClient(t)
			client.On("Stop").Return().Maybe()
			clients = append(clients, client)
			return client, nil
		}
		anki := NewAnki(constructor)
		configCopy := config
		err := anki.ReloadConfig(&configCopy)
		require.NoError(t, err)
		config.APIKey = "mysecond"
		err = anki.ReloadConfig(&config)
		assert.NoError(t, err)
		assert.Equal(t, 2, counter)
		clients[0].AssertCalled(t, "Stop")
		clients[1].AssertNotCalled(t, "Stop")
	})
	t.Run("Error", func(t *testing.T) {
		config := Config{
			Addr:   "myaddr",
			APIKey: "myfirst",
		}
		counter := 0
		constructor := func(_ *Config) (StatefullClient, error) {
			counter++
			if counter == 1 {
				// do not define Stop, so we now if it's called for some reason
				client := NewMockStatefullClient(t)
				return client, nil
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

func Test_Anki_FullStateCheck(t *testing.T) {
	testCases := []struct {
		Name        string
		State       *State
		Error       error
		Expected    *StateResult
		AssertError assert.ErrorAssertionFunc
	}{
		{
			Name:  "Error",
			State: nil,
			Error: errors.New("myerror"),
			AssertError: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "myerror")
			},
		},
		{
			Name: "Version",
			State: &State{
				AnkiState: AnkiState{
					Version: 99,
				},
			},
			Expected: &StateResult{
				Version: 99,
			},
			AssertError: assert.NoError,
		},
		{
			Name: "DeckExists",
			State: &State{
				DeckExists: true,
			},
			Expected: &StateResult{
				DeckExists: true,
			},
			AssertError: assert.NoError,
		},
		{
			Name: "NoteTypeExists",
			State: &State{
				NoteTypeExists: true,
			},
			Expected: &StateResult{
				NoteTypeExists: true,
			},
			AssertError: assert.NoError,
		},
		{
			Name: "NoteHasAllFields",
			State: &State{
				NoteHasAllFields: true,
			},
			Expected: &StateResult{
				NoteHasAllFields: true,
			},
			AssertError: assert.NoError,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			var client *MockStatefullClient
			anki := NewAnki(func(_ *Config) (StatefullClient, error) {
				client = NewMockStatefullClient(t)
				client.On("GetState", mock.Anything).
					Return(tc.State, tc.Error).
					Once()
				return client, nil
			})
			err := anki.ReloadConfig(&Config{})
			require.NoError(t, err)
			actual, err := anki.FullStateCheck(context.Background())
			tc.AssertError(t, err)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_Anki_Decks(t *testing.T) {
	testCases := []struct {
		Name        string
		State       *State
		Error       error
		Expected    []string
		AssertError assert.ErrorAssertionFunc
	}{
		{
			Name:  "Error",
			State: nil,
			Error: errors.New("myerror"),
			AssertError: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "myerror")
			},
		},
		{
			Name: "OK",
			State: &State{
				AnkiState: AnkiState{
					Decks: []string{"a", "b"},
				},
			},
			Expected:    []string{"a", "b"},
			AssertError: assert.NoError,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			var client *MockStatefullClient
			anki := NewAnki(func(_ *Config) (StatefullClient, error) {
				client = NewMockStatefullClient(t)
				client.On("GetState", mock.Anything).
					Return(tc.State, tc.Error).
					Once()
				return client, nil
			})
			err := anki.ReloadConfig(&Config{})
			require.NoError(t, err)
			actual, err := anki.Decks(context.Background())
			tc.AssertError(t, err)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_Anki_NoteTypes(t *testing.T) {
	testCases := []struct {
		Name        string
		State       *State
		Error       error
		Expected    []string
		AssertError assert.ErrorAssertionFunc
	}{
		{
			Name:  "Error",
			State: nil,
			Error: errors.New("myerror"),
			AssertError: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "myerror")
			},
		},
		{
			Name: "OK",
			State: &State{
				AnkiState: AnkiState{
					NoteTypes: []string{"a", "b"},
				},
			},
			Expected:    []string{"a", "b"},
			AssertError: assert.NoError,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			var client *MockStatefullClient
			anki := NewAnki(func(_ *Config) (StatefullClient, error) {
				client = NewMockStatefullClient(t)
				client.On("GetState", mock.Anything).
					Return(tc.State, tc.Error).
					Once()
				return client, nil
			})
			err := anki.ReloadConfig(&Config{})
			require.NoError(t, err)
			actual, err := anki.NoteTypes(context.Background())
			tc.AssertError(t, err)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_Anki_NoteTypeFields(t *testing.T) {
	testCases := []struct {
		Name        string
		State       *State
		Error       error
		Expected    []string
		AssertError assert.ErrorAssertionFunc
	}{
		{
			Name:  "Error",
			State: nil,
			Error: errors.New("myerror"),
			AssertError: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "myerror")
			},
		},
		{
			Name: "OK",
			State: &State{
				NoteTypeExists: true,
				CurrentFields:  []string{"a", "b"},
			},
			Expected:    []string{"a", "b"},
			AssertError: assert.NoError,
		},
		{
			Name: "NoteType not exists",
			State: &State{
				CurrentFields: []string{"a", "b"},
			},
			Expected: nil,
			AssertError: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrNoteTypeNotExists)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			var client *MockStatefullClient
			anki := NewAnki(func(_ *Config) (StatefullClient, error) {
				client = NewMockStatefullClient(t)
				client.On("GetState", mock.Anything).
					Return(tc.State, tc.Error).
					Once()
				return client, nil
			})
			err := anki.ReloadConfig(&Config{})
			require.NoError(t, err)
			actual, err := anki.NoteTypeFields(context.Background())
			tc.AssertError(t, err)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_Anki_CreateDeck(t *testing.T) {
	var client *MockStatefullClient
	expectedErr := errors.New("myerror")
	anki := NewAnki(func(_ *Config) (StatefullClient, error) {
		client = NewMockStatefullClient(t)
		client.On("CreateDeck", mock.Anything, "newdeck").Return(expectedErr).Once()
		return client, nil
	})
	err := anki.ReloadConfig(&Config{})
	require.NoError(t, err)
	actualErr := anki.CreateDeck(context.Background(), "newdeck")
	assert.ErrorIs(t, actualErr, expectedErr)
}

func Test_Anki_CreateDefaultNote(t *testing.T) {
	var client *MockStatefullClient
	expectedErr := errors.New("myerror")
	anki := NewAnki(func(_ *Config) (StatefullClient, error) {
		client = NewMockStatefullClient(t)
		client.On("CreateDefaultNoteType", mock.Anything, "newnotetype").Return(expectedErr).Once()
		return client, nil
	})
	err := anki.ReloadConfig(&Config{})
	require.NoError(t, err)
	actualErr := anki.CreateDefaultNote(context.Background(), "newnotetype")
	assert.ErrorIs(t, actualErr, expectedErr)
}

func Test_Anki_PrepareProjectedLemma(t *testing.T) {
	readyState := &State{
		DeckExists:       true,
		NoteTypeExists:   true,
		NoteHasAllFields: true,
		OrderDefined:     true,
	}
	testCases := []struct {
		Name        string
		InitClient  func(config *Config, client *MockStatefullClient)
		Config      *Config
		AssertError assert.ErrorAssertionFunc
		Expected    *AddNoteRequest
	}{
		{
			Name: "GetState Error",
			InitClient: func(_ *Config, client *MockStatefullClient) {
				client.On("GetState", mock.Anything).
					Return(nil, errors.New("GetState error"))
			},
			AssertError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "GetState error")
			},
		},
		{
			Name: "state not ready",
			InitClient: func(_ *Config, client *MockStatefullClient) {
				client.On("GetState", mock.Anything).
					Return(&State{}, nil)
			},
			AssertError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, ErrIncompleteConfiguration)
			},
		},
		{
			Name: "mapping error",
			InitClient: func(conf *Config, client *MockStatefullClient) {
				newState := *readyState
				newState.CurrentFields = []string{"a"}
				client.On("GetState", mock.Anything).
					Return(&newState, nil)
				client.On("Config").Return(conf)
			},
			Config: &Config{
				Mapping: map[string]*Template{
					"a": {
						Tmpl: template.Must(template.New("").Parse("{{.NotExists}}")),
					},
				},
			},
			AssertError: assert.Error,
		},
		{
			Name: "default fields",
			InitClient: func(conf *Config, client *MockStatefullClient) {
				newState := *readyState
				newState.CurrentFields = []string{"a", "b"}
				client.On("GetState", mock.Anything).
					Return(&newState, nil)
				client.On("Config").Return(conf)
			},
			Config: &Config{
				Mapping: map[string]*Template{},
			},
			AssertError: assert.NoError,
			Expected: &AddNoteRequest{
				Fields: []AddNoteField{
					{
						Name:  "a",
						Value: "",
					},
					{
						Name:  "b",
						Value: "",
					},
				},
			},
		},
		{
			Name: "mapping",
			InitClient: func(conf *Config, client *MockStatefullClient) {
				newState := *readyState
				newState.CurrentFields = []string{"a", "b"}
				client.On("GetState", mock.Anything).
					Return(&newState, nil)
				client.On("Config").Return(conf)
			},
			Config: &Config{
				Mapping: mustConvertMapping(t, map[string]string{
					"a": "{{.Slug.Word}}",
				}),
			},
			AssertError: assert.NoError,
			Expected: &AddNoteRequest{
				Fields: []AddNoteField{
					{
						Name:  "a",
						Value: "一二わ三はい",
					},
					{
						Name:  "b",
						Value: "",
					},
				},
			},
		},
		{
			Name: "audio assets",
			InitClient: func(conf *Config, client *MockStatefullClient) {
				newState := *readyState
				newState.CurrentFields = []string{"a", "b"}
				client.On("GetState", mock.Anything).
					Return(&newState, nil)
				client.On("Config").Return(conf)
			},
			Config: &Config{
				AudioField:         "Audio",
				AudioPreferredType: "mpeg",
			},
			AssertError: assert.NoError,
			Expected: &AddNoteRequest{
				Fields: []AddNoteField{
					{
						Name: "a",
					},
					{
						Name: "b",
					},
				},
				AudioAssets: []AddNoteAudioAsset{
					{
						Field:    "Audio",
						Filename: "一二わ三はい-いちにわさんはい.mp3",
						URL:      "https://example.com/somelink/mp3",
					},
					{
						Field:    "Audio",
						Filename: "一二わ三はい-いちにわさんはい.oga",
						URL:      "https://example.com/somelink/ogg",
					},
				},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			anki := NewAnki(func(conf *Config) (StatefullClient, error) {
				client := NewMockStatefullClient(t)
				tc.InitClient(conf, client)
				return client, nil
			})
			err := anki.ReloadConfig(tc.Config)
			require.NoError(t, err)
			actual, err := anki.PrepareProjectedLemma(context.Background(), &DefaultExampleLemma)
			tc.AssertError(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_prepareFieldsForNoteRequest(t *testing.T) {
	testCases := []struct {
		Name          string
		CurrentFields []string
		Mapping       TemplateMapping
		Expected      []AddNoteField
		AssertError   assert.ErrorAssertionFunc
	}{
		{
			Name:        "empty",
			Expected:    []AddNoteField{},
			AssertError: assert.NoError,
		},
		{
			Name:          "fields without mapping",
			CurrentFields: []string{"foo", "bar"},
			Expected: []AddNoteField{
				{Name: "foo"},
				{Name: "bar"},
			},
			AssertError: assert.NoError,
		},
		{
			Name:          "field with mapping",
			CurrentFields: []string{"foo"},
			Mapping: mustConvertMapping(t, map[string]string{
				"foo": "{{.Slug.Word}}",
			}),
			Expected: []AddNoteField{
				{Name: "foo", Value: DefaultExampleLemma.Slug.Word},
			},
			AssertError: assert.NoError,
		},
		{
			Name:          "error in template",
			CurrentFields: []string{"foo"},
			Mapping: mustConvertMapping(t, map[string]string{
				"foo": `{{ if gt (len .Slug.Word) 0 }}{{ fail "intentional fail" }}{{ end }}`,
			}),
			AssertError: assert.Error,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			actual, err := prepareFieldsForNoteRequest(&DefaultExampleLemma, tc.CurrentFields, tc.Mapping)
			tc.AssertError(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_prepareAudiosForNoteRequest(t *testing.T) {
	assets := prepareAudiosForNoteRequest(
		&lemma.ProjectedLemma{
			Slug: lemma.Word{
				Word:     "foo",
				Hiragana: "bar",
			},
			Audio: []lemma.Audio{
				{
					MediaType: "audio/ogg",
					Source:    "ogglink",
				},
				{
					MediaType: "audio/mpeg",
					Source:    "mpeglink",
				},
				{
					MediaType: "",
					Source:    "uknownlink",
				},
			},
		},
		&Config{
			AudioField:         "audiofield",
			AudioPreferredType: "mpeg",
		},
	)
	assert.Equal(
		t,
		[]AddNoteAudioAsset{
			{
				Field:    "audiofield",
				Filename: "foo-bar.mp3",
				URL:      "mpeglink",
			},
			{
				Field:    "audiofield",
				Filename: "foo-bar.oga",
				URL:      "ogglink",
			},
			{
				Field:    "audiofield",
				Filename: "foo-bar",
				URL:      "uknownlink",
			},
		},
		assets,
	)
}

func Test_generateLemmaAudioBasename(t *testing.T) {
	testCases := []struct {
		Name     string
		Word     *lemma.Word
		Expected string
	}{
		{
			Name: "empty",
			Word: &lemma.Word{},
		},
		{
			Name: "word only",
			Word: &lemma.Word{
				Word: "foo",
			},
			Expected: "foo",
		},
		{
			Name: "word and hiragana",
			Word: &lemma.Word{
				Word:     "foo",
				Hiragana: "bar",
			},
			Expected: "foo-bar",
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			actual := generateLemmaAudioBasename(&lemma.ProjectedLemma{
				Slug: *tc.Word,
			})
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_sortLemmaAudios(t *testing.T) {
	testCases := []struct {
		Name          string
		PreferredType string
		Audios        []lemma.Audio
		Expected      []lemma.Audio
	}{
		{
			Name:     "empty",
			Expected: []lemma.Audio{},
		},
		{
			Name:          "no preferred",
			PreferredType: "",
			Audios: []lemma.Audio{
				{
					MediaType: "foo",
					Source:    "a",
				},
				{
					MediaType: "foo",
					Source:    "b",
				},
				{
					MediaType: "",
					Source:    "c",
				},
				{
					MediaType: "bar",
					Source:    "d",
				},
			},
			Expected: []lemma.Audio{
				{
					MediaType: "foo",
					Source:    "a",
				},
				{
					MediaType: "foo",
					Source:    "b",
				},
				{
					MediaType: "",
					Source:    "c",
				},
				{
					MediaType: "bar",
					Source:    "d",
				},
			},
		},
		{
			Name:          "preferred stable",
			PreferredType: "bar",
			Audios: []lemma.Audio{
				{
					MediaType: "foo",
					Source:    "a",
				},
				{
					MediaType: "bar",
					Source:    "b",
				},
				{
					MediaType: "foo",
					Source:    "c",
				},
				{
					MediaType: "bar",
					Source:    "d",
				},
			},
			Expected: []lemma.Audio{
				{
					MediaType: "bar",
					Source:    "b",
				},
				{
					MediaType: "bar",
					Source:    "d",
				},
				{
					MediaType: "foo",
					Source:    "a",
				},
				{
					MediaType: "foo",
					Source:    "c",
				},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			actual := sortLemmaAudios(tc.Audios, tc.PreferredType)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_Anki_AddNote(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		request := &AddNoteRequest{
			Fields: []AddNoteField{
				{
					Name:  "hello",
					Value: "world",
				},
			},
			Tags: []string{"a", "b"},
		}
		anki := NewAnki(func(_ *Config) (StatefullClient, error) {
			client := NewMockStatefullClient(t)
			client.On("AddNote", mock.Anything, request).Return(int64(32), errors.New("myerror")).Once()
			return client, nil
		})
		err := anki.ReloadConfig(&Config{})
		require.NoError(t, err)
		noteID, err := anki.AddNote(context.Background(), request)
		assert.Equal(t, NoteID(32), noteID)
		assert.ErrorContains(t, err, "myerror")
	})
	t.Run("filter audio", func(t *testing.T) {
		request := &AddNoteRequest{
			AudioAssets: []AddNoteAudioAsset{
				{
					Field:    "foo",
					Filename: "a",
				},
				{
					Field:    "bar",
					Filename: "b",
				},
				{
					Field:    "foo",
					Filename: "c",
				},
				{
					Field:    "bar",
					Filename: "d",
				},
			},
		}
		expectedRequest := &AddNoteRequest{
			AudioAssets: []AddNoteAudioAsset{
				{
					Field:    "foo",
					Filename: "a",
				},
				{
					Field:    "bar",
					Filename: "b",
				},
			},
		}
		anki := NewAnki(func(_ *Config) (StatefullClient, error) {
			client := NewMockStatefullClient(t)
			client.On("AddNote", mock.Anything, expectedRequest).Return(int64(32), nil).Once()
			return client, nil
		})
		err := anki.ReloadConfig(&Config{})
		require.NoError(t, err)
		noteID, err := anki.AddNote(context.Background(), request)
		assert.Equal(t, NoteID(32), noteID)
		assert.NoError(t, err)
	})
}

func Test_Anki_SearchProjectedLemmas(t *testing.T) {
	t.Run("get state error", func(t *testing.T) {
		anki := NewAnki(func(_ *Config) (StatefullClient, error) {
			client := NewMockStatefullClient(t)
			client.On("GetState", mock.Anything).
				Return(
					nil,
					errors.New("myerror"),
				)

			return client, nil
		})
		err := anki.ReloadConfig(&Config{})
		require.NoError(t, err)
		_, err = anki.SearchProjectedLemmas(context.Background(), nil)
		assert.ErrorContains(t, err, "myerror")
	})
	t.Run("state not ready", func(t *testing.T) {
		anki := NewAnki(func(_ *Config) (StatefullClient, error) {
			client := NewMockStatefullClient(t)
			client.On("GetState", mock.Anything).
				Return(
					&State{},
					nil,
				)

			return client, nil
		})
		err := anki.ReloadConfig(&Config{})
		require.NoError(t, err)
		_, err = anki.SearchProjectedLemmas(context.Background(), nil)
		assert.ErrorIs(t, err, ErrIncompleteConfiguration)
	})
	t.Run("no current fields", func(t *testing.T) {
		anki := NewAnki(func(_ *Config) (StatefullClient, error) {
			client := NewMockStatefullClient(t)
			client.On("GetState", mock.Anything).
				Return(
					&State{
						DeckExists:       true,
						NoteTypeExists:   true,
						NoteHasAllFields: true,
						OrderDefined:     true,
					},
					nil,
				)

			return client, nil
		})
		err := anki.ReloadConfig(&Config{})
		require.NoError(t, err)
		_, err = anki.SearchProjectedLemmas(context.Background(), nil)
		assert.ErrorIs(t, err, ErrIncompleteConfiguration)
	})
	readyState := &State{
		DeckExists:       true,
		NoteTypeExists:   true,
		NoteHasAllFields: true,
		OrderDefined:     true,
		CurrentFields:    []string{"of"},
	}
	t.Run("generate query returns error", func(t *testing.T) {
		anki := NewAnki(func(_ *Config) (StatefullClient, error) {
			client := NewMockStatefullClient(t)
			client.On("GetState", mock.Anything).
				Return(readyState, nil)
			client.On("Config").
				Return(&Config{
					// it will return error, because order field doesn't have template
					Mapping: TemplateMapping{},
				})

			return client, nil
		})
		err := anki.ReloadConfig(&Config{})
		require.NoError(t, err)
		_, err = anki.SearchProjectedLemmas(context.Background(), nil)
		assert.ErrorIs(t, err, ErrIncompleteConfiguration)
	})
	config := &Config{
		Deck:     "mydeck",
		NoteType: "mynote",
		Mapping: mustConvertMapping(t, map[string]string{
			"of": `{{.Slug.Word}}`,
		}),
	}
	t.Run("generate query returns error", func(t *testing.T) {
		anki := NewAnki(func(_ *Config) (StatefullClient, error) {
			client := NewMockStatefullClient(t)
			client.On("GetState", mock.Anything).
				Return(readyState, nil)
			client.On("Config").
				Return(config)
			client.On("QueryNotes", mock.Anything, `("deck:mydeck" "note:mynote" )`).
				Return(nil, errors.New("myerror"))

			return client, nil
		})
		err := anki.ReloadConfig(&Config{})
		require.NoError(t, err)
		_, err = anki.SearchProjectedLemmas(context.Background(), nil)
		assert.ErrorContains(t, err, "myerror")
	})
	t.Run("OK", func(t *testing.T) {
		anki := NewAnki(func(_ *Config) (StatefullClient, error) {
			client := NewMockStatefullClient(t)
			client.On("GetState", mock.Anything).
				Return(readyState, nil)
			client.On("Config").
				Return(config)
			client.On("QueryNotes", mock.Anything, `("deck:mydeck" "note:mynote" ("of:hello" OR "of:world"))`).
				Return(
					[]*ankiconnect.NoteInfo{
						{
							NoteID: 2,
							Fields: map[string]*ankiconnect.NoteInfoField{
								"of": {
									Value: "world",
								},
							},
						},
						{
							NoteID: 1,
							Fields: map[string]*ankiconnect.NoteInfoField{
								"of": {
									Value: "hello",
								},
							},
						},
					},
					nil,
				)

			return client, nil
		})
		err := anki.ReloadConfig(&Config{})
		require.NoError(t, err)
		actual, err := anki.SearchProjectedLemmas(context.Background(), []*lemma.ProjectedLemma{
			{
				Slug: lemma.Word{
					Word: "hello",
				},
			},
			{
				Slug: lemma.Word{
					Word: "world",
				},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, []NoteID{1, 2}, actual)
	})
}

func Test_generateQueryForNotes_NoOrderTemplate(t *testing.T) {
	_, _, err := generateQueryForNotes(nil, "hello", &Config{})
	assert.ErrorIs(t, err, ErrIncompleteConfiguration)
}

func Test_generateQueryForNotes_DeckNote(t *testing.T) {
	query, _, err := generateQueryForNotes(nil, "hello", &Config{
		Mapping: mustConvertMapping(t, map[string]string{
			"hello": "",
		}),
		NoteType: "note1",
		Deck:     "deck1",
	})
	require.NoError(t, err)
	assert.Equal(t, `("deck:deck1" "note:note1" )`, query)
	query, _, err = generateQueryForNotes(nil, "hello", &Config{
		Mapping: mustConvertMapping(t, map[string]string{
			"hello": "",
		}),
		NoteType: "note2",
		Deck:     "deck2",
	})
	require.NoError(t, err)
	assert.Equal(t, `("deck:deck2" "note:note2" )`, query)
}

func Test_generateQueryForNotes(t *testing.T) {
	deck := "mydeck"
	noteType := "mynote"
	orderField := "of"
	ornamentQuery := func(s string) string {
		return fmt.Sprintf(`("deck:%s" "note:%s" (%s))`, deck, noteType, s)
	}
	testCases := []struct {
		Name           string
		Lemmas         []*lemma.ProjectedLemma
		OrderTemplate  string
		ExpectedQuery  string
		ExpectedValues []string
		ErrorAssert    assert.ErrorAssertionFunc
	}{
		{
			Name: "template execute failed",
			Lemmas: []*lemma.ProjectedLemma{
				// need to fail on this note specifically
				{
					Slug: lemma.Word{
						Word: "hello",
					},
				},
			},
			OrderTemplate: `{{ if  eq .Slug.Word "hello" }}{{ fail "intended to fail" }}{{ end }}`,
			ErrorAssert:   assert.Error,
		},
		{
			Name: "constant template",
			Lemmas: []*lemma.ProjectedLemma{
				{},
				{},
				{},
			},
			OrderTemplate:  `foo`,
			ExpectedQuery:  ornamentQuery(`"of:foo" OR "of:foo" OR "of:foo"`),
			ExpectedValues: []string{"foo", "foo", "foo"},
			ErrorAssert:    assert.NoError,
		},
		{
			Name: "slug word",
			Lemmas: []*lemma.ProjectedLemma{
				{
					Slug: lemma.Word{
						Word: "word1",
					},
				},
				{
					Slug: lemma.Word{
						Word: "word2",
					},
				},
			},
			OrderTemplate:  `{{.Slug.Word}}`,
			ExpectedQuery:  ornamentQuery(`"of:word1" OR "of:word2"`),
			ExpectedValues: []string{"word1", "word2"},
			ErrorAssert:    assert.NoError,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			config := &Config{
				Deck:     deck,
				NoteType: noteType,
				Mapping: mustConvertMapping(t, map[string]string{
					orderField: tc.OrderTemplate,
				}),
			}
			actualQuery, actualValues, err := generateQueryForNotes(tc.Lemmas, orderField, config)
			tc.ErrorAssert(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.ExpectedQuery, actualQuery)
			assert.Equal(t, tc.ExpectedValues, actualValues)
		})
	}
}

func Test_confirmFoundNotes(t *testing.T) {
	const orderField = "of"
	testCases := []struct {
		Name        string
		Notes       []*ankiconnect.NoteInfo
		OrderValues []string
		Expected    []NoteID
	}{
		{
			Name: "wrong order field",
			Notes: []*ankiconnect.NoteInfo{
				{},
			},
			Expected: []NoteID{},
		},
		{
			Name: "one for all",
			Notes: []*ankiconnect.NoteInfo{
				{
					NoteID: 1,
					Fields: map[string]*ankiconnect.NoteInfoField{
						"of": {
							Value: "foo",
						},
					},
				},
			},
			OrderValues: []string{"foo", "notfound", "foo"},
			Expected:    []NoteID{1, 0, 1},
		},
		{
			Name: "found",
			Notes: []*ankiconnect.NoteInfo{
				{
					NoteID: 3,
					Fields: map[string]*ankiconnect.NoteInfoField{
						"of": {
							Value: "foo",
						},
					},
				},
				{
					NoteID: 2,
					Fields: map[string]*ankiconnect.NoteInfoField{
						"of": {
							Value: "bar",
						},
					},
				},
				{
					NoteID: 1,
					Fields: map[string]*ankiconnect.NoteInfoField{
						"of": {
							Value: "foobar",
						},
					},
				},
			},
			OrderValues: []string{"foobar", "bar", "foo", "notfound"},
			Expected:    []NoteID{1, 2, 3, 0},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			actual := confirmFoundNotes(tc.Notes, orderField, tc.OrderValues)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_Anki_Stop(t *testing.T) {
	var client *MockStatefullClient
	anki := NewAnki(func(_ *Config) (StatefullClient, error) {
		client = NewMockStatefullClient(t)
		client.On("Stop").Return().Once()
		return client, nil
	})
	err := anki.ReloadConfig(&Config{})
	require.NoError(t, err)
	anki.Stop()
}
