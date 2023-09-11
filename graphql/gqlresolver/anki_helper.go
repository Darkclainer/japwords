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
	if errors.Is(err, anki.ErrPermissionDenied) {
		return &gqlmodel.AnkiPermissionError{
			Message: err.Error(),
		}, nil
	}
	if errors.Is(err, anki.ErrUnauthorized) {
		return &gqlmodel.AnkiUnauthorizedError{
			Message: err.Error(),
		}, nil
	}
	var connErr *ankiconnect.ConnectionError
	if !errors.As(err, &connErr) {
		return &gqlmodel.AnkiUnauthorizedError{
			Message: err.Error(),
		}, nil
	}
	return &gqlmodel.AnkiConnectionError{
		Message: err.Error(),
	}, nil
}
