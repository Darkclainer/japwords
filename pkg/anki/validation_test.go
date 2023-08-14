package anki

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateMappingKeys(t *testing.T) {
	testCases := []struct {
		Name        string
		Mapping     map[string]string
		ErrorAssert assert.ValueAssertionFunc
	}{
		{
			Name:        "empty",
			Mapping:     map[string]string{},
			ErrorAssert: assert.Empty,
		},
		{
			Name: "ok",
			Mapping: map[string]string{
				"hello": "",
			},
			ErrorAssert: assert.Empty,
		},
		{
			Name: "single error",
			Mapping: map[string]string{
				"hel:lo": "",
			},
			ErrorAssert: func(tt assert.TestingT, val interface{}, _ ...interface{}) bool {
				return assert.Len(t, val, 1)
			},
		},
		{
			Name: "two errors",
			Mapping: map[string]string{
				"hel:lo": "",
				"hello":  "",
				"hel{lo": "",
			},
			ErrorAssert: func(tt assert.TestingT, val interface{}, _ ...interface{}) bool {
				return assert.Len(t, val, 2)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			errs := validateMappingKeys(tc.Mapping)
			tc.ErrorAssert(t, errs)
		})
	}
}

func Test_validateAddr(t *testing.T) {
	testCases := []struct {
		Addr        string
		ErrorAssert assert.ErrorAssertionFunc
	}{
		// positive
		{
			Addr:        "127.0.0.1:8765",
			ErrorAssert: assert.NoError,
		},
		{
			Addr:        "127.0.0.1:1",
			ErrorAssert: assert.NoError,
		},
		{
			Addr:        "127.0.0.1:65535",
			ErrorAssert: assert.NoError,
		},
		{
			Addr:        "example.com:8765",
			ErrorAssert: assert.NoError,
		},
		{
			Addr:        "localhost:8765",
			ErrorAssert: assert.NoError,
		},
		{
			Addr:        "localhost:000008765",
			ErrorAssert: assert.NoError,
		},
		// negative
		{
			Addr:        "",
			ErrorAssert: assert.Error,
		},
		{
			Addr:        "127.0.0.1",
			ErrorAssert: assert.Error,
		},
		{
			Addr:        ":8765",
			ErrorAssert: assert.Error,
		},
		{
			Addr:        "127.0.0.1:0",
			ErrorAssert: assert.Error,
		},
		{
			Addr:        "127.0.0.1:65536",
			ErrorAssert: assert.Error,
		},
		{
			Addr:        "127.0.0.1:-80",
			ErrorAssert: assert.Error,
		},
		{
			Addr:        "127.0.0.1::80",
			ErrorAssert: assert.Error,
		},
		{
			Addr:        "example..com::80",
			ErrorAssert: assert.Error,
		},
		{
			Addr:        "examplecom.::80",
			ErrorAssert: assert.Error,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Addr, func(t *testing.T) {
			err := validateAddr(tc.Addr)
			tc.ErrorAssert(t, err, "%s", tc.Addr)
		})
	}
}

func Test_validateDeckName(t *testing.T) {
	testCases := getDeckNoteTypeTestCases()
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Src, func(t *testing.T) {
			err := validateDeckName(tc.Src)
			tc.ErrorAssert(t, err, "%s", tc.Src)
		})
	}
}

func Test_validateNoteType(t *testing.T) {
	testCases := getDeckNoteTypeTestCases()
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Src, func(t *testing.T) {
			err := validateNoteType(tc.Src)
			tc.ErrorAssert(t, err, "%s", tc.Src)
		})
	}
}

func getDeckNoteTypeTestCases() []struct {
	Src         string
	ErrorAssert assert.ErrorAssertionFunc
} {
	return []struct {
		Src         string
		ErrorAssert assert.ErrorAssertionFunc
	}{
		// positive
		{
			Src:         "a",
			ErrorAssert: assert.NoError,
		},
		{
			Src:         "ab",
			ErrorAssert: assert.NoError,
		},
		{
			Src:         "a  b",
			ErrorAssert: assert.NoError,
		},
		{
			Src:         "a#!b",
			ErrorAssert: assert.NoError,
		},
		{
			Src:         "a{b",
			ErrorAssert: assert.NoError,
		},
		{
			Src:         "a}b",
			ErrorAssert: assert.NoError,
		},
		{
			Src:         "a:b",
			ErrorAssert: assert.NoError,
		},
		// negative
		{
			Src:         "",
			ErrorAssert: assert.Error,
		},
		{
			Src:         " ",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "  ",
			ErrorAssert: assert.Error,
		},
		{
			Src:         " a",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "a ",
			ErrorAssert: assert.Error,
		},
		{
			Src:         " a ",
			ErrorAssert: assert.Error,
		},
		{
			Src:         " :",
			ErrorAssert: assert.Error,
		},
		{
			Src:         ": ",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "a\"b",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "a: ",
			ErrorAssert: assert.Error,
		},
		{
			Src:         " :a",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "a\nb",
			ErrorAssert: assert.Error,
		},
	}
}

func Test_validateFieldName(t *testing.T) {
	testCases := []struct {
		Src         string
		ErrorAssert assert.ErrorAssertionFunc
	}{
		// positive
		{
			Src:         "a",
			ErrorAssert: assert.NoError,
		},
		{
			Src:         "ab",
			ErrorAssert: assert.NoError,
		},
		{
			Src:         "a  b",
			ErrorAssert: assert.NoError,
		},
		{
			Src:         "a#!b",
			ErrorAssert: assert.NoError,
		},
		// negative
		{
			Src:         "",
			ErrorAssert: assert.Error,
		},
		{
			Src:         " ",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "  ",
			ErrorAssert: assert.Error,
		},
		{
			Src:         " a",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "a ",
			ErrorAssert: assert.Error,
		},
		{
			Src:         " a ",
			ErrorAssert: assert.Error,
		},
		{
			Src:         ":",
			ErrorAssert: assert.Error,
		},
		{
			Src:         " :",
			ErrorAssert: assert.Error,
		},
		{
			Src:         ": ",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "a:b",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "a}b",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "a{b",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "a\"b",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "a: ",
			ErrorAssert: assert.Error,
		},
		{
			Src:         " :a",
			ErrorAssert: assert.Error,
		},
		{
			Src:         "a\nb",
			ErrorAssert: assert.Error,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Src, func(t *testing.T) {
			err := validateFieldName(tc.Src)
			tc.ErrorAssert(t, err, "%s", tc.Src)
		})
	}
}
