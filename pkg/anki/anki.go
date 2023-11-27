package anki

import (
	"context"
	"errors"
	"slices"
	"sync"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
)

var (
	// redefine this errors from client, because we probably want to expose them as API
	ErrForbiddenOrigin       = errors.New("anki-connect forbid request from client origin")
	ErrInvalidAPIKey         = errors.New("anki-connect rejected request because api key is invalid")
	ErrCollectionUnavailable = errors.New("anki-connect is not ready for specified action")

	ErrDeckAlreadyExists = errors.New("deck with the same name already exists")
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

	CreateDeck(ctx context.Context, name string) (int64, error)
	CreateModel(ctx context.Context, parameters *ankiconnect.CreateModelRequest) (int64, error)
	AddNote(ctx context.Context, params *ankiconnect.AddNoteParams, opts *ankiconnect.AddNoteOptions) (int64, error)
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
	Version int

	DeckExists       bool
	NoteTypeExists   bool
	NoteHasAllFields bool
}

// FullStateCheck checks that anki is available, decks and note types exists, also FieldMapping is possible
func (a *Anki) FullStateCheck(ctx context.Context) (*StateResult, error) {
	client, config := a.getClient()
	result := &StateResult{}
	permissions, err := client.RequestPermission(ctx)
	if err != nil {
		return result, redefineClientError(err)
	}
	result.Version = permissions.Version
	if permissions.Permission != ankiconnect.PermissionGranted {
		return result, ErrForbiddenOrigin
	}
	decks, err := client.DeckNames(ctx)
	if err != nil {
		return result, redefineClientError(err)
	}
	deckExists := slices.ContainsFunc(decks, func(e string) bool { return e == config.Deck })
	result.DeckExists = deckExists
	noteTypes, err := client.ModelNames(ctx)
	if err != nil {
		return result, redefineClientError(err)
	}
	noteTypeExists := slices.ContainsFunc(noteTypes, func(e string) bool { return e == config.NoteType })
	result.NoteTypeExists = noteTypeExists
	if noteTypeExists {
		noteFields, err := client.ModelFieldNames(ctx, config.NoteType)
		if err != nil {
			return result, redefineClientError(err)
		}
		setFields := map[string]struct{}{}
		for _, field := range noteFields {
			setFields[field] = struct{}{}
		}
		hasAllFieds := true
		for field := range config.Mapping {
			_, ok := setFields[field]
			if !ok {
				hasAllFieds = false
				break
			}
		}
		result.NoteHasAllFields = hasAllFieds
	}
	return result, nil
}

func (a *Anki) Decks(ctx context.Context) ([]string, error) {
	client, _ := a.getClient()
	decks, err := client.DeckNames(ctx)
	if err != nil {
		return nil, err
	}
	return decks, nil
}

func (a *Anki) NoteTypes(ctx context.Context) ([]string, error) {
	client, _ := a.getClient()
	noteTypes, err := client.ModelNames(ctx)
	if err != nil {
		return nil, err
	}
	return noteTypes, nil
}

func (a *Anki) NoteTypeFields(ctx context.Context, name string) ([]string, error) {
	client, _ := a.getClient()
	fields, err := client.ModelFieldNames(ctx, name)
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func (a *Anki) CreateDeck(ctx context.Context, name string) error {
	if err := validateDeckName(name); err != nil {
		return &ValidationError{Msg: err.Error()}
	}
	client, _ := a.getClient()
	decks, err := client.DeckNames(ctx)
	if err != nil {
		return err
	}
	if slices.Contains(decks, name) {
		// actually Anki-Connect doesn't care if there already is deck with the same name, it just do nothing
		return ErrDeckAlreadyExists
	}
	_, err = client.CreateDeck(ctx, name)
	return err
}

func (a *Anki) CreateDefaultNote(ctx context.Context, name string) error {
	if err := validateNoteType(name); err != nil {
		return &ValidationError{Msg: err.Error()}
	}
	client, _ := a.getClient()
	modelRequest := defaultCreateModelRequest()
	modelRequest.ModelName = name
	_, err := client.CreateModel(ctx, modelRequest)
	return err
}

func (a *Anki) getClient() (AnkiClient, *Config) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.client, a.config
}

func redefineClientError(err error) error {
	var serverError *ankiconnect.ServerError
	if !errors.As(err, &serverError) {
		return err
	}
	switch serverError.Err {
	case ankiconnect.ErrCollectionUnavailable:
		return ErrCollectionUnavailable
	case ankiconnect.ErrInvalidAPIKey:
		return ErrInvalidAPIKey
	default:
		return err
	}
}
