package jisho

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"

	"github.com/Darkclainer/japwords/pkg/lemma"
)

func parseHTMLBytes(html []byte) ([]*lemma.Lemma, error) {
	buffer := bytes.NewBuffer(html)
	return parseHTML(buffer)
}

// More complex selector is needed to ensure that we are parsing correct page
var mainResultsMatcher = singleMatcher("#page_container #main_results")

func parseHTML(src io.Reader) ([]*lemma.Lemma, error) {
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

func parseMainResults(sel *goquery.Selection) ([]*lemma.Lemma, error) {
	primary := sel.FindMatcher(primaryMatcher)
	conceptLight := primary.Children().ChildrenMatcher(conceptLightMatcher)
	var (
		lemmas []*lemma.Lemma
		errs   []error
	)
	conceptLight.Each(func(i int, sel *goquery.Selection) {
		lemma, err := parseConceptLight(sel)
		if err != nil {
			errs = append(
				errs,
				&LemmaError{
					ID:  i,
					Err: err,
				},
			)
		}
		if lemma != nil {
			lemmas = append(lemmas, lemma)
		}
	})
	if len(errs) != 0 {
		return lemmas, &LemmaBatchError{Errs: errs}
	}
	return lemmas, nil
}

var (
	conceptLightWrapperMatcher        = singleMatcher("div.concept_light-wrapper")
	conceptLightRepresentationMatcher = singleMatcher("div.concept_light-representation")
	conceptLightStatusMatcher         = singleMatcher("div.concept_light-status")
	conceptLightMeaningsMatcher       = singleMatcher("div.concept_light-meanings")
)

func parseConceptLight(sel *goquery.Selection) (*lemma.Lemma, error) {
	lightWrapper := sel.ChildrenMatcher(conceptLightWrapperMatcher)
	representation := lightWrapper.FindMatcher(conceptLightRepresentationMatcher)
	slug, err := parseRepresentation(representation)
	if err != nil {
		return nil, err
	}
	longRepresentation := lightWrapper.HasClass("concept_light-long_representation")
	var status *goquery.Selection
	if longRepresentation {
		status = lightWrapper.NextMatcher(conceptLightStatusMatcher)
	} else {
		status = lightWrapper.ChildrenMatcher(conceptLightStatusMatcher)
	}
	audio, tags := parseStatus(status)

	var meanings *goquery.Selection
	if longRepresentation {
		meanings = sel.ChildrenMatcher(conceptLightMeaningsMatcher).Children()
	} else {
		meanings = lightWrapper.NextMatcher(conceptLightMeaningsMatcher).Children()
	}
	wordSenses, otherForms := parseMeanings(meanings)

	return &lemma.Lemma{
		Slug:   slug,
		Tags:   tags,
		Forms:  otherForms,
		Senses: wordSenses,
		Audio:  audio,
	}, nil
}

var (
	textMatcher     = singleMatcher(".text")
	furiganaMatcher = singleMatcher(".furigana")
)

func parseRepresentation(sel *goquery.Selection) (lemma.Word, error) {
	text := strings.TrimSpace(sel.FindMatcher(textMatcher).Text())
	if text == "" {
		return lemma.Word{}, errors.New("no text representation found for word")
	}
	var furigana lemma.Furigana
	var reading strings.Builder
	furiganaSpans := sel.FindMatcher(furiganaMatcher).Children()
	i := 0
	for _, r := range text {
		f := strings.TrimSpace(furiganaSpans.Eq(i).Text())
		rChar := string(r)
		if f != "" {
			furigana = append(furigana, lemma.FuriganaChar{
				Kanji:    rChar,
				Hiragana: f,
			})
			rChar = f
		} else {
			furigana = append(furigana, lemma.FuriganaChar{
				Hiragana: rChar,
			})
		}
		reading.WriteString(rChar)
		i++
	}
	// text is equal to reading only when for no furigana given for any
	// character in text
	if text == reading.String() {
		furigana = nil
		reading.Reset()

	}
	return lemma.Word{
		Word:     text,
		Furigana: furigana,
		Hiragana: reading.String(),
	}, nil
}

var conceptLightTagMatcher = matcher(".concept_light-tag")

func parseStatus(sel *goquery.Selection) (audio []lemma.Audio, tags []string) {
	tagSel := sel.ChildrenMatcher(conceptLightTagMatcher)
	tagSel.Each(func(_ int, sel *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(sel.Text()))
	})
	sourceSel := sel.ChildrenFiltered("audio").Children()
	if sourceSel.Length() == 0 {
		return
	}
	for _, source := range sourceSel.Nodes {
		var audioSrc, audioType string
		for _, attr := range source.Attr {
			switch attr.Key {
			case "src":
				audioSrc = attr.Val
				if strings.HasPrefix(audioSrc, "//") {
					// TODO: for better handling of this
					// we need to introduce state for parse* functions
					// and store html location to resolve
					// relative links properly.
					audioSrc = "https:" + audioSrc
				}

			case "type":
				audioType = attr.Val
			}
		}
		if audioSrc != "" {
			if audioType == "" {
				audioType = "unknown"
			}
			audio = append(audio, lemma.Audio{
				Type:   audioType,
				Source: audioSrc,
			})
		}
	}
	return
}

