package anki

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
)

//go:generate $MOCKERY_TOOL --name AnkiClient --testonly=true --inpackage=true
type AnkiClient interface {
	RequestPermission(ctx context.Context) (*ankiconnect.RequestPermissionResponse, error)
	DeckNames(ctx context.Context) ([]string, error)
	ModelNames(ctx context.Context) ([]string, error)
	ModelFieldNames(ctx context.Context, modelName string) ([]string, error)
	FindNotes(ctx context.Context, query string) ([]int64, error)
	NotesInfo(ctx context.Context, ids []int64) ([]*ankiconnect.NoteInfo, error)

	CreateDeck(ctx context.Context, name string) (int64, error)
	CreateModel(ctx context.Context, parameters *ankiconnect.CreateModelRequest) (int64, error)
	AddNote(ctx context.Context, params *ankiconnect.AddNoteParams, opts *ankiconnect.AddNoteOptions) (int64, error)
}

type AnkiState struct {
	Version   int
	Decks     []string
	NoteTypes []string
	// NoteFields is dictionary between NoteType name and known fields
	// expect that this completed only for config.NoteType
	NoteFields map[string][]string
}

type State struct {
	AnkiState

	LastError        error
	CurrentFields    []string
	DeckExists       bool
	NoteTypeExists   bool
	NoteHasAllFields bool
	OrderDefined     bool
	AudioFieldExists bool
}

func (s *State) IsReadyToAddNote() bool {
	return s.LastError == nil && s.DeckExists && s.NoteTypeExists && s.NoteHasAllFields && s.OrderDefined && s.AudioFieldExists
}

func (state *State) updateFromAnkiState(config *Config) {
	state.DeckExists = slices.ContainsFunc(state.Decks, func(e string) bool { return e == config.Deck })
	// this is redudant, but necessary, because NoteFields is more like all already known fields and can be outdated
	noteTypeExists := slices.ContainsFunc(state.NoteTypes, func(e string) bool { return e == config.NoteType })
	var fields []string
	if noteTypeExists {
		fields, noteTypeExists = state.NoteFields[config.NoteType]
	}
	state.NoteTypeExists = noteTypeExists
	if noteTypeExists {
		state.CurrentFields = fields
		fieldsSet := map[string]struct{}{}
		for _, field := range fields {
			fieldsSet[field] = struct{}{}
		}
		hasAllFieds := true
		for field := range config.Mapping {
			_, ok := fieldsSet[field]
			if !ok {
				hasAllFieds = false
				break
			}
		}
		if len(fields) > 0 {
			_, ok := config.Mapping[fields[0]]
			state.OrderDefined = ok
		}
		state.NoteHasAllFields = hasAllFieds

		if config.AudioField == "" {
			// if audio field is not defined, it's ok
			state.AudioFieldExists = true
		} else {
			_, ok := fieldsSet[config.AudioField]
			state.AudioFieldExists = ok
		}
	}
}

const (
	StatefullClientErrorUpdateTimeout   = time.Second * 3
	StatefullClientDefaultUpdateTimeout = time.Second * 10
)

func (s *State) nextUpdateTimeout() time.Duration {
	if s.LastError != nil {
		return StatefullClientErrorUpdateTimeout
	}
	return StatefullClientDefaultUpdateTimeout
}

type statefullClient struct {
	exited            chan struct{}
	exitContext       context.Context
	exitContextCancel context.CancelFunc
	config            *Config

	// we can be more optimistic with mutex and lock only state change, but it's not
	// like we have any problem with throughput when there is only one user
	mu     sync.Mutex
	client AnkiClient
	state  *State

	// after is for testing only, in production it is time.After
	after func(time.Duration) <-chan time.Time
}

func newStatefullClient(client AnkiClient, config *Config) *statefullClient {
	statefullClient := newStatefullClientImpl(client, config, &statefullClientOptions{
		After: time.After,
	})
	statefullClient.init()
	return statefullClient
}

type statefullClientOptions struct {
	After func(d time.Duration) <-chan time.Time
}

