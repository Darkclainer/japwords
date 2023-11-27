// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlmodel

import (
	"github.com/Darkclainer/japwords/pkg/lemma"
)

type AnkiError interface {
	IsAnkiError()
}

type CreateAnkiDeckError interface {
	IsCreateAnkiDeckError()
}

type CreateDefaultAnkiNoteError interface {
	IsCreateDefaultAnkiNoteError()
}

type Error interface {
	IsError()
	GetMessage() string
}

type Anki struct {
	Decks      *AnkiDecksResult      `json:"decks"`
	Notes      *AnkiNotesResult      `json:"notes"`
	NoteFields *AnkiNoteFieldsResult `json:"noteFields"`
}

type AnkiCollectionUnavailable struct {
	Message string `json:"message"`
	Version int    `json:"version"`
}

func (AnkiCollectionUnavailable) IsError()                {}
func (this AnkiCollectionUnavailable) GetMessage() string { return this.Message }

func (AnkiCollectionUnavailable) IsAnkiError() {}

type AnkiConfig struct {
	Addr     string                `json:"addr"`
	APIKey   string                `json:"apiKey"`
	Deck     string                `json:"deck"`
	NoteType string                `json:"noteType"`
	Mapping  []*AnkiMappingElement `json:"mapping"`
}

type AnkiConfigMappingElementError struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

type AnkiConfigMappingElementInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AnkiConfigMappingError struct {
	FieldErrors []*AnkiConfigMappingElementError `json:"fieldErrors,omitempty"`
	ValueErrors []*AnkiConfigMappingElementError `json:"valueErrors,omitempty"`
	Message     string                           `json:"message"`
}

func (AnkiConfigMappingError) IsError()                {}
func (this AnkiConfigMappingError) GetMessage() string { return this.Message }

type AnkiConfigState struct {
	Version          int  `json:"version"`
	DeckExists       bool `json:"deckExists"`
	NoteTypeExists   bool `json:"noteTypeExists"`
	NoteHasAllFields bool `json:"noteHasAllFields"`
}

type AnkiConfigStateResult struct {
	AnkiConfigState *AnkiConfigState `json:"ankiConfigState,omitempty"`
	Error           AnkiError        `json:"error,omitempty"`
}

type AnkiConnectionError struct {
	Message string `json:"message"`
}

func (AnkiConnectionError) IsError()                {}
func (this AnkiConnectionError) GetMessage() string { return this.Message }

func (AnkiConnectionError) IsAnkiError() {}

type AnkiDecksResult struct {
	Decks []string  `json:"decks,omitempty"`
	Error AnkiError `json:"error,omitempty"`
}

type AnkiForbiddenOrigin struct {
	Message string `json:"message"`
}

func (AnkiForbiddenOrigin) IsError()                {}
func (this AnkiForbiddenOrigin) GetMessage() string { return this.Message }

func (AnkiForbiddenOrigin) IsAnkiError() {}

type AnkiInvalidAPIKey struct {
	Message string `json:"message"`
	Version int    `json:"version"`
}

func (AnkiInvalidAPIKey) IsError()                {}
func (this AnkiInvalidAPIKey) GetMessage() string { return this.Message }

func (AnkiInvalidAPIKey) IsAnkiError() {}

type AnkiMappingElement struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AnkiNoteFieldsResult struct {
	NoteFields []string  `json:"noteFields,omitempty"`
	Error      AnkiError `json:"error,omitempty"`
}

type AnkiNotesResult struct {
	Notes []string  `json:"notes,omitempty"`
	Error AnkiError `json:"error,omitempty"`
}

type AnkiUnknownError struct {
	Message string `json:"message"`
}

func (AnkiUnknownError) IsError()                {}
func (this AnkiUnknownError) GetMessage() string { return this.Message }

func (AnkiUnknownError) IsAnkiError() {}

type Audio struct {
	Type   string `json:"type"`
	Source string `json:"source"`
}

type CreateAnkiDeckAlreadyExists struct {
	Message string `json:"message"`
}

func (CreateAnkiDeckAlreadyExists) IsError()                {}
func (this CreateAnkiDeckAlreadyExists) GetMessage() string { return this.Message }

func (CreateAnkiDeckAlreadyExists) IsCreateAnkiDeckError() {}

type CreateAnkiDeckInput struct {
	Name string `json:"name"`
}

type CreateAnkiDeckResult struct {
	AnkiError AnkiError           `json:"ankiError,omitempty"`
	Error     CreateAnkiDeckError `json:"error,omitempty"`
}

type CreateDefaultAnkiNoteAlreadyExists struct {
	Message string `json:"message"`
}

func (CreateDefaultAnkiNoteAlreadyExists) IsError()                {}
func (this CreateDefaultAnkiNoteAlreadyExists) GetMessage() string { return this.Message }

func (CreateDefaultAnkiNoteAlreadyExists) IsCreateDefaultAnkiNoteError() {}

type CreateDefaultAnkiNoteInput struct {
	Name string `json:"name"`
}

type CreateDefaultAnkiNoteResult struct {
	AnkiError AnkiError                  `json:"ankiError,omitempty"`
	Error     CreateDefaultAnkiNoteError `json:"error,omitempty"`
}

type Lemmas struct {
	Lemmas []*lemma.Lemma `json:"lemmas"`
}

type RenderedField struct {
	Field  string  `json:"field"`
	Result string  `json:"result"`
	Error  *string `json:"error,omitempty"`
}

type RenderedFields struct {
	Template      string           `json:"template"`
	TemplateError *string          `json:"templateError,omitempty"`
	Fields        []*RenderedField `json:"fields"`
}

type SetAnkiConfigConnectionInput struct {
	Addr   string `json:"addr"`
	APIKey string `json:"apiKey"`
}

type SetAnkiConfigConnectionResult struct {
	Error *ValidationError `json:"error,omitempty"`
}

type SetAnkiConfigDeckInput struct {
	Name string `json:"name"`
}

type SetAnkiConfigDeckResult struct {
	Error *ValidationError `json:"error,omitempty"`
}

type SetAnkiConfigMappingInput struct {
	Mapping []*AnkiConfigMappingElementInput `json:"mapping"`
}

type SetAnkiConfigMappingResult struct {
	Error *AnkiConfigMappingError `json:"error,omitempty"`
}

type SetAnkiConfigNote struct {
	Name string `json:"name"`
}

type SetAnkiConfigNoteResult struct {
	Error *ValidationError `json:"error,omitempty"`
}

type ValidationError struct {
	Paths   []string `json:"paths"`
	Message string   `json:"message"`
}

func (ValidationError) IsCreateAnkiDeckError() {}

func (ValidationError) IsCreateDefaultAnkiNoteError() {}

func (ValidationError) IsError()                {}
func (this ValidationError) GetMessage() string { return this.Message }
