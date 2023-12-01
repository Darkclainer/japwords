package query

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Render(t *testing.T) {
	testCases := []struct {
		Name     string
		Query    Query
		Expected string
	}{
		{
			Name:     "exact",
			Query:    Exact("hello", "world"),
			Expected: `"hello:world"`,
		},
		{
			Name:     "and",
			Query:    And(Exact("hello", "a"), Exact("world", "b")),
			Expected: `("hello:a" "world:b")`,
		},
		{
			Name:     "or",
			Query:    Or(Exact("hello", "a"), Exact("world", "b")),
			Expected: `("hello:a" OR "world:b")`,
		},
		{
			Name: "complex",
			Query: And(
				Exact("deck", "my&deck"),
				Exact("note", "my< note"),
				Or(
					Exact("field*1", "va<>lu e"),
					Exact("field*2", "value"),
				),
			),
			Expected: `("deck:my&deck" "note:my< note" ("field\*1:va&lt;&gt;lu e" OR "field\*2:value"))`,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			actual := Render(tc.Query)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_BinaryQuery_write(t *testing.T) {
	testCases := []struct {
		Name     string
		Query    *BinaryQuery
		Expected string
	}{
		{
			Name: "simple and",
			Query: &BinaryQuery{
				Operation: AndOp,
				Args: []Query{
					Exact("hello", "a"),
					Exact("world", "b"),
					Exact("foo", "bar"),
				},
			},
			Expected: `("hello:a" "world:b" "foo:bar")`,
		},
		{
			Name: "reduced and",
			Query: &BinaryQuery{
				Operation: AndOp,
				Args: []Query{
					Exact("hello", "a"),
				},
			},
			Expected: `"hello:a"`,
		},
		{
			Name: "simple or",
			Query: &BinaryQuery{
				Operation: OrOp,
				Args: []Query{
					Exact("hello", "a"),
					Exact("world", "b"),
					Exact("foo", "bar"),
				},
			},
			Expected: `("hello:a" OR "world:b" OR "foo:bar")`,
		},
		{
			Name: "reduced or",
			Query: &BinaryQuery{
				Operation: OrOp,
				Args: []Query{
					Exact("hello", "a"),
				},
			},
			Expected: `"hello:a"`,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			var buffer bytes.Buffer
			tc.Query.write(&buffer)
			assert.Equal(t, tc.Expected, buffer.String())
		})
	}
}

func Test_ExactQuery_write(t *testing.T) {
	testCases := []struct {
		Name     string
		Query    *ExactQuery
		Expected string
	}{
		{
			Name: "no escape",
			Query: &ExactQuery{
				Field: "hello",
				Value: "world",
			},
			Expected: `"hello:world"`,
		},
		{
			Name: "no field",
			Query: &ExactQuery{
				Value: "world",
			},
			Expected: `"world"`,
		},
		{
			Name: "no field escape value",
			Query: &ExactQuery{
				Value: `w"or&l<d`,
			},
			Expected: `"w\"or&amp;l&lt;d"`,
		},
		{
			Name: "escape field",
			Query: &ExactQuery{
				Field: `h_e&l*l:o`,
				Value: "world",
			},
			Expected: `"h\_e&l\*l\:o:world"`,
		},
		{
			Name: "escape field and value",
			Query: &ExactQuery{
				Field: `h_e&l*l:o`,
				Value: "wo<rl*d",
			},
			Expected: `"h\_e&l\*l\:o:wo&lt;rl\*d"`,
		},
		{
			Name: "escape for deck",
			Query: &ExactQuery{
				Field: `deck`,
				Value: "wo<rl*d",
			},
			Expected: `"deck:wo<rl\*d"`,
		},
		{
			Name: "escape for note",
			Query: &ExactQuery{
				Field: `note`,
				Value: "wo<rl*d",
			},
			Expected: `"note:wo<rl\*d"`,
		},
		{
			Name: "escape for tag",
			Query: &ExactQuery{
				Field: `tag`,
				Value: "wo<rl*d",
			},
			Expected: `"tag:wo<rl\*d"`,
		},
		{
			Name: "escape for card",
			Query: &ExactQuery{
				Field: `card`,
				Value: "wo<rl*d",
			},
			Expected: `"card:wo<rl\*d"`,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			var buffer bytes.Buffer
			tc.Query.write(&buffer)
			assert.Equal(t, tc.Expected, buffer.String())
		})
	}
}

func Test_escape(t *testing.T) {
	testCases := []struct {
		Src      string
		Expected string
	}{
		{
			Src:      "nothing to escape",
			Expected: "nothing to escape",
		},
		{
			Src:      `escape all "*_:\`,
			Expected: `escape all \"\*\_\:\\`,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Src, func(t *testing.T) {
			var buffer bytes.Buffer
			escape(&buffer, tc.Src)
			assert.Equal(t, tc.Expected, buffer.String())
		})
	}
}

func Test_escapeField(t *testing.T) {
	testCases := []struct {
		Src      string
		Expected string
	}{
		{
			Src:      "nothing to escape",
			Expected: "nothing to escape",
		},
		{
			Src:      `escape all "*_:\`,
			Expected: `escape all \"\*\_\:\\`,
		},
		{
			Src:      "escape for fields &<>",
			Expected: "escape for fields &amp;&lt;&gt;",
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Src, func(t *testing.T) {
			var buffer bytes.Buffer
			escapeField(&buffer, tc.Src)
			assert.Equal(t, tc.Expected, buffer.String())
		})
	}
}
