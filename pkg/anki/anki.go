package anki

import (
	"context"
	"sync"

	"golang.org/x/exp/slices"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
)

type Anki struct {
	constructor ClientConstructorFn

	mu      sync.Mutex
	wrapper *ankiWrapper
}

type ankiWrapper struct {
	client AnkiClient
	config *Config
}

//go:generate $MOCKERY_TOOL --name AnkiClient --testonly=true --inpackage=true
type AnkiClient interface {
	RequestPermission(ctx context.Context) (*ankiconnect.RequestPermissionResponse, error)
	DeckNames(ctx context.Context) ([]string, error)
	ModelNames(ctx context.Context) ([]string, error)
	ModelFieldNames(ctx context.Context, modelName string) ([]string, error)
}

type ClientConstructorFn func(*ankiconnect.Options) (AnkiClient, error)

func DefaultClientConstructor(o *ankiconnect.Options) (AnkiClient, error) {
	return ankiconnect.New(o)
}

// NewAnki return uninitialized Anki instance.
// It should be inited before use with ReloadConfig.
// It is intended to be used with ConfigReloader.
func NewAnki(constuctor ClientConstructorFn) *Anki {
	return &Anki{
		constructor: constuctor,
	}
}

func (a *Anki) ReloadConfig(config *Config) error {
	client, err := a.constructor(config.options())
	if err != nil {
		return err
	}
	wrapper := &ankiWrapper{
		client: client,
		config: config,
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.wrapper = wrapper
	return nil
}

type StateResult struct {
	Version int

	PermissionGranted bool
	APIKeyRequired    bool

	DeckExists        bool
	NoteTypeExists    bool
	NoteMissingFields []string
}

// FullStateCheck checks that anki is available, decks and note types exists, also FieldMapping is possible
func (a *Anki) FullStateCheck(ctx context.Context) (*StateResult, error) {
	wrapper := a.getWrapper()
	result := &StateResult{}
	permissions, err := wrapper.client.RequestPermission(ctx)
	if err != nil {
		return nil, err
	}
	result.PermissionGranted = permissions.Permission == ankiconnect.PermissionGranted
	result.Version = permissions.Version
	result.APIKeyRequired = permissions.RequireAPIKey
	if !result.PermissionGranted {
		return result, nil
	}

	decks, err := wrapper.client.DeckNames(ctx)
	if err != nil {
		return result, err
	}
	deckExists := slices.ContainsFunc(decks, func(e string) bool { return e == wrapper.config.Deck })
	result.DeckExists = deckExists

	noteTypes, err := wrapper.client.ModelNames(ctx)
	if err != nil {
		return result, err
	}
	noteTypeExists := slices.ContainsFunc(noteTypes, func(e string) bool { return e == wrapper.config.NoteType })
	result.NoteTypeExists = noteTypeExists
	if noteTypeExists {
		noteFields, err := wrapper.client.ModelFieldNames(ctx, wrapper.config.NoteType)
		if err != nil {
			return result, err
		}
		setFields := map[string]struct{}{}
		for _, field := range noteFields {
			setFields[field] = struct{}{}
		}
		for field := range wrapper.config.Mapping {
			_, ok := setFields[field]
			if !ok {
				result.NoteMissingFields = append(result.NoteMissingFields, field)
			}
		}
	}
	return result, nil
}

func (a *Anki) getWrapper() *ankiWrapper {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.wrapper
}