// newStatefullClientImpl can be used to mock time.After for tests
func newStatefullClientImpl(client AnkiClient, config *Config, opts *statefullClientOptions) *statefullClient {
	exitContext, exitContextCancel := context.WithCancel(context.Background())
	return &statefullClient{
		client:            client,
		config:            config,
		exited:            make(chan struct{}),
		exitContext:       exitContext,
		exitContextCancel: exitContextCancel,
		after:             opts.After,
	}
}

func (sc *statefullClient) init() {
	sc.state = sc.getNewState(sc.exitContext)
	go sc.run()
}

func (sc *statefullClient) run() {
	sc.mu.Lock()
	sleepTimeout := sc.state.nextUpdateTimeout()
	sc.mu.Unlock()
	for {
		select {
		case <-sc.exitContext.Done():
			close(sc.exited)
			return
		case <-sc.after(sleepTimeout):
			sc.mu.Lock()
			ctx, cancel := context.WithCancel(sc.exitContext)
			newState := sc.getNewState(ctx)
			cancel()
			sleepTimeout = newState.nextUpdateTimeout()
			sc.state = newState
			sc.mu.Unlock()
		}
	}
}

func (sc *statefullClient) withClient(fn func(client AnkiClient, config *Config, state *State) (*State, error)) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	// early exited if our last state is error
	if sc.state.LastError != nil {
		return sc.state.LastError
	}
	stateCopy := (*sc.state)
	newState, err := fn(sc.client, sc.config, &stateCopy)
	err, isAnkiError := convertAnkiError(err)
	if isAnkiError {
		// new state is error state
		newState = &State{
			LastError: err,
		}
	}
	// update new state if get one from error or not
	if newState != nil {
		sc.state = newState
	}
	return err
}

func (sc *statefullClient) Config() *Config {
	return sc.config
}

func (sc *statefullClient) Stop() {
	sc.exitContextCancel()
	<-sc.exited
}

func (sc *statefullClient) getNewState(ctx context.Context) *State {
	client, config := sc.client, sc.config
	state := &State{
		AnkiState: AnkiState{
			NoteFields: map[string][]string{},
		},
	}
	errorState := func(err error) *State {
		convertedErr, _ := convertAnkiError(err)
		return &State{
			LastError: convertedErr,
		}
	}

	permissions, err := client.RequestPermission(ctx)
	if err != nil {
		return errorState(err)
	}
	state.Version = permissions.Version
	if permissions.Permission != ankiconnect.PermissionGranted {
		return errorState(ErrForbiddenOrigin)
	}
	decks, err := client.DeckNames(ctx)
	if err != nil {
		return errorState(err)
	}
	state.Decks = decks
	noteTypes, err := client.ModelNames(ctx)
	if err != nil {
		return errorState(err)
	}
	state.NoteTypes = noteTypes
	noteFields, err := client.ModelFieldNames(ctx, config.NoteType)
	if err != nil {
		var serverError *ankiconnect.ServerError
		if !errors.As(err, &serverError) || !strings.HasPrefix(serverError.Message, "model was not found:") {
			return errorState(err)
		}
	} else {
		state.NoteFields[config.NoteType] = noteFields
	}
	state.updateFromAnkiState(config)
	return state
}

func (sc *statefullClient) GetState(ctx context.Context) (*State, error) {
	var currentState *State
	err := sc.withClient(func(_ AnkiClient, _ *Config, state *State) (*State, error) {
		currentState = state
		return nil, nil
	})
	return currentState, err
}

func (sc *statefullClient) CreateDeck(ctx context.Context, name string) error {
	if err := validateDeckName(name); err != nil {
		return &ValidationError{Msg: err.Error()}
	}
	err := sc.withClient(func(client AnkiClient, config *Config, state *State) (*State, error) {
		if slices.Contains(state.Decks, name) {
			// actually Anki-Connect doesn't care if there already is a deck with the same name, it just do nothing
			return nil, ErrDeckAlreadyExists
		}
		_, err := client.CreateDeck(ctx, name)
		if err != nil {
			return nil, err
		}
		decks, err := client.DeckNames(ctx)
		if err != nil {
			return nil, err
		}
		state.Decks = decks
		state.updateFromAnkiState(config)
		return state, nil
	})
	return err
}

