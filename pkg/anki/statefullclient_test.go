package anki

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
)

func Test_statefullClient_getState(t *testing.T) {
	testCases := []struct {
		Name            string
		Config          *Config
		Permissions     *ankiconnect.RequestPermissionResponse
		DeckNames       []string
		ModelNames      []string
		ModelFieldNames []string
		Expected        *State
		AssertError     assert.ErrorAssertionFunc
	}{
		{
			Name:   "permission denied",
			Config: &Config{},
			Permissions: &ankiconnect.RequestPermissionResponse{
				Permission:    ankiconnect.PermissionDenied,
				RequireAPIKey: true,
				Version:       5,
			},
			Expected: &State{},
			AssertError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, ErrForbiddenOrigin)
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
			Expected: &State{
				DeckExists: false,
			},
			AssertError: assert.NoError,
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
			Expected: &State{
				AnkiState: AnkiState{
					Decks: []string{"mydeck", "testdeck"},
				},
				DeckExists: true,
			},
			AssertError: assert.NoError,
		},
		{
			Name: "note type exists",
			Config: &Config{
				NoteType: "testnote",
			},
			Permissions: &ankiconnect.RequestPermissionResponse{
				Permission: ankiconnect.PermissionGranted,
			},
			ModelNames:      []string{"mynote", "testnote"},
			ModelFieldNames: []string{"key1"},
			Expected: &State{
				AnkiState: AnkiState{
					NoteTypes: []string{"mynote", "testnote"},
					NoteFields: map[string][]string{
						"testnote": {"key1"},
					},
				},
				CurrentFields:    []string{"key1"},
				NoteTypeExists:   true,
				NoteHasAllFields: true,
			},
			AssertError: assert.NoError,
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
			Expected: &State{
				AnkiState: AnkiState{
					NoteTypes: []string{"testnote"},
					NoteFields: map[string][]string{
						"testnote": {"key2"},
					},
				},
				NoteTypeExists: true,
				CurrentFields:  []string{"key2"},
				OrderDefined:   true,
			},
			AssertError: assert.NoError,
		},
		{
			Name: "missing fields order undefined",
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
			ModelFieldNames: []string{"keyUknown", "key1"},
			Expected: &State{
				AnkiState: AnkiState{
					NoteTypes: []string{"testnote"},
					NoteFields: map[string][]string{
						"testnote": {"keyUknown", "key1"},
					},
				},
				NoteTypeExists: true,
				CurrentFields:  []string{"keyUknown", "key1"},
				OrderDefined:   false,
			},
			AssertError: assert.NoError,
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
			Expected: &State{
				AnkiState: AnkiState{
					NoteTypes: []string{"testnote"},
					NoteFields: map[string][]string{
						"testnote": {"key1", "key2", "key3", "key4"},
					},
				},
				NoteTypeExists:   true,
				CurrentFields:    []string{"key1", "key2", "key3", "key4"},
				NoteHasAllFields: true,
				OrderDefined:     true,
			},
			AssertError: assert.NoError,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ankiClient := NewMockAnkiClient(t)
			ankiClient.On("RequestPermission", mock.Anything).
				Return(tc.Permissions, nil).
				Once()
			ankiClient.On("DeckNames", mock.Anything).
				Return(tc.DeckNames, nil).
				Maybe()
			ankiClient.On("ModelNames", mock.Anything).
				Return(tc.ModelNames, nil).
				Maybe()
			if len(tc.ModelFieldNames) != 0 {
				ankiClient.On("ModelFieldNames", mock.Anything, tc.Config.NoteType).
					Return(tc.ModelFieldNames, nil).
					Maybe()
			} else {
				ankiClient.On("ModelFieldNames", mock.Anything, tc.Config.NoteType).
					Return(nil, &ankiconnect.ServerError{
						Message: "model was not found:",
					}).
					Maybe()
			}
			client := newStatefullClientImpl(ankiClient, tc.Config, &statefullClientOptions{})
			actual := client.getNewState(context.Background())
			tc.AssertError(t, actual.LastError)
			actual.LastError = nil
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

// Test_statefullClient_getState_errors checks that if ankiclient return error, it is converted and assigned in State
func Test_statefullClient_getState_errors(t *testing.T) {
	// methods that will be called in their order
	methods := []string{"RequestPermission", "DeckNames", "ModelNames", "ModelFieldNames"}
	methodsParams := map[string]struct {
		Params []any
		Return []any
	}{
		"RequestPermission": {
			Return: []any{
				&ankiconnect.RequestPermissionResponse{
					Permission: ankiconnect.PermissionGranted,
				},
				nil,
			},
		},
		"DeckNames": {
			Return: []any{
				[]string{"testdeck"},
				nil,
			},
		},
		"ModelNames": {
			Return: []any{
				[]string{"testnote"},
				nil,
			},
		},
		"ModelFieldNames": {
			Params: []any{
				"testnote",
			},
			Return: []any{
				[]string{"key1"},
				nil,
			},
		},
	}
	config := &Config{
		Deck:     "testdeck",
		NoteType: "testnote",
		Mapping: TemplateMapping{
			"key1": nil,
		},
	}

	testCases := []struct {
		Name        string
		Error       error
		AssertError assert.ErrorAssertionFunc
	}{
		// there can be same cases as for convertAnkiError, but let's be not
		{
			Name: "Connection error",
			Error: &ankiconnect.ConnectionError{
				Err: errors.New("some connection error"),
			},
			AssertError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				var connErr *ConnectionError
				return assert.ErrorAs(tt, err, &connErr)
			},
		},
		{
			Name: "Unknown error",
			Error: &ankiconnect.ServerError{
				Err: errors.New("some uknown error"),
			},
			AssertError: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, ErrUnknownServerError)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		for j := range methods {
			errorMethod := methods[j]
			t.Run(tc.Name+"-"+errorMethod, func(t *testing.T) {
				ankiClient := NewMockAnkiClient(t)
				for _, method := range methods {
					methodMeta := methodsParams[method]
					call := ankiClient.On(method, append([]any{mock.Anything}, methodMeta.Params...)...)
					if method == errorMethod {
						call.
							Return(nil, tc.Error).
							Once()
						break
					} else {
						call.Return(methodMeta.Return...).
							Once()
					}
				}
				client := newStatefullClientImpl(ankiClient, config, &statefullClientOptions{})
				actual := client.getNewState(context.Background())
				tc.AssertError(t, actual.LastError)
			})
		}
	}
}

func Test_statefullClient_run(t *testing.T) {
	t.Run("start in error state", func(t *testing.T) {
		ankiClient := NewMockAnkiClient(t)
		ankiClient.On("RequestPermission", mock.Anything).
			Return(nil, errors.New("first")).Once()

		afterChan := make(chan time.Time)
		client := newStatefullClientImpl(ankiClient, &Config{}, &statefullClientOptions{
			After: func(d time.Duration) <-chan time.Time {
				assert.Equal(t, StatefullClientErrorUpdateTimeout, d)
				return afterChan
			},
		})
		client.init()
		client.Stop()
		assert.ErrorContains(t, client.state.LastError, "first")
	})
	t.Run("start in normal state", func(t *testing.T) {
		ankiClient := NewMockAnkiClient(t)
		ankiClient.On("RequestPermission", mock.Anything).
			Return(&ankiconnect.RequestPermissionResponse{
				Permission: "granted",
			}, nil).Once()

		ankiClient.On("DeckNames", mock.Anything).
			Return(nil, nil).
			Once()
		ankiClient.On("ModelNames", mock.Anything).
			Return(nil, nil).
			Once()
		ankiClient.On("ModelFieldNames", mock.Anything, mock.Anything).
			Return(nil, nil).
			Once()
		afterChan := make(chan time.Time)
		client := newStatefullClientImpl(ankiClient, &Config{}, &statefullClientOptions{
			After: func(d time.Duration) <-chan time.Time {
				assert.Equal(t, StatefullClientDefaultUpdateTimeout, d)
				return afterChan
			},
		})
		client.init()
		client.Stop()
		assert.NoError(t, client.state.LastError)
	})
	t.Run("start in normal state and then go to error", func(t *testing.T) {
		ankiClient := NewMockAnkiClient(t)
		ankiClient.On("RequestPermission", mock.Anything).
			Return(&ankiconnect.RequestPermissionResponse{
				Permission: "granted",
			}, nil).Once()
		ankiClient.On("DeckNames", mock.Anything).
			Return(nil, nil).
			Once()
		ankiClient.On("ModelNames", mock.Anything).
			Return(nil, nil).
			Once()
		ankiClient.On("ModelFieldNames", mock.Anything, mock.Anything).
			Return(nil, nil).
			Once()
		// this is for second update
		ankiClient.On("RequestPermission", mock.Anything).
			Return(nil, errors.New("permerror")).
			Once()
		afterChan := make(chan time.Time)
		called := 0
		client := newStatefullClientImpl(ankiClient, &Config{}, &statefullClientOptions{
			After: func(d time.Duration) <-chan time.Time {
				called++
				if called == 1 {
					assert.Equal(t, StatefullClientDefaultUpdateTimeout, d)
				} else {
					assert.Equal(t, StatefullClientErrorUpdateTimeout, d)
				}
				return afterChan
			},
		})
		client.init()
		afterChan <- time.Now()
		client.Stop()
		assert.ErrorContains(t, client.state.LastError, "permerror")
		assert.Equal(t, 2, called)
	})
}

func Test_statefullClient_Config(t *testing.T) {
	expectedConfig := &Config{
		Addr:     "testaddr",
		APIKey:   "testapikey",
		Deck:     "testdeck",
		NoteType: "testnotetype",
		Mapping: map[string]*Template{
			"a": {
				Src: "testasrc",
			},
		},
	}
	client := newStatefullClientImpl(nil, expectedConfig, &statefullClientOptions{})
	assert.Equal(t, expectedConfig, client.Config())
}

func newTestStatefullClient(t *testing.T, config *Config, init func(client *MockAnkiClient)) (*statefullClient, *MockAnkiClient, chan time.Time) {
	ankiClient := NewMockAnkiClient(t)
	init(ankiClient)
	afterChan := make(chan time.Time)
	client := newStatefullClientImpl(ankiClient, config, &statefullClientOptions{
		After: func(d time.Duration) <-chan time.Time {
			return afterChan
		},
	})
	client.init()
	return client, ankiClient, afterChan
}

// newTestNormalStatefullClient is basically newTestStatefullClient but returns initialised client in normal (non-error) state
func newTestNormalStatefullClient(t *testing.T, config *Config) (*statefullClient, *MockAnkiClient, chan time.Time) {
	return newTestStatefullClient(t, config, func(client *MockAnkiClient) {
		client.On("RequestPermission", mock.Anything).
			Return(&ankiconnect.RequestPermissionResponse{
				Permission: "granted",
			}, nil).Once()
		client.On("DeckNames", mock.Anything).
			Return([]string{"deck1", "deck2"}, nil).
			Once()
		client.On("ModelNames", mock.Anything).
			Return([]string{"note1", "note2"}, nil).
			Once()
		client.On("ModelFieldNames", mock.Anything, mock.AnythingOfType("string")).
			Return([]string{"field1", "field2"}, nil).
			Maybe()
	})
}

// newTestNormalStatefullClient is basically newTestStatefullClient but returns initialised client in error state
func newTestErrorStatefullClient(t *testing.T, config *Config) (*statefullClient, *MockAnkiClient, chan time.Time) {
	return newTestStatefullClient(t, config, func(client *MockAnkiClient) {
		client.On("RequestPermission", mock.Anything).
			Return(&ankiconnect.RequestPermissionResponse{
				Permission: "denied",
			}, nil).Once()
	})
}

// in fact we test newStatefullClient, but without mock it is not possible
func Test_statefullClient_init(t *testing.T) {
	// we will test that client got first initial test, and that it running update cycle,
	// so RequestPermission will return two error, for initial getState and for update.
	client, _, updateCh := newTestStatefullClient(t, &Config{}, func(client *MockAnkiClient) {
		client.On("RequestPermission", mock.Anything).
			Return(nil, errors.New("first")).Once()
		client.On("RequestPermission", mock.Anything).
			Return(nil, errors.New("second")).Once()
	})
	assert.ErrorContains(t, client.state.LastError, "first")
	// now we will force update
	updateCh <- time.Now()
	// we don't know when update will be finished, but we can call Stop
	client.Stop()
	assert.ErrorContains(t, client.state.LastError, "second")
}

func Test_statefullClient_Stop(t *testing.T) {
	client, _, updateCh := newTestErrorStatefullClient(t, &Config{})
	client.Stop()
	// we can't gurantee that this will always detect problem, but this test will give only False Negative
	select {
	case updateCh <- time.Now():
		t.Fatalf("unexpected write to channel")
	default:
	}
}

func Test_statefullClient_withClient(t *testing.T) {
	// if we start from error state, we get previous error and can't update state
	t.Run("from error state", func(t *testing.T) {
		client, _, _ := newTestErrorStatefullClient(t, &Config{})
		err := client.withClient(func(client AnkiClient, config *Config, state *State) (*State, error) {
			t.Fatal("with client callback must not be called")
			return nil, nil
		})
		assert.ErrorIs(t, err, ErrForbiddenOrigin)
		client.Stop()
	})
	t.Run("normal state do nothing", func(t *testing.T) {
		client, _, _ := newTestNormalStatefullClient(t, &Config{
			Deck:     "deck1",
			NoteType: "note1",
		})
		err := client.withClient(func(client AnkiClient, config *Config, state *State) (*State, error) {
			return nil, nil
		})
		assert.NoError(t, err)
		client.Stop()
		assert.Equal(t, []string{"deck1", "deck2"}, client.state.Decks)
	})
	t.Run("normal state change state", func(t *testing.T) {
		client, _, _ := newTestNormalStatefullClient(t, &Config{
			Deck:     "deck1",
			NoteType: "note1",
		})
		expectedState := &State{
			AnkiState: AnkiState{
				Decks: []string{"mytestdeck"},
			},
		}
		err := client.withClient(func(client AnkiClient, config *Config, state *State) (*State, error) {
			return expectedState, nil
		})
		assert.NoError(t, err)
		client.Stop()
		assert.Equal(t, expectedState, client.state)
	})
	t.Run("normal state permanent error", func(t *testing.T) {
		client, _, _ := newTestNormalStatefullClient(t, &Config{
			Deck:     "deck1",
			NoteType: "note1",
		})
		baitState := &State{
			AnkiState: AnkiState{
				Decks: []string{"mytestdeck"},
			},
		}
		err := client.withClient(func(client AnkiClient, config *Config, state *State) (*State, error) {
			return baitState, &ankiconnect.ServerError{
				Err: ankiconnect.ErrCollectionUnavailable,
			}
		})
		assert.ErrorIs(t, err, ErrCollectionUnavailable)
		client.Stop()
		assert.ErrorIs(t, client.state.LastError, ErrCollectionUnavailable)
	})
	t.Run("normal state non permanent error", func(t *testing.T) {
		client, _, _ := newTestNormalStatefullClient(t, &Config{
			Deck:     "deck1",
			NoteType: "note1",
		})
		expectedState := &State{
			AnkiState: AnkiState{
				Decks: []string{"mytestdeck"},
			},
		}
		expectedError := errors.New("myerror")
		err := client.withClient(func(client AnkiClient, config *Config, state *State) (*State, error) {
			return expectedState, expectedError
		})
		assert.ErrorIs(t, err, expectedError)
		client.Stop()
		assert.Equal(t, expectedState, client.state)
	})
}

func Test_statefullClient_withClient_race(t *testing.T) {
	client, _, _ := newTestNormalStatefullClient(t, &Config{
		Deck:     "deck1",
		NoteType: "note1",
	})
	expectedState1 := &State{
		AnkiState: AnkiState{
			Decks: []string{"mytestdeck"},
		},
	}
	expectedState2 := &State{
		AnkiState: AnkiState{
			Decks: []string{"myotherdeck"},
		},
	}
	expectedStateCh := make(chan *State, 2)
	go func() {
		err := client.withClient(func(client AnkiClient, config *Config, state *State) (*State, error) {
			return expectedState1, nil
		})
		require.NoError(t, err)
		expectedStateCh <- expectedState1
	}()
	go func() {
		err := client.withClient(func(client AnkiClient, config *Config, state *State) (*State, error) {
			return expectedState2, nil
		})
		require.NoError(t, err)
		expectedStateCh <- expectedState2
	}()
	<-expectedStateCh
	expectedState := <-expectedStateCh
	client.Stop()
	assert.Equal(t, expectedState, client.state)
}

func Test_statefullClient_GetState(t *testing.T) {
	client, _, _ := newTestNormalStatefullClient(t, &Config{
		NoteType: "note1",
	})
	state, err := client.GetState(context.Background())
	client.Stop()
	assert.NoError(t, err)
	assert.Equal(t, &State{
		AnkiState: AnkiState{
			NoteTypes: []string{"note1", "note2"},
			NoteFields: map[string][]string{
				"note1": {"field1", "field2"},
			},
			Decks: []string{"deck1", "deck2"},
		},
		NoteTypeExists:   true,
		NoteHasAllFields: true,
		CurrentFields:    []string{"field1", "field2"},
	}, state)
}

func Test_statefullClient_CreateDeck(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		client, _, _ := newTestNormalStatefullClient(t, &Config{})
		err := client.CreateDeck(context.Background(), "\"hello")
		client.Stop()
		var validationErr *ValidationError
		assert.ErrorAs(t, err, &validationErr)
	})
	t.Run("deck already exists", func(t *testing.T) {
		client, _, _ := newTestNormalStatefullClient(t, &Config{})
		err := client.CreateDeck(context.Background(), "deck1")
		client.Stop()
		assert.ErrorIs(t, err, ErrDeckAlreadyExists)
	})
	t.Run("create failed", func(t *testing.T) {
		client, ankiClient, _ := newTestNormalStatefullClient(t, &Config{})
		ankiClient.On("CreateDeck", mock.Anything, "deck3").
			Return(int64(0), &ankiconnect.ServerError{
				Err: ankiconnect.ErrCollectionUnavailable,
			}).
			Once()
		err := client.CreateDeck(context.Background(), "deck3")
		client.Stop()
		assert.ErrorIs(t, err, ErrCollectionUnavailable)
	})
	t.Run("decknames failed", func(t *testing.T) {
		client, ankiClient, _ := newTestNormalStatefullClient(t, &Config{})
		ankiClient.On("CreateDeck", mock.Anything, "deck3").
			Return(int64(1), nil).
			Once()
		ankiClient.On("DeckNames", mock.Anything).
			Return(nil, &ankiconnect.ServerError{
				Err: ankiconnect.ErrCollectionUnavailable,
			}).
			Once()
		err := client.CreateDeck(context.Background(), "deck3")
		client.Stop()
		assert.ErrorIs(t, err, ErrCollectionUnavailable)
	})
	t.Run("ok", func(t *testing.T) {
		client, ankiClient, _ := newTestNormalStatefullClient(t, &Config{})
		ankiClient.On("CreateDeck", mock.Anything, "deck3").
			Return(int64(1), nil).
			Once()
		ankiClient.On("DeckNames", mock.Anything).
			Return([]string{"deck1", "deck2", "deck3"}, nil).
			Once()
		err := client.CreateDeck(context.Background(), "deck3")
		client.Stop()
		assert.NoError(t, err)
		assert.Equal(t, []string{"deck1", "deck2", "deck3"}, client.state.Decks)
		assert.False(t, client.state.DeckExists)
	})
	t.Run("ok state update", func(t *testing.T) {
		client, ankiClient, _ := newTestNormalStatefullClient(t, &Config{
			Deck: "deck3",
		})
		state, err := client.GetState(context.Background())
		require.NoError(t, err)
		assert.False(t, state.DeckExists)
		ankiClient.On("CreateDeck", mock.Anything, "deck3").
			Return(int64(1), nil).
			Once()
		ankiClient.On("DeckNames", mock.Anything).
			Return([]string{"deck1", "deck2", "deck3"}, nil).
			Once()
		err = client.CreateDeck(context.Background(), "deck3")
		client.Stop()
		assert.NoError(t, err)
		assert.Equal(t, []string{"deck1", "deck2", "deck3"}, client.state.Decks)
		assert.True(t, client.state.DeckExists)
	})
}

