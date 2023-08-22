package anki

import (
	"context"
	"slices"
	"sync"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
)

type Anki struct {
	constructor ClientConstructorFn

	mu     sync.Mutex
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

// ReloadConfig intialize internal client with config.
func (a *Anki) ReloadConfig(config *Config) error {
	client, err := a.constructor(config.options())
	if err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.client = client
	a.config = config
	return nil
}

type StateResult struct {
	Connected bool
	Version   int

	PermissionGranted bool
	APIKeyRequired    bool

	DeckExists        bool
	NoteTypeExists    bool
	NoteMissingFields []string
}

// FullStateCheck checks that anki is available, decks and note types exists, also FieldMapping is possible
func (a *Anki) FullStateCheck(ctx context.Context) (*StateResult, error) {
	client, config := a.getClient()
	result := &StateResult{}
	permissions, err := client.RequestPermission(ctx)
	// TODO: deal with connection errors?
	if err != nil {
		return result, err
	}
	result.Connected = true
	result.PermissionGranted = permissions.Permission == ankiconnect.PermissionGranted
	result.Version = permissions.Version
	result.APIKeyRequired = permissions.RequireAPIKey
	if !result.PermissionGranted {
		return result, nil
	}

	decks, err := client.DeckNames(ctx)
	if err != nil {
		return result, err
	}
	deckExists := slices.ContainsFunc(decks, func(e string) bool { return e == config.Deck })
	result.DeckExists = deckExists

	noteTypes, err := client.ModelNames(ctx)
	if err != nil {
		return result, err
	}
	noteTypeExists := slices.ContainsFunc(noteTypes, func(e string) bool { return e == config.NoteType })
	result.NoteTypeExists = noteTypeExists
	if noteTypeExists {
		noteFields, err := client.ModelFieldNames(ctx, config.NoteType)
		if err != nil {
			return result, err
		}
		setFields := map[string]struct{}{}
		for _, field := range noteFields {
			setFields[field] = struct{}{}
		}
		for field := range config.Mapping {
			_, ok := setFields[field]
			if !ok {
				result.NoteMissingFields = append(result.NoteMissingFields, field)
			}
		}
	}
	return result, nil
}

func (a *Anki) getClient() (AnkiClient, *Config) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.client, a.config
}