var meaningDefinitionMatcher = singleMatcher("div.meaning-definition")

// parseMeanings return senses and other forms from `.meaning-wrapper`
func parseMeanings(sel *goquery.Selection) (senses []lemma.WordSense, forms []lemma.Word) {
	// get all tag-meaning pairs
	sel = sel.Children().First()
	for sel.Length() != 0 {
		var partOfSpeech []string
		if sel.HasClass("meaning-tags") {
			tagsText := strings.TrimSpace(sel.Text())
			if tagsText == "Other forms" {
				sel = sel.Next()
				break
			}
			partOfSpeech = splitPartOfSpeech(tagsText)
			sel = sel.Next()
		}
		if sel.HasClass("meaning-tags") {
			continue
		}
		if sel.HasClass("meaning-wrapper") {
			meaningDefinition := sel.ChildrenMatcher(meaningDefinitionMatcher)
			definitions, tags := parseMeaningDefinition(meaningDefinition)
			if len(definitions) > 0 {
				senses = append(senses, lemma.WordSense{
					Definition:   definitions,
					PartOfSpeech: partOfSpeech,
					Tags:         tags,
				})
			}
		}
		sel = sel.Next()
	}
	// no Other forms provided
	if sel.Length() == 0 || !sel.HasClass("meaning-wrapper") {
		return
	}
	forms = parseDefinitionOtherForms(
		sel.ChildrenMatcher(meaningDefinitionMatcher),
	)
	return
}

var (
	meaningMeaningMatcher   = singleMatcher(".meaning-meaning")
	supplementalInfoMatcher = singleMatcher(".supplemental_info")
	tagTagMatcher           = matcher(".tag-tag")
)

func parseMeaningDefinition(sel *goquery.Selection) (definitions []string, tags []string) {
	rawDefinition := strings.TrimSpace(sel.ChildrenMatcher(meaningMeaningMatcher).Text())
	for _, definition := range strings.Split(rawDefinition, "; ") {
		definition = strings.TrimSpace(definition)
		if definition == "" {
			continue
		}
		definitions = append(definitions, definition)
	}
	tagsSel := sel.
		ChildrenMatcher(supplementalInfoMatcher).
		ChildrenMatcher(tagTagMatcher)
	tags = tagsSel.Map(func(_ int, s *goquery.Selection) string {
		return strings.TrimSpace(s.Text())
	})
	return
}

var breakUnitMatcher = matcher(".break-unit")

// parseMeainingOtherForms extracts other forms from last `.meaning-definition`
func parseDefinitionOtherForms(sel *goquery.Selection) (forms []lemma.Word) {
	meaningSel := sel.ChildrenMatcher(meaningMeaningMatcher)
	breakUnits := meaningSel.ChildrenMatcher(breakUnitMatcher)
	breakUnits.Each(func(_ int, s *goquery.Selection) {
		raw := strings.TrimSpace(s.Text())
		word, exist := parseOtherForm(raw)
		if exist {
			forms = append(forms, word)
		}
	})
	return
}

var partOfSpeechDelimiter = regexp.MustCompile(`, \p{Lu}`)

func splitPartOfSpeech(src string) []string {
	matches := partOfSpeechDelimiter.FindAllStringIndex(src, -1)
	var result []string
	start := 0
	for _, match := range matches {
		if match[0]-start > 0 {
			result = append(result, src[start:match[0]])
		}
		// we want uppercase letter to get in our next
		// selection
		start = match[1] - 1
	}
	if len(src) > start {
		result = append(result, src[start:])
	}
	return result
}

var otherFormRegex = regexp.MustCompile(`([^【\t\n\f\r ]+)(?:\s*【([^】]+)】)?`)

func parseOtherForm(src string) (lemma.Word, bool) {
	matches := otherFormRegex.FindStringSubmatch(src)
	if len(matches) == 0 {
		return lemma.Word{}, false
	}
	return lemma.Word{
		Word:     matches[1],
		Hiragana: matches[2],
	}, true
}

func matcher(src string) goquery.Matcher {
	return cascadia.MustCompile(src)
}

func singleMatcher(src string) goquery.Matcher {
	return goquery.SingleMatcher(matcher(src))
}
