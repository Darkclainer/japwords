package anki

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetDefaultExampleLemmaJSON_nopanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("no panic expected: %s", err)
		}
	}()
	result := GetDefaultExampleLemmaJSON()
	assert.True(t, len(result) > 0)
}
