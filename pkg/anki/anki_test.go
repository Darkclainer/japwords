package anki

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Darkclainer/japwords/pkg/anki/ankiconnect"
)

func Test_Anki_New(t *testing.T) {
	t.Run("DefaultOK", func(t *testing.T) {
		anki := NewAnki(DefaultClientConstructor)
		assert.Nil(t, anki.wrapper)
	})
}

func TestAnki_ReloadConfig(t *testing.T) {
	t.Run("Options", func(t *testing.T) {
		config := &Config{
			Addr:   "myaddr",
			APIKey: "mykey",
		}
		constructor := func(opts *ankiconnect.Options) (AnkiClient, error) {
			assert.Equal(t, "myaddr", opts.URL)
			assert.Equal(t, "mykey", opts.APIKey)
			return nil, nil
		}
		anki := NewAnki(constructor)
		err := anki.ReloadConfig(config)
		require.NoError(t, err)
		assert.Equal(t, config, anki.wrapper.config)
	})
	t.Run("OK", func(t *testing.T) {
		config := Config{
			Addr:   "myaddr",
			APIKey: "myfirst",
		}
		counter := 0
		constructor := func(opts *ankiconnect.Options) (AnkiClient, error) {
			counter++
			assert.Equal(t, "myaddr", opts.URL)
			if counter == 1 {
				assert.Equal(t, "myfirst", opts.APIKey)
			} else {
				assert.Equal(t, "mysecond", opts.APIKey)
			}
			return nil, nil
		}
		anki := NewAnki(constructor)
		configCopy := config
		err := anki.ReloadConfig(&configCopy)
		require.NoError(t, err)
		config.APIKey = "mysecond"
		err = anki.ReloadConfig(&config)
		assert.NoError(t, err)
		assert.Equal(t, 2, counter)
	})
	t.Run("Error", func(t *testing.T) {
		config := Config{
			Addr:   "myaddr",
			APIKey: "myfirst",
		}
		counter := 0
		constructor := func(opts *ankiconnect.Options) (AnkiClient, error) {
			counter++
			if counter == 1 {
				return DefaultClientConstructor(opts)
			} else {
				return nil, errors.New("testerr")
			}
		}
		anki := NewAnki(constructor)
		err := anki.ReloadConfig(&config)
		require.NoError(t, err)
		wrapper := anki.wrapper
		err = anki.ReloadConfig(&config)
		assert.Error(t, err)
		assert.Same(t, wrapper, anki.wrapper)
	})
}
