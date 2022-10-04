package gqlresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"japwords/graphql/gqlgenerated"
	"japwords/graphql/gqlmodel"
)

// JapaneseWords is the resolver for the japaneseWords field.
func (r *queryResolver) JapaneseWords(ctx context.Context, query string) (*gqlmodel.JapaneseWords, error) {
	panic(fmt.Errorf("not implemented: JapaneseWords - japaneseWords"))
}

// Query returns gqlgenerated.QueryResolver implementation.
func (r *Resolver) Query() gqlgenerated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
