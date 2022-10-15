package jisho

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
	"go.uber.org/multierr"
)

func parseHTMLBytes(html []byte) ([]*Lemma, error) {
	buffer := bytes.NewBuffer(html)
	return parseHTML(buffer)
}

// More complex selector is needed to ensure that we are parsing correct page
var mainResultsMatcher = singleMatcher("#page_container #main_results")

func parseHTML(src io.Reader) ([]*Lemma, error) {
	document, err := goquery.NewDocumentFromReader(src)
	if err != nil {
		return nil, fmt.Errorf("can not parse page: %w", err)
	}

	mainResults := document.FindMatcher(mainResultsMatcher)
	// we consider page correct only if it contains one node in selection
	// mainResultsMatcher
	if mainResults.Length() != 1 {
		return nil, errors.New("page doesn't contain main_results or contains more than one")
	}
	return parseMainResults(mainResults.Eq(0))
}

var (
	primaryMatcher      = singleMatcher("#primary")
	conceptLightMatcher = matcher("div.concept_light")
)

func parseMainResults(sel *goquery.Selection) ([]*Lemma, error) {
	primary := sel.FindMatcher(primaryMatcher)
	conceptLight := primary.Children().ChildrenMatcher(conceptLightMatcher)
	var (
		lemmas []*Lemma
		errs   []error
	)
	conceptLight.Each(func(i int, sel *goquery.Selection) {
		lemma, err := parseConceptLight(sel)
		if err != nil {
			errs = append(
				errs,
				fmt.Errorf("failed to process concept %d: %w", i, err),
			)
		}
		if lemma != nil {
			lemmas = append(lemmas, lemma)
		}
	})
	return lemmas, multierr.Combine(errs...)
}

var (
	conceptLightWrapperMatcher        = singleMatcher("div.concept_light-wrapper")
	conceptLightRepresentationMatcher = singleMatcher("div.concept_light-representation")
	conceptLightStatusMatcher         = singleMatcher("div.concept_light-status")
)

func parseConceptLight(sel *goquery.Selection) (*Lemma, error) {
	var lemma Lemma
	lightWrapper := sel.FindMatcher(conceptLightWrapperMatcher)
	representation := lightWrapper.FindMatcher(conceptLightRepresentationMatcher)
	slug, err := parseRepresentation(representation)
	if err != nil {
		return nil, err
	}
	lemma.Slug = slug
	status := lightWrapper.ChildrenMatcher(conceptLightWrapperMatcher)
	audio, tags := parseStatus(status)
	lemma.Audio = audio
	lemma.Tags = tags
	return &lemma, nil
}

var (
	textMatcher     = singleMatcher(".text")
	furiganaMatcher = singleMatcher(".furigana")
)

func parseRepresentation(sel *goquery.Selection) (Word, error) {
	text := strings.TrimSpace(sel.FindMatcher(textMatcher).Text())
	if text == "" {
		return Word{}, errors.New("no text representation found for word")
	}
	var furigana Furigana
	var reading strings.Builder
	furiganaSpans := sel.FindMatcher(furiganaMatcher).Children()
	i := 0
	for _, r := range text {
		f := strings.TrimSpace(furiganaSpans.Eq(i).Text())
		rChar := string(r)
		furigana = append(furigana, FuriganaChar{
			Kanji:    rChar,
			Hiragana: f,
		})
		if f != "" {
			rChar = f
		}
		reading.WriteString(rChar)
		i++
	}
	return Word{
		Word:     text,
		Furigana: furigana,
		Reading:  reading.String(),
	}, nil
}

var conceptLightTagMatcher = matcher("span.concept_light-tag")

func parseStatus(sel *goquery.Selection) (audio map[string]string, tags []string) {
	tagSel := sel.ChildrenMatcher(conceptLightTagMatcher)
	tagSel.Each(func(_ int, sel *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(sel.Text()))
	})
	sourceSel := sel.ChildrenFiltered("audio").Children()
	audio = make(map[string]string)
	for _, source := range sourceSel.Nodes {
		var audioSrc, audioType string
		for _, attr := range source.Attr {
			switch attr.Key {
			case "src":
				audioSrc = attr.Val
			case "type":
				audioType = attr.Val
			}
		}
		if audioSrc != "" {
			if audioType == "" {
				audioType = "unknown"
			}
			audio[audioType] = audioSrc
		}
	}
	return
}

func matcher(src string) goquery.Matcher {
	return cascadia.MustCompile(src)
}

func singleMatcher(src string) goquery.Matcher {
	return goquery.SingleMatcher(matcher(src))
}
