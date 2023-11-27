package anki

import (
	"context"
	"sync"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
)

//go:generate $MOCKERY_TOOL --name StatefullClient --testonly=true --inpackage=true
type StatefullClient interface {
	Stop()
	Config() *Config
	GetState(ctx context.Context) (*State, error)
	CreateDeck(ctx context.Context, name string) error
	CreateDefaultNoteType(ctx context.Context, name string) error
}

type StatefullClientConstructorFn func(*Config) (StatefullClient, error)

func DefaultStatefullClientConstructor(conf *Config) (StatefullClient, error) {
	statelessClient, err := ankiconnect.New(conf.options())
	if err != nil {
		return nil, err
	}
	client := newStatefullClient(statelessClient, conf)
	return client, nil
}

// Anki is wrapper the main purpose is to support config reloading
type Anki struct {
	constructor StatefullClientConstructorFn

	mu     sync.Mutex
	client StatefullClient
}

// NewAnki return uninitialized Anki instance.
// It should be inited before use with ReloadConfig.
// It is intended to be used with ConfigReloader.
func NewAnki(constuctor StatefullClientConstructorFn) *Anki {
	return &Anki{
		constructor: constuctor,
	}
}

// ReloadConfig intialize internal client with config.
func (a *Anki) ReloadConfig(config *Config) error {
	statefullClient, err := a.constructor(config)
	if err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.client != nil {
		// Stop here doesn't invalidate client, it is fine if someone use it right now
		a.client.Stop()
	}
	a.client = statefullClient
	return nil
}

type StateResult struct {
	Version int

	DeckExists       bool
	NoteTypeExists   bool
	NoteHasAllFields bool
}

// FullStateCheck checks that anki is available, decks and note types exists, also FieldMapping is possible
func (a *Anki) FullStateCheck(ctx context.Context) (*StateResult, error) {
	state, err := a.getClient().GetState(ctx)
	if err != nil {
		return nil, err
	}
	return &StateResult{
		Version:          state.Version,
		DeckExists:       state.DeckExists,
		NoteTypeExists:   state.NoteTypeExists,
		NoteHasAllFields: state.NoteHasAllFields,
	}, nil
}

func (a *Anki) Decks(ctx context.Context) ([]string, error) {
	state, err := a.getClient().GetState(ctx)
	if err != nil {
		return nil, err
	}
	return state.Decks, nil
}

func (a *Anki) NoteTypes(ctx context.Context) ([]string, error) {
	state, err := a.getClient().GetState(ctx)
	if err != nil {
		return nil, err
	}
	return state.NoteTypes, nil
}

// TODO: remove name parameter
func (a *Anki) NoteTypeFields(ctx context.Context) ([]string, error) {
	state, err := a.getClient().GetState(ctx)
	if err != nil {
		return nil, err
	}
	if !state.NoteTypeExists {
		return nil, ErrNoteTypeNotExists
	}
	return state.CurrentFields, nil
}

func (a *Anki) CreateDeck(ctx context.Context, name string) error {
	return a.getClient().CreateDeck(ctx, name)
}

func (a *Anki) CreateDefaultNote(ctx context.Context, name string) error {
	return a.getClient().CreateDefaultNoteType(ctx, name)
}

func (a *Anki) Stop() {
	a.mu.Lock()
	a.client.Stop()
	a.mu.Unlock()
}

func (a *Anki) getClient() StatefullClient {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.client
}
