package gqlresolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"japwords/graphql/gqlgenerated"
	"japwords/graphql/gqlmodel"
)

func pstring(v string) *string {
	return &v
}

// JapaneseWords is the resolver for the japaneseWords field.
func (r *queryResolver) JapaneseWords(ctx context.Context, query string) (*gqlmodel.JapaneseWords, error) {
	return &gqlmodel.JapaneseWords{
		Words: []*gqlmodel.JapaneseWord{
			{
				Kanji:    "犬",
				Furigana: pstring("犬[いぬ]"),
				Hiragana: "いぬ",
				Acents:   "",
				Meaning:  "dog",
				Audio:    []string{},
				Examples: []string{
					"some example",
				},
			},
			{
				Kanji:    "来年",
				Furigana: pstring("来[らい]年[ねん]"),
				Hiragana: "らいねん",
				Acents:   "",
				Meaning:  "next year",
				Audio:    []string{},
				Examples: []string{},
			},
		},
	}, nil
}

// Query returns gqlgenerated.QueryResolver implementation.
func (r *Resolver) Query() gqlgenerated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
