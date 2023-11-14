package anki

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/Darkclainer/japwords/pkg/lemma"
)

type Template struct {
	Src  string
	Tmpl *template.Template
}

type TemplateMapping map[string]*Template

func (tm TemplateMapping) Equal(otm TemplateMapping) bool {
	if len(tm) != len(otm) {
		return false
	}
	for field, tmpl := range tm {
		otherTmpl, ok := otm[field]
		if !ok || otherTmpl.Src != tmpl.Src {
			return false
		}
	}
	return true
}

// RenderRawTemplate renders template with specified lemma, intended to use for API.
func RenderRawTemplate(templateSrc string, lemma *Lemma) (string, error) {
	tmpl := template.New("")
	err := initTemplate(tmpl, templateSrc)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, lemma)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func convertMapping(mapping map[string]string) (TemplateMapping, []*MappingValidationError) {
	root := template.New("")
	result := map[string]*Template{}
	var errs []*MappingValidationError
	for field, src := range mapping {
		fieldTemplate := root.New(field)
		err := initTemplate(fieldTemplate, src)
		if err != nil {
			errs = append(errs, &MappingValidationError{
				Key: field,
				Msg: err.Error(),
			})
		}
		result[field] = &Template{
			Src:  src,
			Tmpl: fieldTemplate,
		}
	}
	return result, errs
}

var (
	templateFuncsCached template.FuncMap
	templateFuncsSync   sync.Once
)

func templateFuncs() template.FuncMap {
	templateFuncsSync.Do(func() {
		sprigFuncs := sprig.TxtFuncMap()
		newFuncs := map[string]any{
			"renderFurigana": renderFuriganaTemplate,
			"renderPitch":    renderPitchTemplate,
		}
		for name, f := range sprigFuncs {
			newFuncs[name] = f
		}
		templateFuncsCached = newFuncs
	})
	return templateFuncsCached
}

func initTemplate(tmpl *template.Template, src string) error {
	tmpl.Funcs(templateFuncs())
	_, err := tmpl.Parse(src)
	if err != nil {
		return err
	}
	tmpl.Option("missingkey=error")
	return checkTemplate(tmpl)
}

func checkTemplate(tmpl *template.Template) error {
	lemma := &Lemma{}
	return tmpl.Execute(io.Discard, lemma)
}

// renderFuriganaTemplate is template functions that return string representation of
// furigana that can be understood by Anki.
func renderFuriganaTemplate(word *Word) string {
	furigana := word.Furigana
	var buffer strings.Builder
	for _, part := range furigana {
		// if only hiragana presented, print it. Else use special syntax, supported by anki
		kanji, hiragana := part.Kanji != "", part.Hiragana != ""
		if (kanji || hiragana) && !(kanji && hiragana) {
			// either of them has value
			buffer.WriteString(part.Kanji)
			buffer.WriteString(part.Hiragana)
		} else {
			buffer.WriteString(part.Kanji)
			buffer.WriteByte('[')
			buffer.WriteString(part.Hiragana)
			buffer.WriteByte(']')
		}
	}
	return buffer.String()
}

func renderPitchTemplate(word *Word, tag string, up string, right string, down string, left string) (string, error) {
	return renderPitch(word, tag, []string{
		up,
		right,
		down,
		left,
	})
}

func renderPitch(word *Word, tag string, directionClasses []string) (string, error) {
	if tag == "" {
		return "", errors.New("tag should be non empty string")
	}
	if len(directionClasses) < 4 {
		return "", errors.New("renderPitch should be called with 4 direction classes")
	}
	pitchShapes := word.Pitches
	var buffer strings.Builder
	for _, shape := range pitchShapes {
		buffer.WriteByte('<')
		buffer.WriteString(tag)
		buffer.WriteString(" class=\"")
		for i, direction := range shape.Directions {
			var directionString string
			switch direction {
			case lemma.AccentDirectionUp:
				directionString = directionClasses[0]
			case lemma.AccentDirectionRight:
				directionString = directionClasses[1]
			case lemma.AccentDirectionDown:
				directionString = directionClasses[2]
			case lemma.AccentDirectionLeft:
				directionString = directionClasses[3]
			default:
				panic("unreachable")
			}
			buffer.WriteString(directionString)
			if i+1 < len(shape.Directions) {
				buffer.WriteByte(' ')
			}
		}
		buffer.WriteString("\">")
		buffer.WriteString(shape.Hiragana)
		buffer.WriteString("</")
		buffer.WriteString(tag)
		buffer.WriteByte('>')
	}
	return buffer.String(), nil
}
