package anki

import (
	"context"
	"errors"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
				Mapping: map[string]*Template{
					"a": {
						Tmpl: template.Must(template.New("").Parse("{{.Slug.Word}}")),
					},
				},
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

func Test_Anki_AddNote(t *testing.T) {
	request := &AddNoteRequest{
		Fields: []AddNoteField{
			{
				Name:  "hello",
				Value: "world",
			},
		},
		Tags:     []string{"a", "b"},
		AudioURL: "",
	}
	anki := NewAnki(func(_ *Config) (StatefullClient, error) {
		client := NewMockStatefullClient(t)
		client.On("AddNote", mock.Anything, request).Return(errors.New("myerror")).Once()
		return client, nil
	})
	err := anki.ReloadConfig(&Config{})
	require.NoError(t, err)
	err = anki.AddNote(context.Background(), request)
	assert.ErrorContains(t, err, "myerror")
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
