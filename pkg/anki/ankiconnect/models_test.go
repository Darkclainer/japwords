package ankiconnect

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Anki_ModelNames(t *testing.T) {
	testCases := []struct {
		Name        string
		Handlers    []http.Handler
		Expected    []string
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "some models",
			Handlers: []http.Handler{
				handlerAssertRequest(t, &fullRequest{
					Action: "modelNames",
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
			Name: "no models",
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
			result, err := a.ModelNames(ctx)
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func Test_Anki_ModelFieldNames(t *testing.T) {
	testCases := []struct {
		Name        string
		Handlers    []http.Handler
		ModelName   string
		Expected    []string
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "OK",
			Handlers: []http.Handler{
				handlerAssertRequest(t, &fullRequest{
					Action: "modelFieldNames",
					Params: map[string]any{
						"modelName": "mymodelname",
					},
				}),
				handlerRespondJSON(t, &fullResponse{
					Result: []string{"foo", "bar"},
				}),
			},
			ModelName:   "mymodelname",
			Expected:    []string{"foo", "bar"},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "empty",
			Handlers: []http.Handler{
				handlerRespondJSON(t, &fullResponse{
					Result: nil,
				}),
			},
			ModelName:   "mymodelname",
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
			result, err := a.ModelFieldNames(ctx, tc.ModelName)
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func Test_Anki_CreateModel(t *testing.T) {
	testCases := []struct {
		Name        string
		Handlers    []http.Handler
		Request     CreateModelRequest
		Expected    int64
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "OK",
			Handlers: []http.Handler{
				handlerAssertRequest(t, &fullRequest{
					Action: "createModel",
					Params: map[string]any{
						"modelName":     "mymodel",
						"inOrderFields": []any{"foo", "bar"},
						"css":           "mycss",
						"cardTemplates": []any{
							map[string]any{
								"Name":  "footemplate",
								"Front": "foofront",
								"Back":  "fooback",
							},
							map[string]any{
								"Name":  "bartemplate",
								"Front": "barfront",
								"Back":  "barback",
							},
						},
					},
				}),
				handlerRespondJSON(t, &fullResponse{
					Result: map[string]any{
						"id": int64(23323),
					},
				}),
			},
			Request: CreateModelRequest{
				ModelName: "mymodel",
				Fields:    []string{"foo", "bar"},
				CSS:       "mycss",
				CardTemplates: []CreateModelCardTemplate{
					{
						Name:  "footemplate",
						Front: "foofront",
						Back:  "fooback",
					},
					{
						Name:  "bartemplate",
						Front: "barfront",
						Back:  "barback",
					},
				},
			},
			Expected:    23323,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "error",
			Handlers: []http.Handler{
				handlerRespondJSON(t, &fullResponse{
					Error: "myspecificerr",
				}),
			},
			Expected: 0,
			ErrorAssert: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, "myspecificerr")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ctx, a := prepareMockServer(t, tc.Handlers...)
			result, err := a.CreateModel(ctx, &tc.Request)
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}
