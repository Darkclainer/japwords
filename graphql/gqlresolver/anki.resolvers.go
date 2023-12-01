package gqlresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"slices"

	"github.com/Darkclainer/japwords/graphql/gqlgenerated"
	"github.com/Darkclainer/japwords/graphql/gqlmodel"
	"github.com/Darkclainer/japwords/pkg/anki"
	"github.com/Darkclainer/japwords/pkg/lemma"
)

// Decks is the resolver for the decks field.
func (r *ankiResolver) Decks(ctx context.Context, obj *gqlmodel.Anki) (*gqlmodel.AnkiDecksResult, error) {
	decks, err := r.ankiClient.Decks(ctx)
	if err != nil {
		if ankiErr, _ := convertAnkiError(err); ankiErr != nil {
			return &gqlmodel.AnkiDecksResult{
				Error: ankiErr,
			}, nil
		}
		return nil, err
	}
	return &gqlmodel.AnkiDecksResult{
		Decks: decks,
	}, nil
}

// Notes is the resolver for the notes field.
func (r *ankiResolver) Notes(ctx context.Context, obj *gqlmodel.Anki) (*gqlmodel.AnkiNotesResult, error) {
	noteTypes, err := r.ankiClient.NoteTypes(ctx)
	if err != nil {
		if ankiErr, _ := convertAnkiError(err); ankiErr != nil {
			return &gqlmodel.AnkiNotesResult{
				Error: ankiErr,
			}, nil
		}
		return nil, err
	}
	return &gqlmodel.AnkiNotesResult{
		Notes: noteTypes,
	}, nil
}

// NoteFields is the resolver for the noteFields field.
func (r *ankiResolver) NoteFields(ctx context.Context, obj *gqlmodel.Anki, name string) (*gqlmodel.AnkiNoteFieldsResult, error) {
	fields, err := r.ankiClient.NoteTypeFields(ctx)
	if err != nil {
		if ankiErr, _ := convertAnkiError(err); ankiErr != nil {
			return &gqlmodel.AnkiNoteFieldsResult{
				Error: ankiErr,
			}, nil
		}
		return nil, err
	}
	return &gqlmodel.AnkiNoteFieldsResult{
		NoteFields: fields,
	}, nil
}

// SetAnkiConfigConnection is the resolver for the setAnkiConfigConnection field.
func (r *mutationResolver) SetAnkiConfigConnection(ctx context.Context, input gqlmodel.SetAnkiConfigConnectionInput) (*gqlmodel.SetAnkiConfigConnectionResult, error) {
	err := r.ankiConfig.UpdateConnection(input.Addr, input.APIKey)
	if validationErr, _ := convertAnkiValidationError(ctx, err); validationErr != nil {
		return &gqlmodel.SetAnkiConfigConnectionResult{
			Error: validationErr,
		}, nil
	}
	return &gqlmodel.SetAnkiConfigConnectionResult{}, err
}

// SetAnkiConfigDeck is the resolver for the setAnkiConfigDeck field.
func (r *mutationResolver) SetAnkiConfigDeck(ctx context.Context, input gqlmodel.SetAnkiConfigDeckInput) (*gqlmodel.SetAnkiConfigDeckResult, error) {
	err := r.ankiConfig.UpdateDeck(input.Name)
	if validationErr, _ := convertAnkiValidationError(ctx, err); validationErr != nil {
		return &gqlmodel.SetAnkiConfigDeckResult{
			Error: validationErr,
		}, nil
	}
	return &gqlmodel.SetAnkiConfigDeckResult{}, err
}

// SetAnkiConfigNote is the resolver for the setAnkiConfigNote field.
func (r *mutationResolver) SetAnkiConfigNote(ctx context.Context, input gqlmodel.SetAnkiConfigNote) (*gqlmodel.SetAnkiConfigNoteResult, error) {
	err := r.ankiConfig.UpdateNoteType(input.Name)
	if validationErr, _ := convertAnkiValidationError(ctx, err); validationErr != nil {
		return &gqlmodel.SetAnkiConfigNoteResult{
			Error: validationErr,
		}, nil
	}
	return &gqlmodel.SetAnkiConfigNoteResult{}, err
}

