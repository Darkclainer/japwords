package anki

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MediaAssetOptions_applyToRequest(t *testing.T) {
	testCases := []struct {
		Name     string
		Options  *MediaAssetOptions
		Expected *MediaAssetRequest
	}{
		{
			Name:     "empty",
			Options:  &MediaAssetOptions{},
			Expected: &MediaAssetRequest{},
		},
		{
			Name: "filename",
			Options: &MediaAssetOptions{
				Filename: "filename",
			},
			Expected: &MediaAssetRequest{
				Filename: "filename",
			},
		},
		{
			Name: "delete existing",
			Options: &MediaAssetOptions{
				DeleteExisting: true,
			},
			Expected: &MediaAssetRequest{
				DeleteExisting: true,
			},
		},
		{
			Name: "fields",
			Options: &MediaAssetOptions{
				Fields: []string{"hello", "world"},
			},
			Expected: &MediaAssetRequest{
				Fields: []string{"hello", "world"},
			},
		},
		{
			Name: "skip hash",
			Options: &MediaAssetOptions{
				SkipHash: []byte{0x23, 0x47},
			},
			Expected: &MediaAssetRequest{
				SkipHash: "2347",
			},
		},
		{
			Name: "everything",
			Options: &MediaAssetOptions{
				Filename:       "filename",
				DeleteExisting: true,
				Fields:         []string{"hello", "world"},
				SkipHash:       []byte{0x23, 0x47},
			},
			Expected: &MediaAssetRequest{
				Filename:       "filename",
				Fields:         []string{"hello", "world"},
				SkipHash:       "2347",
				DeleteExisting: true,
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			request := &MediaAssetRequest{}
			tc.Options.applyToRequest(request)
			assert.Equal(t, tc.Expected, request)
		})
	}
}

func Test_NewMediaFile(t *testing.T) {
	testCases := []struct {
		Name     string
		Filename string
		Options  *MediaAssetOptions
		Expected *MediaAssetRequest
	}{
		{
			Name:     "empty",
			Options:  &MediaAssetOptions{},
			Expected: &MediaAssetRequest{},
		},
		{
			Name:     "only filename",
			Filename: "myfilename",
			Options:  &MediaAssetOptions{},
			Expected: &MediaAssetRequest{
				Path: "myfilename",
			},
		},
		{
			Name:     "filename with options",
			Filename: "myfilename",
			Options: &MediaAssetOptions{
				Filename: "optiosfilename",
			},
			Expected: &MediaAssetRequest{
				Path:     "myfilename",
				Filename: "optiosfilename",
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			request := NewMediaFile(tc.Filename, tc.Options)
			assert.Equal(t, tc.Expected, request)
		})
	}
}

func Test_NewMediaURL(t *testing.T) {
	testCases := []struct {
		Name     string
		URL      string
		Options  *MediaAssetOptions
		Expected *MediaAssetRequest
	}{
		{
			Name:     "empty",
			Options:  &MediaAssetOptions{},
			Expected: &MediaAssetRequest{},
		},
		{
			Name:    "url",
			URL:     "some url",
			Options: &MediaAssetOptions{},
			Expected: &MediaAssetRequest{
				URL: "some url",
			},
		},
		{
			Name: "url with options",
			URL:  "some url",
			Options: &MediaAssetOptions{
				Fields: []string{"myfield"},
			},
			Expected: &MediaAssetRequest{
				URL:    "some url",
				Fields: []string{"myfield"},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			request := NewMediaURL(tc.URL, tc.Options)
			assert.Equal(t, tc.Expected, request)
		})
	}
}

func Test_NewMediaBlob(t *testing.T) {
	testCases := []struct {
		Name     string
		Data     []byte
		Options  *MediaAssetOptions
		Expected *MediaAssetRequest
	}{
		{
			Name:     "empty",
			Options:  &MediaAssetOptions{},
			Expected: &MediaAssetRequest{},
		},
		{
			Name:    "with data",
			Data:    []byte("hello world"),
			Options: &MediaAssetOptions{},
			Expected: &MediaAssetRequest{
				Data: "aGVsbG8gd29ybGQ=", // used base64 gnu tool
			},
		},
		{
			Name: "with data",
			Data: []byte("hello world"),
			Options: &MediaAssetOptions{
				Fields: []string{"myfield"},
			},
			Expected: &MediaAssetRequest{
				Data:   "aGVsbG8gd29ybGQ=", // used base64 gnu tool
				Fields: []string{"myfield"},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			request := NewMediaBlob(tc.Data, tc.Options)
			assert.Equal(t, tc.Expected, request)
		})
	}
}
