package anki

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NoteID_String(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		assert.Equal(t, "", NoteID(0).String())
	})
	t.Run("non zero", func(t *testing.T) {
		assert.Equal(t, "13923", NoteID(13923).String())
	})
}