func Test_statefullClient_CreateDefaultNoteType(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		client, _, _ := newTestNormalStatefullClient(t, &Config{})
		err := client.CreateDefaultNoteType(context.Background(), "\"hello")
		client.Stop()
		var validationErr *ValidationError
		assert.ErrorAs(t, err, &validationErr)
	})
	t.Run("note type already exists", func(t *testing.T) {
		client, ankiClient, _ := newTestNormalStatefullClient(t, &Config{})
		ankiClient.On("CreateModel", mock.Anything, mock.Anything).
			Return(int64(0), &ankiconnect.ServerError{
				Message: "Model name already exists",
			})
		err := client.CreateDefaultNoteType(context.Background(), "note3")
		client.Stop()
		assert.ErrorIs(t, err, ErrNoteTypeAlreadyExists)
	})
	t.Run("create failed", func(t *testing.T) {
		client, ankiClient, _ := newTestNormalStatefullClient(t, &Config{})
		ankiClient.On("CreateModel", mock.Anything, mock.Anything).
			Return(int64(0), &ankiconnect.ServerError{
				Err: ankiconnect.ErrCollectionUnavailable,
			}).
			Once()
		err := client.CreateDefaultNoteType(context.Background(), "note3")
		client.Stop()
		assert.ErrorIs(t, err, ErrCollectionUnavailable)
	})
	t.Run("modelnames failed", func(t *testing.T) {
		client, ankiClient, _ := newTestNormalStatefullClient(t, &Config{})
		ankiClient.On("CreateModel", mock.Anything, mock.Anything).
			Return(int64(1), nil).
			Once()
		ankiClient.On("ModelNames", mock.Anything).
			Return(nil, &ankiconnect.ServerError{
				Err: ankiconnect.ErrCollectionUnavailable,
			}).
			Once()
		err := client.CreateDefaultNoteType(context.Background(), "note3")
		client.Stop()
		assert.ErrorIs(t, err, ErrCollectionUnavailable)
	})
	t.Run("ok", func(t *testing.T) {
		client, ankiClient, _ := newTestNormalStatefullClient(t, &Config{
			NoteType: "note2",
		})
		ankiClient.On("CreateModel", mock.Anything, mock.Anything).
			Return(int64(1), nil).
			Once()
		ankiClient.On("ModelNames", mock.Anything).
			Return([]string{"note1", "note2", "note3"}, nil).
			Once()
		err := client.CreateDefaultNoteType(context.Background(), "note3")
		client.Stop()
		assert.NoError(t, err)
		assert.Equal(t, []string{"note1", "note2", "note3"}, client.state.NoteTypes)
		assert.True(t, client.state.NoteTypeExists)
	})
	t.Run("ok note type dissapeared", func(t *testing.T) {
		client, ankiClient, _ := newTestNormalStatefullClient(t, &Config{
			NoteType: "note2",
		})
		ankiClient.On("CreateModel", mock.Anything, mock.Anything).
			Return(int64(1), nil).
			Once()
		ankiClient.On("ModelNames", mock.Anything).
			Return([]string{"note1", "note3"}, nil).
			Once()
		err := client.CreateDefaultNoteType(context.Background(), "note3")
		client.Stop()
		assert.NoError(t, err)
		assert.Equal(t, []string{"note1", "note3"}, client.state.NoteTypes)
		assert.False(t, client.state.NoteTypeExists)
	})
	t.Run("ok note type appeared", func(t *testing.T) {
		client, ankiClient, _ := newTestNormalStatefullClient(t, &Config{
			NoteType: "note3",
		})
		ankiClient.On("CreateModel", mock.Anything, mock.Anything).
			Return(int64(1), nil).
			Once()
		ankiClient.On("ModelNames", mock.Anything).
			Return([]string{"note1", "note2", "note3"}, nil).
			Once()
		err := client.CreateDefaultNoteType(context.Background(), "note3")
		client.Stop()
		assert.NoError(t, err)
		assert.Equal(t, []string{"note1", "note2", "note3"}, client.state.NoteTypes)
		assert.True(t, client.state.NoteTypeExists)
	})
}
