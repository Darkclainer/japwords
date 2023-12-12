// lemma is like model package. It contains structure for jisho and wadoku
// packages and also can compose results from them.
package lemma

type Lemma struct {
	Slug   Word        `json:"Slug,omitempty"`
	Tags   []string    `json:"Tags,omitempty"`
	Forms  []Word      `json:"Forms,omitempty"`
	Senses []WordSense `json:"Senses,omitempty"`
	Audio  []Audio     `json:"Audio,omitempty"`
}

type Word struct {
	Word     string   `json:"Word,omitempty"`
	Hiragana string   `json:"Hiragana,omitempty"`
	Furigana Furigana `json:"Furigana,omitempty"`
	// PitchShapes is encoded pitch for moras
	// TODO: implement normalization: left direction is prefered,
	// so we say generally say that (a|)(b) should be written as (a)(|b) or
	// (AccentDirectionUD1 AccentDirectionRight)-(AccentDirectionUD2 should be substitute by
	// (AccentDirectionUD1)-(AccentDirectionUD2 AccentDirectionLeft)
	PitchShapes []PitchShape `json:"Pitches,omitempty"`
}

type Audio struct {
	// MediaType is one of the specified audo media types by iana:
	// https://www.iana.org/assignments/media-types/media-types.xhtml#audio
	MediaType string `json:"Format,omitempty"`
	Source    string `json:"Source,omitempty"`
}

type Furigana []FuriganaChar

type FuriganaChar struct {
	Kanji    string `json:"Kanji,omitempty"`
	Hiragana string `json:"Hiragana,omitempty"`
}

//go:generate $ENUMER_TOOL -type=AccentDirection -trimprefix=AccentDirection -transform=upper -text -gqlgen
type AccentDirection int

const (
	AccentDirectionUp AccentDirection = iota
	AccentDirectionRight
	AccentDirectionDown
	AccentDirectionLeft
)

type PitchShape struct {
	Hiragana   string            `json:"Hiragana,omitempty"`
	Directions []AccentDirection `json:"Directions,omitempty"`
}

type WordSense struct {
	// Definition is slice of synonymous definitions in english
	Definition   []string `json:"Definition,omitempty"`
	PartOfSpeech []string `json:"PartOfSpeech,omitempty"`
	Tags         []string `json:"Tags,omitempty"`
}

// ProjectedLemma is more specific variant of Lemma.
// This structure include only one meaning.
type ProjectedLemma struct {
	Slug          Word     `json:"Slug,omitempty"`
	Tags          []string `json:"Tags,omitempty"`
	Forms         []Word   `json:"Forms,omitempty"`
	Definitions   []string `json:"Definitions,omitempty"`
	PartsOfSpeech []string `json:"PartsOfSpeech,omitempty"`
	SenseTags     []string `json:"SenseTags,omitempty"`
	Audio         []Audio  `json:"Audio,omitempty"`
}
