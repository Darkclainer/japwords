package anki

import (
	"bytes"
	"context"
	"strings"
	"sync"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
	"github.com/Darkclainer/japwords/pkg/anki/query"
	"github.com/Darkclainer/japwords/pkg/lemma"
	"github.com/Darkclainer/japwords/pkg/mediatypes"
)

//go:generate $MOCKERY_TOOL --name StatefullClient --testonly=true --inpackage=true
type StatefullClient interface {
	Stop()
	Config() *Config
	GetState(ctx context.Context) (*State, error)
	CreateDeck(ctx context.Context, name string) error
	CreateDefaultNoteType(ctx context.Context, name string) error
	AddNote(ctx context.Context, note *AddNoteRequest) (int64, error)
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

type AddNoteAudioAsset struct {
	Field    string
	Filename string
	URL      string
	// Data is base64 encoded audio
	Data string
}

type AddNoteRequest struct {
	Fields []AddNoteField
	Tags   []string
	// AudioAssets is list of possible audio assets. Choice must be made between
	// assets that have same Field value
	AudioAssets []AddNoteAudioAsset
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
	fields, err := prepareFieldsForNoteRequest(lemma, state.CurrentFields, config.Mapping)
	if err != nil {
		// Probably best to leave as unexported error
		return nil, err
	}
	audioAssets := prepareAudiosForNoteRequest(lemma, config)
	return &AddNoteRequest{
		Fields:      fields,
		AudioAssets: audioAssets,
		// TODO: add tags
	}, nil
}

func prepareFieldsForNoteRequest(lemma *lemma.ProjectedLemma, currentFields []string, mapping TemplateMapping) ([]AddNoteField, error) {
	fields := make([]AddNoteField, len(currentFields))
	var buffer bytes.Buffer
	for i, fieldName := range currentFields {
		fields[i].Name = fieldName
		fieldTemplate, ok := mapping[fieldName]
		if !ok {
			continue
		}
		buffer.Reset()
		err := fieldTemplate.Tmpl.Execute(&buffer, lemma)
		if err != nil {
			return nil, err
		}
		fields[i].Value = buffer.String()
	}
	return fields, nil
}

func prepareAudiosForNoteRequest(lemma *lemma.ProjectedLemma, config *Config) []AddNoteAudioAsset {
	if config.AudioField == "" {
		return nil
	}
	audios := sortLemmaAudios(lemma.Audio, config.AudioPreferredType)
	result := make([]AddNoteAudioAsset, len(audios))
	baseFilename := generateLemmaAudioBasename(lemma)
	for i := range audios {
		extension := mediatypes.GetExtensionByMediaType(audios[i].MediaType)
		result[i] = AddNoteAudioAsset{
			Field:    config.AudioField,
			Filename: baseFilename + extension,
			URL:      audios[i].Source,
		}
	}
	return result
}

func generateLemmaAudioBasename(lemma *lemma.ProjectedLemma) string {
	var buffer strings.Builder
	_, _ = buffer.WriteString(lemma.Slug.Word)
	if lemma.Slug.Hiragana != "" {
		_ = buffer.WriteByte('-')
		_, _ = buffer.WriteString(lemma.Slug.Hiragana)
	}
	return buffer.String()
}

func sortLemmaAudios(audios []lemma.Audio, preferredType string) []lemma.Audio {
	// we want stable sort, so this is more complicated then could be
	result := make([]lemma.Audio, len(audios))
	j := 0
	for i := range audios {
		if strings.Contains(audios[i].MediaType, preferredType) {
			result[j] = audios[i]
			j++
		}
	}
	for i := range audios {
		if !strings.Contains(audios[i].MediaType, preferredType) {
			result[j] = audios[i]
			j++
		}
	}
	return result
}

// AddNote sends request to anki-connect to add specified note.
// If there are assets with equal Field name, then only first of these asset is saved.
func (a *Anki) AddNote(ctx context.Context, note *AddNoteRequest) (NoteID, error) {
	// create copy with filtered out assets that has duplicated field
	noteCopy := *note
	var assets []AddNoteAudioAsset
	usedFields := map[string]struct{}{}
	for _, asset := range note.AudioAssets {
		_, ok := usedFields[asset.Field]
		if ok {
			continue
		}
		usedFields[asset.Field] = struct{}{}
		assets = append(assets, asset)
	}
	noteCopy.AudioAssets = assets
	// TODO: we can also predownload resource ourselves, instead of passing this task to Anki
	// probably we can do better caching and error handling, but it require to configure another
	// http client.
	id, err := a.getClient().AddNote(ctx, &noteCopy)
	return NoteID(id), err
}

// SearchProjectedLemmas returns note id for specified lemmas, empty string means not found.
// Note id represented as string, because it's not fully documented what id actually means
func (a *Anki) SearchProjectedLemmas(ctx context.Context, lemmas []*lemma.ProjectedLemma) ([]NoteID, error) {
	client := a.getClient()
	state, err := client.GetState(ctx)
	if err != nil {
		return nil, err
	}
	if !state.IsReadyToAddNote() {
		return nil, ErrIncompleteConfiguration
	}
	// actually checks bellow is redudant because "readiness" should include it
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
func confirmFoundNotes(notes []*ankiconnect.NoteInfo, orderField string, orderValues []string) []NoteID {
	foundIds := make([]NoteID, len(orderValues))
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
			foundIds[i] = NoteID(id)
		}
	}
	return foundIds
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
