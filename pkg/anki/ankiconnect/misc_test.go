package ankiconnect

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Anki_Version(t *testing.T) {
	testCases := []struct {
		Name        string
		Handlers    []http.Handler
		Expected    int
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "OK",
			Handlers: []http.Handler{
				handlerAssertRequest(t, &fullRequest{
					Action: "version",
					Params: nil,
				}),
				handlerRespondJSON(t, &fullResponse{
					Result: 7,
				}),
			},
			Expected:    7,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "error",
			Handlers: []http.Handler{
				handlerRespondJSON(t, &fullResponse{
					Error: "myspecialerror",
				}),
			},
			Expected: 0,
			ErrorAssert: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, "myspecialerror")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ctx, a := prepareMockServer(t, tc.Handlers...)
			result, err := a.Version(ctx)
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func Test_Anki_RequestPermission(t *testing.T) {
	testCases := []struct {
		Name        string
		Handlers    []http.Handler
		Expected    *RequestPermissionResponse
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "granted",
			Handlers: []http.Handler{
				handlerAssertRequest(t, &fullRequest{
					Action: "requestPermission",
					Params: nil,
				}),
				handlerRespondJSON(t, &fullResponse{
					Result: map[string]any{
						"permission":    PermissionGranted,
						"requireApiKey": true,
						"version":       99,
					},
				}),
			},
			Expected: &RequestPermissionResponse{
				Permission:    PermissionGranted,
				RequireAPIKey: true,
				Version:       99,
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "denied",
			Handlers: []http.Handler{
				handlerRespondJSON(t, &fullResponse{
					Result: map[string]any{
						"permission": PermissionDenied,
						"version":    apiVersion,
					},
				}),
			},
			Expected: &RequestPermissionResponse{
				Permission: PermissionDenied,
				Version:    apiVersion,
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "error",
			Handlers: []http.Handler{
				handlerRespondJSON(t, &fullResponse{
					Error: "myspecificerror",
				}),
			},
			ErrorAssert: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, "myspecificerror")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ctx, a := prepareMockServer(t, tc.Handlers...)
			result, err := a.RequestPermission(ctx)
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func Test_Anki_LoadProfile(t *testing.T) {
	testCases := []struct {
		Name        string
		Handlers    []http.Handler
		ProfileName string
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name:        "success",
			ProfileName: "myprofile",
			Handlers: []http.Handler{
				handlerAssertRequest(t, &fullRequest{
					Action: "loadProfile",
					Params: map[string]any{
						"name": "myprofile",
					},
				}),
				handlerRespondJSON(t, &fullResponse{
					Result: true,
				}),
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name:        "failed",
			ProfileName: "myprofile",
			Handlers: []http.Handler{
				handlerRespondJSON(t, &fullResponse{
					Result: false,
				}),
			},
			ErrorAssert: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, "profile load failed")
			},
		},
		{
			Name: "error",
			Handlers: []http.Handler{
				handlerRespondJSON(t, &fullResponse{
					Error: "myspecificerror",
				}),
			},
			ErrorAssert: assert.Error,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			ctx, a := prepareMockServer(t, tc.Handlers...)
			err := a.LoadProfile(ctx, tc.ProfileName)
			tc.ErrorAssert(t, err)
		})
	}
}
