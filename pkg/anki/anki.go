package anki

import (
	"bytes"
	"context"
	"sync"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
	"github.com/Darkclainer/japwords/pkg/anki/query"
	"github.com/Darkclainer/japwords/pkg/lemma"
)

//go:generate $MOCKERY_TOOL --name StatefullClient --testonly=true --inpackage=true
type StatefullClient interface {
	Stop()
	Config() *Config
	GetState(ctx context.Context) (*State, error)
	CreateDeck(ctx context.Context, name string) error
	CreateDefaultNoteType(ctx context.Context, name string) error
	AddNote(ctx context.Context, note *AddNoteRequest) error
	QueryNotes(ctx context.Context, query string) ([]*ankiconnect.NoteInfo, error)
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
	OrderDefined     bool
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
		OrderDefined:     state.OrderDefined,
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

type AddNoteField struct {
	Name  string
	Value string
}

type AddNoteRequest struct {
	Fields   []AddNoteField
	Tags     []string
	AudioURL string
}

func (a *Anki) PrepareProjectedLemma(ctx context.Context, lemma *lemma.ProjectedLemma) (*AddNoteRequest, error) {
	client := a.getClient()
	state, err := client.GetState(ctx)
	if err != nil {
		return nil, err
	}
	if !state.IsReadyToAddNote() {
		return nil, ErrIncompleteConfiguration
	}
	config := client.Config()
	fields := make([]AddNoteField, len(state.CurrentFields))
	var buffer bytes.Buffer
	for i, fieldName := range state.CurrentFields {
		fields[i].Name = fieldName
		fieldTemplate, ok := config.Mapping[fieldName]
		if !ok {
			continue
		}
		buffer.Reset()
		err := fieldTemplate.Tmpl.Execute(&buffer, lemma)
		if err != nil {
			// Probably best to leave as unexported error
			return nil, err
		}
		fields[i].Value = buffer.String()
	}
	return &AddNoteRequest{
		Fields: fields,
		// TODO: add audio and tags
	}, nil
}

// SearchProjectedLemmas returns note id for specified lemmas, empty string means not found.
// Note id represented as string, because it's not fully documented what id actually means
func (a *Anki) SearchProjectedLemmas(ctx context.Context, lemmas []*lemma.ProjectedLemma) ([]int64, error) {
	client := a.getClient()
	state, err := client.GetState(ctx)
	if err != nil {
		return nil, err
	}
	if !state.IsReadyToAddNote() {
		return nil, ErrIncompleteConfiguration
	}
	// actually checks bellow redudant because "readiness" should include it
	if len(state.CurrentFields) < 1 {
		return nil, ErrIncompleteConfiguration
	}
	orderField := state.CurrentFields[0]
	searchQuery, orderValues, err := generateQueryForNotes(lemmas, orderField, client.Config())
	if err != nil {
		return nil, err
	}
	notes, err := client.QueryNotes(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	// because query return notes in no particular order, we need to rebuild result
	return confirmFoundNotes(notes, orderField, orderValues), nil
}

// generateQueryForNotes returns search query for anki, slice of values of expected values of order field and errors
func generateQueryForNotes(lemmas []*lemma.ProjectedLemma, orderField string, config *Config) (string, []string, error) {
	orderTemplate, ok := config.Mapping[orderField]
	if !ok {
		return "", nil, ErrIncompleteConfiguration
	}
	var buffer bytes.Buffer
	// what values in orderField we will search
	fieldQueries := make([]query.Query, len(lemmas))
	// we associate note with value in order field, this list will help
	// us understand what notes anki actually has
	orderValues := make([]string, len(lemmas))
	for i, lemma := range lemmas {
		buffer.Reset()
		err := orderTemplate.Tmpl.Execute(&buffer, lemma)
		if err != nil {
			return "", nil, err
		}
		v := buffer.String()
		fieldQueries[i] = query.Exact(orderField, v)
		orderValues[i] = v
	}
	searchQuery := query.And(
		// this must be exactly duplication settings in add note function
		query.Exact("deck", config.Deck),
		query.Exact("note", config.NoteType),
		// search for any fields
		query.Or(fieldQueries...),
	)
	return query.Render(searchQuery), orderValues, nil
}

// confirmFoundNotes returns actually found notes using expected values of order field
func confirmFoundNotes(notes []*ankiconnect.NoteInfo, orderField string, orderValues []string) []int64 {
	foundIds := make([]int64, len(orderValues))
	// orderField value to noteId
	foundNotes := make(map[string]int64, len(orderValues))
	for _, note := range notes {
		field, ok := note.Fields[orderField]
		if !ok {
			// redudant
			continue
		}
		foundNotes[field.Value] = note.NoteID
	}
	for i, orderValue := range orderValues {
		id, ok := foundNotes[orderValue]
		if ok {
			foundIds[i] = id
		}
	}
	return foundIds
}

func (a *Anki) AddNote(ctx context.Context, note *AddNoteRequest) error {
	return a.getClient().AddNote(ctx, note)
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
