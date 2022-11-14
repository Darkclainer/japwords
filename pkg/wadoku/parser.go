package wadoku

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"

	"japwords/pkg/lemma"
)

type reading struct {
	Hiragana string
	Pitches  []lemma.Pitch
}

func parseHTMLBytes(html []byte) ([]*lemma.WadokuLemma, error) {
	buffer := bytes.NewBuffer(html)
	return parseHTML(buffer)
}

var contentSectionMatcher = singleMatcher("section#content")

func parseHTML(src io.Reader) ([]*lemma.WadokuLemma, error) {
	document, err := goquery.NewDocumentFromReader(src)
	if err != nil {
		return nil, fmt.Errorf("can not parse page: %w", err)
	}

	contentSection := document.FindMatcher(contentSectionMatcher)
	// we consider page correct only if it contains one node in selection
	// contentSectionMatcher
	if contentSection.Length() != 1 {
		return nil, errors.New("page doesn't contain content section or contains more than one")
	}
	return parseContentSection(contentSection)
}

var (
	tableResultBodyMatcher = singleMatcher("#resulttable > tbody")
	rowResultMatcher       = matcher("tr.resultline")
)

func parseContentSection(sel *goquery.Selection) ([]*lemma.WadokuLemma, error) {
	tableBody := sel.FindMatcher(tableResultBodyMatcher)
	rows := tableBody.ChildrenMatcher(rowResultMatcher)
	if rows.Length() == 0 {
		return nil, nil
	}
	var lemmas []*lemma.WadokuLemma
	var errs []*LemmaError
	rows.Each(func(i int, row *goquery.Selection) {
		newLemmas, err := parseRowResult(row)
		if err != nil {
			errs = append(
				errs,
				&LemmaError{
					ID:  i,
					Err: err,
				},
			)
		}
		if len(newLemmas) != 0 {
			lemmas = append(lemmas, newLemmas...)
		}
	})
	if len(errs) != 0 {
		return lemmas, &LemmaBatchError{Errs: errs}
	}
	return lemmas, nil
}

var (
	resultDetailMatcher = singleMatcher(".resultdetail")
	japaneseMatcher     = singleMatcher(".japanese")
	readingMather       = singleMatcher(".reading")
)

func parseRowResult(sel *goquery.Selection) ([]*lemma.WadokuLemma, error) {
	resultDetail := sel.ChildrenMatcher(resultDetailMatcher)
	japaneseSel := resultDetail.ChildrenMatcher(japaneseMatcher)
	japaneseVariants, err := parseJapanese(japaneseSel)
	if err != nil {
		return nil, err
	}
	if len(japaneseVariants) == 0 {
		// we filtered out all results, that's not an error
		return nil, nil
	}
	readingSel := resultDetail.ChildrenMatcher(readingMather)
	reading, err := parseReading(readingSel)
	if err != nil {
		return nil, err
	}
	if reading == nil {
		// we also don't need result with filtered readings
		return nil, nil
	}
	lemmas := make([]*lemma.WadokuLemma, 0, len(japaneseVariants))
	for _, variant := range japaneseVariants {
		lemmas = append(lemmas, &lemma.WadokuLemma{
			Slug:     variant,
			Hiragana: reading.Hiragana,
			Pitches:  reading.Pitches,
		})
	}
	return lemmas, nil
}

var orthMatcher = singleMatcher(".orth")

func parseJapanese(sel *goquery.Selection) ([]string, error) {
	orthChild := sel.FindMatcher(orthMatcher).Contents()
	var buffer strings.Builder
	var variants []string
	for i, node := range orthChild.Nodes {
		element := orthChild.Eq(i)
		switch node.Type {
		case html.TextNode:
			buffer.WriteString(strings.TrimSpace(element.Text()))
		case html.ElementNode:
			if element.HasClass("divider") {
				variants = append(variants, buffer.String())
				buffer.Reset()
				continue
			}
			buffer.WriteString(strings.TrimSpace(element.Text()))
		}
	}
	if buffer.Len() != 0 {
		variants = append(variants, buffer.String())
	}
	if len(variants) == 0 {
		// this probably indication of bug for us
		return nil, errors.New("no japanese slug found")
	}
	filtered := variants[:0]
	for _, variant := range variants {
		if strings.HasPrefix(variant, "…") || strings.HasSuffix(variant, "…") {
			continue
		}
		filtered = append(filtered, variant)
	}
	if len(filtered) == 0 {
		return nil, nil
	}
	return filtered, nil
}

var pronAccentMatcher = singleMatcher(".pron.accent")

func parseReading(sel *goquery.Selection) (*reading, error) {
	accentParts := sel.FindMatcher(pronAccentMatcher).Children()
	if accentParts.Length() == 0 {
		return nil, nil
	}
	var (
		buffer   strings.Builder
		pitches  []lemma.Pitch
		position = 0
	)
	accentParts.Each(func(_ int, s *goquery.Selection) {
		readingPart := extractReading(s)
		n, _ := buffer.WriteString(readingPart)
		position += n
		pitches = append(pitches, lemma.Pitch{
			Position: position,
			IsHigh:   s.HasClass("t"), // `t` stands for top
		})
	})
	if accentParts.Last().HasClass("r") {
		pitches = append(pitches, lemma.Pitch{
			Position: position,
			IsHigh:   !pitches[len(pitches)-1].IsHigh,
		})
	}
	if buffer.Len() == 0 {
		// that means we wasn't able to extract actual reading, I see this as a bug
		return nil, errors.New("no reading found despite section with accent was")
	}
	return &reading{
		Hiragana: buffer.String(),
		Pitches:  pitches,
	}, nil
}

// extractReading extracts reading from pitched span (one that contains
// classes like `t` or `b` and removes extra parts like `|` or `･`
func extractReading(sel *goquery.Selection) string {
	contents := sel.Contents()
	var buffer strings.Builder
	for i, node := range contents.Nodes {
		element := contents.Eq(i)
		switch node.Type {
		case html.TextNode:
			trimmed := strings.TrimSpace(element.Text())
			buffer.WriteString(
				strings.ReplaceAll(trimmed, "･", ""),
			)
		case html.ElementNode:
			continue
		}
	}
	return buffer.String()
}

func matcher(src string) goquery.Matcher {
	return cascadia.MustCompile(src)
}

func singleMatcher(src string) goquery.Matcher {
	return goquery.SingleMatcher(matcher(src))
}
