package gqlresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"

	"github.com/Darkclainer/japwords/graphql/gqlgenerated"
	"github.com/Darkclainer/japwords/pkg/lemma"
)

// Furigana is the resolver for the furigana field.
func (r *wordResolver) Furigana(ctx context.Context, obj *lemma.Word) ([]*lemma.FuriganaChar, error) {
	return sliceToPointers(obj.Furigana), nil
}

// Furigana is the resolver for the furigana field.
func (r *wordInputResolver) Furigana(ctx context.Context, obj *lemma.Word, data []*lemma.FuriganaChar) error {
	obj.Furigana = sliceToValues(data)
	return nil
}

// Word returns gqlgenerated.WordResolver implementation.
func (r *Resolver) Word() gqlgenerated.WordResolver { return &wordResolver{r} }

// WordInput returns gqlgenerated.WordInputResolver implementation.
func (r *Resolver) WordInput() gqlgenerated.WordInputResolver { return &wordInputResolver{r} }

type wordResolver struct{ *Resolver }
type wordInputResolver struct{ *Resolver }