// SetAnkiConfigMapping is the resolver for the setAnkiConfigMapping field.
func (r *mutationResolver) SetAnkiConfigMapping(ctx context.Context, input gqlmodel.SetAnkiConfigMappingInput) (*gqlmodel.SetAnkiConfigMappingResult, error) {
	mapping := map[string]string{}
	for _, element := range input.Mapping {
		mapping[element.Key] = element.Value
	}
	err := r.ankiConfig.UpdateMapping(mapping)
	if err != nil {
		var ankiMappingErrs *anki.MappingValidationErrors
		if !errors.As(err, &ankiMappingErrs) {
			return nil, err
		}
		mappingErrs := &gqlmodel.AnkiConfigMappingError{
			Message: "invalid mapping",
		}
		for _, err := range ankiMappingErrs.KeyErrors {
			mappingErrs.FieldErrors = append(mappingErrs.FieldErrors,
				&gqlmodel.AnkiConfigMappingElementError{
					Key:     err.Key,
					Message: err.Msg,
				},
			)
		}
		for _, err := range ankiMappingErrs.ValueErrors {
			mappingErrs.ValueErrors = append(mappingErrs.ValueErrors,
				&gqlmodel.AnkiConfigMappingElementError{
					Key:     err.Key,
					Message: err.Msg,
				},
			)
		}
		return &gqlmodel.SetAnkiConfigMappingResult{Error: mappingErrs}, nil
	}
	return &gqlmodel.SetAnkiConfigMappingResult{}, nil
}

// CreateAnkiDeck is the resolver for the createAnkiDeck field.
func (r *mutationResolver) CreateAnkiDeck(ctx context.Context, input *gqlmodel.CreateAnkiDeckInput) (*gqlmodel.CreateAnkiDeckResult, error) {
	err := r.ankiClient.CreateDeck(ctx, input.Name)
	if err != nil {
		if errors.Is(err, anki.ErrDeckAlreadyExists) {
			return &gqlmodel.CreateAnkiDeckResult{
				Error: &gqlmodel.CreateAnkiDeckAlreadyExists{
					Message: err.Error(),
				},
			}, nil
		}
		if validationErr, _ := convertAnkiValidationError(ctx, err); validationErr != nil {
			return &gqlmodel.CreateAnkiDeckResult{
				Error: validationErr,
			}, nil
		}
		if ankiErr, _ := convertAnkiError(err); ankiErr != nil {
			return &gqlmodel.CreateAnkiDeckResult{
				AnkiError: ankiErr,
			}, nil
		}
		return nil, err
	}
	return &gqlmodel.CreateAnkiDeckResult{}, nil
}

// CreateAnkiNote is the resolver for the createAnkiNote field.
func (r *mutationResolver) CreateDefaultAnkiNote(ctx context.Context, input *gqlmodel.CreateDefaultAnkiNoteInput) (*gqlmodel.CreateDefaultAnkiNoteResult, error) {
	err := r.ankiClient.CreateDefaultNote(ctx, input.Name)
	if err != nil {
		if validationErr, _ := convertAnkiValidationError(ctx, err); validationErr != nil {
			return &gqlmodel.CreateDefaultAnkiNoteResult{
				Error: validationErr,
			}, nil
		}
		if ankiErr, _ := convertAnkiError(err); ankiErr != nil {
			return &gqlmodel.CreateDefaultAnkiNoteResult{
				AnkiError: ankiErr,
			}, nil
		}
		return nil, err
	}
	return &gqlmodel.CreateDefaultAnkiNoteResult{}, nil
}

// AddAnkiNote is the resolver for the addAnkiNote field.
func (r *mutationResolver) AddAnkiNote(ctx context.Context, request *anki.AddNoteRequest) (*gqlmodel.AnkiAddNoteResult, error) {
	err := r.ankiClient.AddNote(ctx, request)
	if err != nil {
		if errors.Is(err, anki.ErrDuplicatedNoteFound) {
			return &gqlmodel.AnkiAddNoteResult{
				Error: &gqlmodel.AnkiAddNoteDuplicateFound{
					Message: err.Error(),
				},
			}, nil
		}
		if errors.Is(err, anki.ErrIncompleteConfiguration) {
			return &gqlmodel.AnkiAddNoteResult{
				Error: &gqlmodel.AnkiIncompleteConfiguration{
					Message: err.Error(),
				},
			}, nil
		}
		if ankiErr, _ := convertAnkiError(err); ankiErr != nil {
			return &gqlmodel.AnkiAddNoteResult{
				AnkiError: ankiErr,
			}, nil
		}
		return nil, err

	}
	return &gqlmodel.AnkiAddNoteResult{}, nil
}

// Anki is the resolver for the Anki field.
func (r *queryResolver) Anki(ctx context.Context) (*gqlmodel.Anki, error) {
	return &gqlmodel.Anki{}, nil
}

