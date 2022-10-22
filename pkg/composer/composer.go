package composer

import (
	"context"

	"japwords/pkg/jisho"
)

// Workers?
// Compose word from two dictionary?
// ?

type Composer struct {
	jisho Jisho
}

type Jisho interface {
	Query(context.Context, string) ([]*jisho.Lemma, error)
}
