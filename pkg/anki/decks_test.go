package anki

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Anki_DeckNames(t *testing.T) {
	testCases := []struct {
		Name        string
		Handlers    []http.Handler
		Expected    []string
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "some decks",
			Handlers: []http.Handler{
				handlerAssertRequest(t, &fullRequest{
					Action: "deckNames",
					Params: nil,
				}),
				handlerRespondJSON(t, &fullResponse{
					Result: []string{"foo", "bar"},
				}),
			},
			Expected:    []string{"foo", "bar"},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "no decks",
			Handlers: []http.Handler{
				handlerRespondJSON(t, &fullResponse{
					Result: nil,
				}),
			},
			Expected:    nil,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "error",
			Handlers: []http.Handler{
				handlerRespondJSON(t, &fullResponse{
					Error: "myspecificerr",
				}),
			},
			Expected: nil,
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "myspecificerr")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ctx, a := prepareMockServer(t, tc.Handlers...)
			result, err := a.DeckNames(ctx)
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func Test_Anki_CreateDeck(t *testing.T) {
	testCases := []struct {
		Name        string
		Handlers    []http.Handler
		DeckName    string
		Expected    int
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name:     "ok",
			DeckName: "myname",
			Handlers: []http.Handler{
				handlerAssertRequest(t, &fullRequest{
					Action: "createDeck",
					Params: map[string]interface{}{
						"deck": "myname",
					},
				}),
				handlerRespondJSON(t, &fullResponse{
					Result: 123321,
				}),
			},
			Expected:    123321,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "error",
			Handlers: []http.Handler{
				handlerRespondJSON(t, &fullResponse{
					Error: "myspecificerr",
				}),
			},
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "myspecificerr")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ctx, a := prepareMockServer(t, tc.Handlers...)
			result, err := a.CreateDeck(ctx, tc.DeckName)
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func Test_Anki_DeleteDecks(t *testing.T) {
	testCases := []struct {
		Name        string
		Handlers    []http.Handler
		DeckNames   []string
		CardsToo    bool
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name:      "ok",
			DeckNames: []string{"foo", "bar"},
			CardsToo:  true,
			Handlers: []http.Handler{
				handlerAssertRequest(t, &fullRequest{
					Action: "deleteDecks",
					Params: map[string]interface{}{
						"decks":    []interface{}{"foo", "bar"},
						"cardsToo": true,
					},
				}),
				handlerRespondJSON(t, &fullResponse{
					Result: nil,
				}),
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "error",
			Handlers: []http.Handler{
				handlerRespondJSON(t, &fullResponse{
					Error: "myspecificerr",
				}),
			},
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "myspecificerr")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ctx, a := prepareMockServer(t, tc.Handlers...)
			err := a.DeleteDecks(ctx, tc.DeckNames)
			tc.ErrorAssert(t, err)
		})
	}
}