// AnkiConfigState is the resolver for the AnkiConfigState field.
func (r *queryResolver) AnkiConfigState(ctx context.Context) (*gqlmodel.AnkiConfigStateResult, error) {
	state, err := r.ankiClient.FullStateCheck(ctx)
	if err != nil {
		if ankiErr, _ := convertAnkiError(err); ankiErr != nil {
			return &gqlmodel.AnkiConfigStateResult{
				Error: ankiErr,
			}, nil
		}
		return nil, err
	}
	result := &gqlmodel.AnkiConfigState{
		Version:          state.Version,
		DeckExists:       state.DeckExists,
		NoteTypeExists:   state.NoteTypeExists,
		NoteHasAllFields: state.NoteHasAllFields,
		OrderDefined:     state.OrderDefined,
	}
	return &gqlmodel.AnkiConfigStateResult{
		AnkiConfigState: result,
	}, nil
}

// AnkiConfig is the resolver for the AnkiConfig field.
func (r *queryResolver) AnkiConfig(ctx context.Context) (*gqlmodel.AnkiConfig, error) {
	ankiConfig := r.configManager.Current().Anki
	result := &gqlmodel.AnkiConfig{
		Addr:     ankiConfig.Addr,
		APIKey:   ankiConfig.APIKey,
		Deck:     ankiConfig.Deck,
		NoteType: ankiConfig.NoteType,
		Mapping:  nil,
	}
	mapping := make([]*gqlmodel.AnkiMappingElement, 0, len(ankiConfig.FieldMapping))
	for key, value := range ankiConfig.FieldMapping {
		mapping = append(mapping, &gqlmodel.AnkiMappingElement{
			Key:   key,
			Value: value,
		})
	}
	slices.SortStableFunc(mapping, func(a, b *gqlmodel.AnkiMappingElement) int {
		return cmp.Compare(a.Key, b.Key)
	})
	result.Mapping = mapping
	return result, nil
}

// RenderFields is the resolver for the RenderFields field.
func (r *queryResolver) RenderFields(ctx context.Context, fields []string, template *string) (*gqlmodel.RenderedFields, error) {
	var (
		currentLemma = anki.DefaultExampleLemma
		lemmaSrc     = anki.GetDefaultExampleLemmaJSON()
	)
	if template != nil {
		var newLemma lemma.ProjectedLemma
		err := json.Unmarshal([]byte(*template), &newLemma)
		if err != nil {
			errString := err.Error()
			return &gqlmodel.RenderedFields{
				Template:      lemmaSrc,
				TemplateError: &errString,
			}, nil
		}
		currentLemma = newLemma
		src, err := json.MarshalIndent(&currentLemma, "", "  ")
		if err != nil {
			errString := err.Error()
			return &gqlmodel.RenderedFields{
				Template:      lemmaSrc,
				TemplateError: &errString,
			}, nil
		}
		lemmaSrc = string(src)
	}
	renderedFields := make([]*gqlmodel.RenderedField, len(fields))
	for i, field := range fields {
		renderedField := &gqlmodel.RenderedField{
			Field: field,
		}
		result, err := anki.RenderRawTemplate(field, &currentLemma)
		if err != nil {
			errString := err.Error()
			renderedField.Error = &errString
		} else {
			renderedField.Result = result
		}
		renderedFields[i] = renderedField
	}
	return &gqlmodel.RenderedFields{
		Template: lemmaSrc,
		Fields:   renderedFields,
	}, nil
}

// PrepareProjectedLemma is the resolver for the PrepareProjectedLemma field.
func (r *queryResolver) PrepareProjectedLemma(ctx context.Context, lemma *lemma.ProjectedLemma) (*gqlmodel.PrepareProjectedLemmaResult, error) {
	request, err := r.ankiClient.PrepareProjectedLemma(ctx, lemma)
	if err != nil {
		if errors.Is(err, anki.ErrIncompleteConfiguration) {
			return &gqlmodel.PrepareProjectedLemmaResult{
				Error: &gqlmodel.AnkiIncompleteConfiguration{
					Message: err.Error(),
				},
			}, nil
		}
		if ankiErr, _ := convertAnkiError(err); ankiErr != nil {
			return &gqlmodel.PrepareProjectedLemmaResult{
				AnkiError: ankiErr,
			}, nil
		}

		return nil, err
	}
	return &gqlmodel.PrepareProjectedLemmaResult{
		Request: request,
	}, nil
}

// Anki returns gqlgenerated.AnkiResolver implementation.
func (r *Resolver) Anki() gqlgenerated.AnkiResolver { return &ankiResolver{r} }

// Mutation returns gqlgenerated.MutationResolver implementation.
func (r *Resolver) Mutation() gqlgenerated.MutationResolver { return &mutationResolver{r} }

// Query returns gqlgenerated.QueryResolver implementation.
func (r *Resolver) Query() gqlgenerated.QueryResolver { return &queryResolver{r} }

type (
	ankiResolver     struct{ *Resolver }
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
)
