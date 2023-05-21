package ankiconnect

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Anki_FindNotes(t *testing.T) {
	testCases := []struct {
		Name        string
		Handlers    []http.Handler
		Query       string
		Expected    []int64
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "some notes",
			Handlers: []http.Handler{
				handlerAssertRequest(t, &fullRequest{
					Action: "findNotes",
					Params: map[string]any{
						"query": "some query",
					},
				}),
				handlerRespondJSON(t, &fullResponse{
					Result: []int64{1, 2, 3},
				}),
			},
			Query:       "some query",
			Expected:    []int64{1, 2, 3},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "no notes",
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
			ErrorAssert: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, "myspecificerr")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ctx, a := prepareMockServer(t, tc.Handlers...)
			result, err := a.FindNotes(ctx, tc.Query)
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func Test_Anki_AddNote(t *testing.T) {
	testCases := []struct {
		Name            string
		Parameters      *AddNoteParams
		Options         *AddNoteOptions
		Handlers        []http.Handler
		ExpectedRequest map[string]any
		Response        int64
		Expected        int64
		ErrorAssert     assert.ErrorAssertionFunc
	}{
		{
			Name: "options text fields",
			Options: &AddNoteOptions{
				Deck:          "mydeckname",
				Model:         "mymodelname",
				DuplicateDeck: "myduplicatedeck",
			},
			Parameters: &AddNoteParams{},
			ExpectedRequest: map[string]any{
				"deckName":  "mydeckname",
				"modelName": "mymodelname",
				"options": map[string]any{
					"duplicateScope": "deck",
					"duplicateScopeOptions": map[string]any{
						"deckName": "myduplicatedeck",
					},
				},
			},
			Response:    1,
			Expected:    1,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "options duplicate flags check",
			Options: &AddNoteOptions{
				DuplicateFlags: DuplicateFlagsCheck,
			},
			Parameters: &AddNoteParams{},
			ExpectedRequest: map[string]any{
				"options": map[string]any{
					"allowDuplicate":        true,
					"duplicateScope":        "deck",
					"duplicateScopeOptions": map[string]any{},
				},
			},
			Response:    1,
			Expected:    1,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "options duplicate flags with children",
			Options: &AddNoteOptions{
				DuplicateFlags: DuplicateFlagsWithChildren,
			},
			Parameters: &AddNoteParams{},
			ExpectedRequest: map[string]any{
				"options": map[string]any{
					"duplicateScope": "deck",
					"duplicateScopeOptions": map[string]any{
						"checkChildren": true,
					},
				},
			},
			Response:    1,
			Expected:    1,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "options duplicate flags with models",
			Options: &AddNoteOptions{
				DuplicateFlags: DuplicateFlagsWithModels,
			},
			Parameters: &AddNoteParams{},
			ExpectedRequest: map[string]any{
				"options": map[string]any{
					"duplicateScope": "deck",
					"duplicateScopeOptions": map[string]any{
						"checkAllModels": true,
					},
				},
			},
			Response:    1,
			Expected:    1,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "options duplicate flags check with models",
			Options: &AddNoteOptions{
				DuplicateFlags: DuplicateFlagsWithModels | DuplicateFlagsCheck,
			},
			Parameters: &AddNoteParams{},
			ExpectedRequest: map[string]any{
				"options": map[string]any{
					"allowDuplicate": true,
					"duplicateScope": "deck",
					"duplicateScopeOptions": map[string]any{
						"checkAllModels": true,
					},
				},
			},
			Response:    1,
			Expected:    1,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "options duplicate scope deck",
			Options: &AddNoteOptions{
				DuplicateScope: DuplicateScopeDeck,
			},
			Parameters: &AddNoteParams{},
			ExpectedRequest: map[string]any{
				"options": map[string]any{
					"duplicateScope":        "deck",
					"duplicateScopeOptions": map[string]any{},
				},
			},
			Response:    1,
			Expected:    1,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "options duplicate scope all",
			Options: &AddNoteOptions{
				DuplicateScope: DuplicateScopeEverywhere,
			},
			Parameters: &AddNoteParams{},
			ExpectedRequest: map[string]any{
				"options": map[string]any{
					"duplicateScope":        "all",
					"duplicateScopeOptions": map[string]any{},
				},
			},
			Response:    1,
			Expected:    1,
			ErrorAssert: assert.NoError,
		},
		{
			Name:    "params fields",
			Options: &AddNoteOptions{},
			Parameters: &AddNoteParams{
				Fields: map[string]string{
					"f1": "v1",
					"f2": "v2",
				},
			},
			ExpectedRequest: map[string]any{
				"fields": map[string]any{
					"f1": "v1",
					"f2": "v2",
				},
				"options": map[string]any{
					"duplicateScope":        "deck",
					"duplicateScopeOptions": map[string]any{},
				},
			},
			Response:    1,
			Expected:    1,
			ErrorAssert: assert.NoError,
		},
		{
			Name:    "params tags",
			Options: &AddNoteOptions{},
			Parameters: &AddNoteParams{
				Tags: []string{"t1", "t2"},
			},
			ExpectedRequest: map[string]any{
				"tags": []any{"t1", "t2"},
				"options": map[string]any{
					"duplicateScope":        "deck",
					"duplicateScopeOptions": map[string]any{},
				},
			},
			Response:    1,
			Expected:    1,
			ErrorAssert: assert.NoError,
		},
		{
			Name:    "params asset audio",
			Options: &AddNoteOptions{},
			Parameters: &AddNoteParams{
				Assets: []*AddNoteAsset{
					{
						Asset: *NewMediaURL("https://google.com", &MediaAssetOptions{
							DeleteExisting: true,
						}),
						Type: MediaTypeAudio,
					},
				},
			},
			ExpectedRequest: map[string]any{
				"audio": []any{
					map[string]any{
						"url":            "https://google.com",
						"deleteExisting": true,
					},
				},
				"options": map[string]any{
					"duplicateScope":        "deck",
					"duplicateScopeOptions": map[string]any{},
				},
			},
			Response:    1,
			Expected:    1,
			ErrorAssert: assert.NoError,
		},
		{
			Name:    "params asset video",
			Options: &AddNoteOptions{},
			Parameters: &AddNoteParams{
				Assets: []*AddNoteAsset{
					{
						Asset: *NewMediaURL("https://google.com", &MediaAssetOptions{
							DeleteExisting: true,
						}),
						Type: MediaTypeVideo,
					},
				},
			},
			ExpectedRequest: map[string]any{
				"video": []any{
					map[string]any{
						"url":            "https://google.com",
						"deleteExisting": true,
					},
				},
				"options": map[string]any{
					"duplicateScope":        "deck",
					"duplicateScopeOptions": map[string]any{},
				},
			},
			Response:    1,
			Expected:    1,
			ErrorAssert: assert.NoError,
		},
		{
			Name:    "params asset picture",
			Options: &AddNoteOptions{},
			Parameters: &AddNoteParams{
				Assets: []*AddNoteAsset{
					{
						Asset: *NewMediaURL("https://google.com", &MediaAssetOptions{
							DeleteExisting: true,
						}),
						Type: MediaTypePicture,
					},
				},
			},
			ExpectedRequest: map[string]any{
				"picture": []any{
					map[string]any{
						"url":            "https://google.com",
						"deleteExisting": true,
					},
				},
				"options": map[string]any{
					"duplicateScope":        "deck",
					"duplicateScopeOptions": map[string]any{},
				},
			},
			Response:    1,
			Expected:    1,
			ErrorAssert: assert.NoError,
		},
		{
			Name:       "return no note",
			Options:    &AddNoteOptions{},
			Parameters: &AddNoteParams{},
			ExpectedRequest: map[string]any{
				"options": map[string]any{
					"duplicateScope":        "deck",
					"duplicateScopeOptions": map[string]any{},
				},
			},
			Response: 0,
			Expected: 0,
			ErrorAssert: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, "note creation failed")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			handlers := []http.Handler{
				handlerAssertRequest(t, &fullRequest{
					Action: "addNote",
					Params: map[string]any{
						"note": tc.ExpectedRequest,
					},
				}),
				handlerRespondJSON(t, &fullResponse{
					Result: tc.Response,
				}),
			}
			ctx, a := prepareMockServer(t, handlers...)
			id, err := a.AddNote(ctx, tc.Parameters, tc.Options)
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, id)
		})
	}
}

func Test_Anki_DeleteNotes(t *testing.T) {
	testCases := []struct {
		Name        string
		Handlers    []http.Handler
		IDs         []int64
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "some notes",
			Handlers: []http.Handler{
				handlerAssertRequest(t, &fullRequest{
					Action: "deleteNotes",
					Params: map[string]any{
						"notes": []any{float64(2), float64(3), float64(5)},
					},
				}),
				handlerRespondJSON(t, &fullResponse{}),
			},
			IDs:         []int64{2, 3, 5},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "no notes",
			Handlers: []http.Handler{
				handlerRespondJSON(t, &fullResponse{}),
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
			ErrorAssert: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, "myspecificerr")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ctx, a := prepareMockServer(t, tc.Handlers...)
			err := a.DeleteNotes(ctx, tc.IDs)
			tc.ErrorAssert(t, err)
		})
	}
}
