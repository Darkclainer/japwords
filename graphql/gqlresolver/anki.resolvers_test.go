package gqlresolver

import (
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"

	"github.com/Darkclainer/japwords/graphql/gqlgenerated"
	"github.com/Darkclainer/japwords/graphql/gqlmodel"
	"github.com/Darkclainer/japwords/pkg/anki"
)

func Test_queryResolver_RenderFields(t *testing.T) {
	resolvers := Resolver{}
	c := client.New(handler.NewDefaultServer(gqlgenerated.NewExecutableSchema(gqlgenerated.Config{Resolvers: &resolvers})))
	type Response struct {
		RenderFields gqlmodel.RenderedFields
	}
	testCases := []struct {
		Name     string
		Query    string
		Expected gqlmodel.RenderedFields
	}{
		{
			Name: "Default template",
			Query: `
				query {
					RenderFields() {
						template
					}

				}`,
			Expected: gqlmodel.RenderedFields{
				Template: anki.GetDefaultExampleLemmaJSON(),
			},
		},
		{
			Name: "Pretty print template",
			Query: `
				query {
					RenderFields(template: "{\"Slug\": {\"Word\": \"Hello\"}}") {
						template
					}

				}`,
			Expected: gqlmodel.RenderedFields{
				Template: `{
  "Slug": {
    "Word": "Hello"
  }
}`,
			},
		},
		{
			Name: "Invalid template",
			Query: `
				query {
					RenderFields(template: "{a}") {
						templateError
					}

				}`,
			Expected: gqlmodel.RenderedFields{
				TemplateError: valuePointer("invalid character 'a' looking for beginning of object key string"),
			},
		},
		{
			Name: "Default template render",
			Query: `
				query {
					RenderFields(fields: ["{{.Slug.Word}}", "hello"]) {
						fields {
							field
							result
							error
						}
					}

				}`,
			Expected: gqlmodel.RenderedFields{
				Fields: []*gqlmodel.RenderedField{
					{
						Field:  "{{.Slug.Word}}",
						Result: anki.DefaultExampleLemma.Slug.Word,
					},
					{
						Field:  "hello",
						Result: "hello",
					},
				},
			},
		},
		{
			Name: "Invalid field",
			Query: `
				query {
					RenderFields(fields: ["{{.NotAField}}"]) {
						fields {
							field
							result
							error
						}
					}

				}`,
			Expected: gqlmodel.RenderedFields{
				Fields: []*gqlmodel.RenderedField{
					{
						Field:  "{{.NotAField}}",
						Result: "",
						Error:  valuePointer("template: :1:2: executing \"\" at <.NotAField>: can't evaluate field NotAField in type *lemma.ProjectedLemma"),
					},
				},
			},
		},
		{
			Name: "Render field with custom template",
			Query: `
				query {
					RenderFields(template: "{\"Slug\": {\"Word\": \"hello foo\"}}", fields: ["{{.Slug.Word}}"]) {
						fields {
							result
							error
						}
					}

				}`,
			Expected: gqlmodel.RenderedFields{
				Fields: []*gqlmodel.RenderedField{
					{
						Result: "hello foo",
					},
				},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			var resp Response
			c.MustPost(tc.Query, &resp)
			assert.Equal(t, tc.Expected, resp.RenderFields)
		})
	}
}

func valuePointer[T any](v T) *T {
	return &v
}
