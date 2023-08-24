package gqlresolver

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"

	"github.com/Darkclainer/japwords/graphql/gqlmodel"
	"github.com/Darkclainer/japwords/pkg/anki"
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
