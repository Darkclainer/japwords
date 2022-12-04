package cachedict

import (
	"context"

	lru "github.com/hashicorp/golang-lru/v2"
)

type Dict[T any] interface {
	Query(ctx context.Context, query string) (T, error)
}

type CacheDict[T any] struct {
	cache *lru.Cache[string, result[T]]
	dict  Dict[T]
}

func New[T any](dict Dict[T]) (*CacheDict[T], error) {
	cache, err := lru.New[string, result[T]](256)
	if err != nil {
		return nil, err
	}
	return &CacheDict[T]{
		cache: cache,
		dict:  dict,
	}, nil
}

func (c *CacheDict[T]) Query(ctx context.Context, query string) (T, error) {
	r, ok := c.cache.Get(query)
	if ok {
		return r.Value, r.Err
	}
	value, err := c.dict.Query(ctx, query)
	c.cache.Add(query, result[T]{
		Value: value,
		Err:   err,
	})
	return value, err
}

type result[T any] struct {
	Value T
	Err   error
}
