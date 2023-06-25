package gqlresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"

	"github.com/Darkclainer/japwords/graphql/gqlgenerated"
	"github.com/Darkclainer/japwords/graphql/gqlmodel"
)

// Lemmas is the resolver for the Lemmas field.
func (r *queryResolver) Lemmas(ctx context.Context, query string) (*gqlmodel.Lemmas, error) {
	lemmas, err := r.multiDict.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var gqlLemmas []*gqlmodel.Lemma
	for _, lemma := range lemmas {
		gqlLemmas = append(gqlLemmas, convertLemma(lemma))
	}
	return &gqlmodel.Lemmas{
		Lemmas: gqlLemmas,
	}, nil
}

// Query returns gqlgenerated.QueryResolver implementation.
func (r *Resolver) Query() gqlgenerated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
