package anki

// ATTENTION!
// To run integration tests you need:
//
// * Launch Anki with anki-connect plugin
// * Define variable ANKI_CONNECT_URL and optionally ANKI_CONNECT_API_KEY
//
//	ANKI_CONNECT_URL=http://127.0.0.1:8765 go test
//
// * Create profile "anki-connect-test" (all operation operation that alter data should be contained somehow)

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ankiTestProfile = "anki-connect-test"

func Test_Anki_Version_Integration(t *testing.T) {
	a := aquireAnkiConnect(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	version, err := a.Version(ctx)
	require.NoError(t, err)
	assert.Equal(t, apiVersion, version, "test expected to work with version %d, they may or may not work with another version", apiVersion)
}

func Test_Anki_RequestPermission_Integration(t *testing.T) {
	a := aquireAnkiConnect(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	response, err := a.RequestPermission(ctx)
	require.NoError(t, err)
	assert.Equal(t, "granted", response.Permission)
}

func Test_Anki_LoadProfile_Integration(t *testing.T) {
	a := aquireAnkiConnect(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := a.LoadProfile(ctx, ankiTestProfile)
	require.NoError(t, err)
}

// Test_Anki_Deck_Functions_Integration tests aquiring decks names, creation and deletion of decks in composition
func Test_Anki_Deck_Functions_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	a := aquireAnkiConnectWithProfile(t, ctx)
	initialDecks, err := a.DeckNames(ctx)
	require.NoError(t, err)
	newDeckName := randString("deckfuncs")
	id, err := a.CreateDeck(ctx, newDeckName)
	require.NoError(t, err)
	assert.NotEqual(t, 0, id)
	defer a.DeleteDecks(ctx, []string{newDeckName})
	currentDecks, err := a.DeckNames(ctx)
	require.NoError(t, err)
	initialDecks = append(initialDecks, newDeckName)
	sort.Strings(initialDecks)
	sort.Strings(currentDecks)
	assert.Equal(t, initialDecks, currentDecks)
	err = a.DeleteDecks(ctx, []string{newDeckName})
	assert.NoError(t, err)
}

func randString(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	buffer := make([]byte, 8)
	_, err := rand.Read(buffer)
	if err != nil {
		panic("should not happen")
	}
	return fmt.Sprintf("%s_%s", prefix, base64.StdEncoding.EncodeToString(buffer))
}

func aquireAnkiConnectWithProfile(tb testing.TB, ctx context.Context) *Anki {
	a := aquireAnkiConnect(tb)
	err := a.LoadProfile(ctx, ankiTestProfile)
	require.NoError(tb, err)
	return a
}

func aquireAnkiConnect(tb testing.TB) *Anki {
	url, ok := os.LookupEnv("ANKI_CONNECT_URL")
	if !ok {
		tb.Skipf("test skipped because ANKI_CONNECT_URL is not defined")
	}
	apiKey := os.Getenv("ANKI_CONNECT_API_KEY")
	a, err := New(&Options{
		URL:    url,
		APIKey: apiKey,
	})
	require.NoError(tb, err)
	return a
}
