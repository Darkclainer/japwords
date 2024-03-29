package jisho

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_queryURL(t *testing.T) {
	testCases := []struct {
		Name        string
		Query       string
		BaseURL     string
		ExpectedURL string
	}{
		{
			Name:        "zero query",
			Query:       "",
			ExpectedURL: "https://jisho.org/search/",
		},
		{
			Name:        "english",
			Query:       "inu",
			ExpectedURL: "https://jisho.org/search/inu",
		},
		{
			Name:        "japanese",
			Query:       "東口",
			ExpectedURL: "https://jisho.org/search/%E6%9D%B1%E5%8F%A3",
		},
		{
			Name:        "with slash",
			Query:       "hel/lo",
			ExpectedURL: "https://jisho.org/search/hel%2Flo",
		},
		{
			Name:        "another basename",
			Query:       "hello",
			BaseURL:     "http://localhost:3890/",
			ExpectedURL: "http://localhost:3890/hello",
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			j := New(nil, tc.BaseURL)
			assert.Equal(t, tc.ExpectedURL, j.queryURL(tc.Query))
		})
	}
}