func (sc *statefullClient) CreateDefaultNoteType(ctx context.Context, name string) error {
	if err := validateNoteType(name); err != nil {
		return &ValidationError{Msg: err.Error()}
	}
	err := sc.withClient(func(client AnkiClient, config *Config, state *State) (*State, error) {
		modelRequest := defaultCreateModelRequest()
		modelRequest.ModelName = name
		_, err := client.CreateModel(ctx, modelRequest)
		if err != nil {
			var serverError *ankiconnect.ServerError
			if errors.As(err, &serverError) {
				if serverError.Message == "Model name already exists" {
					return nil, ErrNoteTypeAlreadyExists
				}
			}
			return nil, err
		}
		noteTypes, err := client.ModelNames(ctx)
		if err != nil {
			return nil, err
		}
		state.NoteTypes = noteTypes
		state.NoteFields[name] = modelRequest.Fields
		state.updateFromAnkiState(config)
		return state, nil
	})
	return err
}

func (sc *statefullClient) AddNote(ctx context.Context, note *AddNoteRequest) (int64, error) {
	var noteID int64
	err := sc.withClient(func(client AnkiClient, config *Config, state *State) (*State, error) {
		if !state.IsReadyToAddNote() {
			return nil, ErrIncompleteConfiguration
		}
		// NOTE: we could assert request on known state, but why would we?
		fields := make(map[string]string, len(note.Fields))
		for i := range note.Fields {
			fields[note.Fields[i].Name] = note.Fields[i].Value
		}
		assets, err := convertAddNoteAudioAssets(note.AudioAssets)
		if err != nil {
			return nil, err
		}

		noteID, err = client.AddNote(ctx,
			&ankiconnect.AddNoteParams{
				Fields: fields,
				Assets: assets,
				// TODO:
				// Tags:   note.Tags,
			},
			&ankiconnect.AddNoteOptions{
				Deck:           config.Deck,
				Model:          config.NoteType,
				DuplicateScope: ankiconnect.DuplicateScopeDeck,
				DuplicateFlags: ankiconnect.DuplicateFlagsCheck,
			},
		)
		var serverError *ankiconnect.ServerError
		if errors.As(err, &serverError) && serverError.Message == "cannot create note because it is a duplicate" {
			return nil, ErrDuplicatedNoteFound
		}
		return nil, err
	})
	if err != nil {
		return 0, err
	}
	return noteID, nil
}

func convertAddNoteAudioAssets(noteAssets []AddNoteAudioAsset) ([]*ankiconnect.AddNoteAsset, error) {
	var assets []*ankiconnect.AddNoteAsset
	for _, asset := range noteAssets {
		mediaAsset := &ankiconnect.MediaAssetRequest{
			Filename: asset.Filename,
			Fields: []string{
				asset.Field,
			},
		}
		switch {
		case asset.Data != "":
			md5Hash, err := md5OfEncodedData(asset.Data)
			if err != nil {
				return nil, err
			}
			mediaAsset.Data = asset.Data
			mediaAsset.SkipHash = md5Hash
		case asset.URL != "":
			mediaAsset.URL = asset.URL
		}
		assets = append(assets, &ankiconnect.AddNoteAsset{
			Asset: *mediaAsset,
			Type:  ankiconnect.MediaTypeAudio,
		})

	}
	return assets, nil
}

func md5OfEncodedData(src string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = hash.Write(decoded)
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// QueryNotes get notes by query, it does FindNotes and NoteInfo.
func (sc *statefullClient) QueryNotes(ctx context.Context, query string) ([]*ankiconnect.NoteInfo, error) {
	var notes []*ankiconnect.NoteInfo
	err := sc.withClient(func(client AnkiClient, _ *Config, _ *State) (*State, error) {
		noteIds, err := client.FindNotes(ctx, query)
		if err != nil {
			return nil, err
		}
		notes, err = client.NotesInfo(ctx, noteIds)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return notes, nil
}
