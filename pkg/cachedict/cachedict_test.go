package cachedict

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CacheDict_Query(t *testing.T) {
	var called int
	dict := TestDictHandler(func(_ context.Context, _ string) (string, error) {
		called++
		return "my result", errors.New("my error")
	})
	cacheDict, err := New[string](dict)
	require.NoError(t, err)
	ctx := context.Background()

	result, err := cacheDict.Query(ctx, "query1")
	assert.Equal(t, "my result", result)
	assert.ErrorContains(t, err, "my error")
	assert.Equal(t, 1, called)

	result, err = cacheDict.Query(ctx, "query1")
	assert.Equal(t, "my result", result)
	assert.ErrorContains(t, err, "my error")
	assert.Equal(t, 1, called)

	result, err = cacheDict.Query(ctx, "query2")
	assert.Equal(t, "my result", result)
	assert.ErrorContains(t, err, "my error")
	assert.Equal(t, 2, called)
}

type TestDictHandler func(context.Context, string) (string, error)

func (td TestDictHandler) Query(ctx context.Context, query string) (string, error) {
	return td(ctx, query)
}
