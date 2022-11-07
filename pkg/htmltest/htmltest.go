package htmltest

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

func MustDocument(t *testing.T, src string) *goquery.Document {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(src))
	require.NoError(t, err)
	return doc
}

func MustRootSelection(t *testing.T, src string) *goquery.Selection {
	t.Helper()
	doc := MustDocument(t, src)
	sel := doc.Find("#root")
	if sel.Length() != 1 {
		t.Fatal("document should contain one #root node")
	}
	return sel
}
