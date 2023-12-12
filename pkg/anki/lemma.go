package anki

import (
	"encoding/json"
	"sync"

	"github.com/Darkclainer/japwords/pkg/lemma"
)

var DefaultExampleLemma = lemma.ProjectedLemma{
	Slug: lemma.Word{
		Word:     "一二わ三はい",
		Hiragana: "いちにわさんはい",
		Furigana: []lemma.FuriganaChar{
			{
				Kanji:    "一",
				Hiragana: "いち",
			},
			{
				Kanji:    "二",
				Hiragana: "に",
			},
			{
				Hiragana: "わ",
			},
			{
				Kanji:    "三",
				Hiragana: "さん",
			},
			{
				Hiragana: "は",
			},
			{
				Hiragana: "い",
			},
		},
		PitchShapes: []lemma.PitchShape{
			{
				Hiragana: "いち",
				Directions: []lemma.AccentDirection{
					lemma.AccentDirectionDown,
				},
			},
			{
				Hiragana: "わ",
				Directions: []lemma.AccentDirection{
					lemma.AccentDirectionUp,
					lemma.AccentDirectionLeft,
				},
			},
			{
				Hiragana: "さんは",
				Directions: []lemma.AccentDirection{
					lemma.AccentDirectionDown,
					lemma.AccentDirectionLeft,
				},
			},
			{
				Hiragana: "い",
				Directions: []lemma.AccentDirection{
					lemma.AccentDirectionUp,
					lemma.AccentDirectionLeft,
					lemma.AccentDirectionRight,
				},
			},
		},
	},
	Tags: []string{
		"Common word",
		"JLPT N5",
		"Wanikani level 2",
		"Test database",
	},
	Forms: []lemma.Word{},
	Definitions: []string{
		"this is test lemma, it means nothing, just useful for tests",
		"use it only for tests",
		"don't use it for anything else",
	},
	PartsOfSpeech: []string{
		"Expressions (phrases, clauses, etc.)",
		"Test noun (probably verb)",
		"I-adjective (keiyoushi)",
	},
	SenseTags: []string{
		"Test",
		"Nonsense",
		"Usually written using kana alone",
	},
	Audio: []lemma.Audio{
		{MediaType: "audio/mpeg", Source: "https://example.com/somelink/mp3"},
		{MediaType: "audio/ogg", Source: "https://example.com/somelink/ogg"},
	},
}

var GetDefaultExampleLemmaJSON = sync.OnceValue(func() string {
	src, err := json.MarshalIndent(&DefaultExampleLemma, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(src)
})
