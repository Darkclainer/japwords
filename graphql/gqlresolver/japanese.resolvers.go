package gqlresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"japwords/graphql/gqlgenerated"
	"japwords/graphql/gqlmodel"
)

// Lemmas is the resolver for the Lemmas field.
func (r *queryResolver) Lemmas(ctx context.Context, query string) (*gqlmodel.Lemmas, error) {
	panic(fmt.Errorf("not implemented: Lemmas - Lemmas"))
}

// Query returns gqlgenerated.QueryResolver implementation.
func (r *Resolver) Query() gqlgenerated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
