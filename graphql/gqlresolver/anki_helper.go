package gqlresolver

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"

	"github.com/Darkclainer/japwords/graphql/gqlmodel"
	"github.com/Darkclainer/japwords/pkg/anki"
	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
)

func convertAnkiValidationError(ctx context.Context, err error) (*gqlmodel.ValidationError, error) {
	var validationErr *anki.ValidationError
	if !errors.As(err, &validationErr) {
		return nil, err
	}
	result := gqlmodel.ValidationError{
		Paths:   []string{graphql.GetPath(ctx).String()},
		Message: validationErr.Msg,
	}
	return &result, nil
}

func convertAnkiError(err error) (gqlmodel.AnkiError, error) {
	switch {
	case errors.Is(err, anki.ErrForbiddenOrigin):
		return &gqlmodel.AnkiForbiddenOrigin{
			Message: err.Error(),
		}, nil
	case errors.Is(err, anki.ErrCollectionUnavailable):
		return &gqlmodel.AnkiCollectionUnavailable{
			Message: err.Error(),
		}, nil
	case errors.Is(err, anki.ErrInvalidAPIKey):
		return &gqlmodel.AnkiInvalidAPIKey{
			Message: err.Error(),
		}, nil
	}
	var connErr *ankiconnect.ConnectionError
	if errors.As(err, &connErr) {
		return &gqlmodel.AnkiConnectionError{
			Message: err.Error(),
		}, nil
	}
	return &gqlmodel.AnkiUnknownError{
		Message: err.Error(),
	}, nil
}
