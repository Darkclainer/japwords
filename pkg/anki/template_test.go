package anki

import (
	"fmt"
	"html/template"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestState struct {
	World string
}

func Test_TemplateTODO(t *testing.T) {
	tmpl := template.New("test")
	tmpl, err := tmpl.Parse(`
hello {{ .Worlr}}
`)
	require.NoError(t, err)
	tmpl = tmpl.Option("missingkey=error")
	err = tmpl.Execute(io.Discard, &TestState{})
	fmt.Printf("%T\n", err)
	require.NoError(t, err)
}
