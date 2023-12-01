package ankiconnect

// ATTENTION!
//
// To run integration tests you need:
//
// * Launch Anki with anki-connect plugin
// * Define variable ANKI_CONNECT_URL and optionally ANKI_CONNECT_API_KEY
//
//	ANKI_CONNECT_URL=http://127.0.0.1:8765 go test
//
// * Create profile "anki-connect-test" (all operation operation that alter data should be contained somehow)
//
// Also note that model can not be deleted with anki-connect API for current moment, so
// profiles will be anyway polluted with test models (note types).

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
	assert.Equal(t, PermissionGranted, response.Permission)
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
	defer func() { assert.NoError(t, a.DeleteDecks(ctx, []string{newDeckName})) }()
	currentDecks, err := a.DeckNames(ctx)
	require.NoError(t, err)
	initialDecks = append(initialDecks, newDeckName)
	sort.Strings(initialDecks)
	sort.Strings(currentDecks)
	assert.Equal(t, initialDecks, currentDecks)
	err = a.DeleteDecks(ctx, []string{newDeckName})
	assert.NoError(t, err)
}

// Test_Anki_Model_Functions_Integration tests creation, and search of models
func Test_Anki_Model_Functions_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	a := aquireAnkiConnectWithProfile(t, ctx)
	initialModels, err := a.ModelNames(ctx)
	require.NoError(t, err)
	newModelName := randString("modelfuncs")
	id, err := a.CreateModel(ctx, &CreateModelRequest{
		ModelName: newModelName,
		Fields:    []string{"foo", "bar"},
		CSS:       "",
		CardTemplates: []CreateModelCardTemplate{
			{
				Name:  "mynote",
				Front: "{{foo}}",
				Back:  "{{bar}}",
			},
		},
	})
	require.NoError(t, err)
	assert.NotEqual(t, 0, id)
	currentModels, err := a.ModelNames(ctx)
	require.NoError(t, err)
	initialModels = append(initialModels, newModelName)
	sort.Strings(initialModels)
	sort.Strings(currentModels)
	assert.Equal(t, initialModels, currentModels)
	modelFields, err := a.ModelFieldNames(ctx, newModelName)
	require.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, modelFields)
}

// Test_Anki_Note_Functions_Integrations tests creation, deletion and search of notes
func Test_Anki_Note_Functions_Integrations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	a := aquireAnkiConnectWithProfile(t, ctx)
	// create model for tests
	modelName := randString("notesfuncs")
	_, err := a.CreateModel(ctx, &CreateModelRequest{
		ModelName: modelName,
		Fields:    []string{"foo", "bar"},
		CSS:       "",
		CardTemplates: []CreateModelCardTemplate{
			{
				Name:  "mynote",
				Front: "{{foo}}",
				Back:  "{{bar}}",
			},
		},
	})
	require.NoError(t, err)
	deckName := randString("deckfuncs")
	_, err = a.CreateDeck(ctx, deckName)
	require.NoError(t, err)
	defer func() { assert.NoError(t, a.DeleteDecks(ctx, []string{deckName})) }()
	noteID, err := a.AddNote(
		ctx,
		&AddNoteParams{
			Fields: map[string]string{
				"foo": "hello",
				"bar": "world",
			},
		},
		&AddNoteOptions{
			Deck:  deckName,
			Model: modelName,
		},
	)
	require.NoError(t, err)
	searchedIds, err := a.FindNotes(ctx, fmt.Sprintf("nid:%d", noteID))
	require.NoError(t, err)
	assert.Equal(t, []int64{noteID}, searchedIds)
	requstedNotes, err := a.NotesInfo(ctx, []int64{noteID})
	require.NoError(t, err)
	assert.Equal(t,
		[]*NoteInfo{
			{
				NoteID:    noteID,
				ModelName: modelName,
				Tags:      []string{},
				Fields: map[string]*NoteInfoField{
					"foo": {
						Value: "hello",
						Order: 0,
					},
					"bar": {
						Value: "world",
						Order: 1,
					},
				},
			},
		},
		requstedNotes)
	err = a.DeleteNotes(ctx, []int64{noteID})
	require.NoError(t, err)
	searchedIds, err = a.FindNotes(ctx, fmt.Sprintf("nid:%d", noteID))
	require.NoError(t, err)
	assert.Len(t, searchedIds, 0)
}

func randString(prefix string) string {
	source := rand.New(rand.NewSource(time.Now().Unix()))
	buffer := make([]byte, 8)
	_, err := source.Read(buffer)
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
