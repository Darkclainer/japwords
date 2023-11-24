// lemma is like model package. It contains structure for jisho and wadoku
// packages and also can compose results from them.
package lemma

type Lemma struct {
	Slug   Word        `json:"Slug,omitempty"`
	Tags   []string    `json:"Tags,omitempty"`
	Forms  []Word      `json:"Forms,omitempty"`
	Senses []WordSense `json:"Senses,omitempty"`
	// Audio is array of links to audio files.
	// Key is format
	Audio map[string]string `json:"Audio,omitempty"`
}

type Word struct {
	Word     string   `json:"Word,omitempty"`
	Hiragana string   `json:"Hiragana,omitempty"`
	Furigana Furigana `json:"Furigana,omitempty"`
	// Pitches are encoded japanese pitch accent.
	// Every element of Pitches describe what accent (high or low)
	// should be used from previous element of Pitches up to and including
	// current specified position. Last elemnt can have virtual position
	// past length of Hiragana string, in that case it indicates
	// pitch accent of particle that follow word.
	//
	// For example consider world: 紙 「かみ」
	// The pitch accent of this word should be encoded next way:
	//
	// []Pitch{
	// 	{ Position: 3, IsHigh: false },
	// 	{ Position: 6, IsHigh: true  }
	// 	{ Position: 9, IsHigh: false },
	// }
	//
	// That means that 紙が will read as:
	// か (low) み (high) が (low)
	//
	// Note: Tokyo dialect can be encoded in much easier way, but I left
	// flexibility in case dictionary contain some entries that doesn't
	// follow Tokyo dialect rules.
	Pitches []Pitch `json:"Pitches,omitempty"`
}

func (w *Word) PitchShapes() []PitchShape {
	return ConvertPitchToShapes(w.Hiragana, w.Pitches)
}

func ConvertPitchToShapes(hiragana string, pitches []Pitch) []PitchShape {
	var pitchShapes []PitchShape
	// we copy first pitch, so we always think, that there was some pitch before current,
	// for invariant
	previousPitch := Pitch{
		Position: 0,
	}
	if len(pitches) > 0 {
		previousPitch.IsHigh = pitches[0].IsHigh
	}
	for i, currentPitch := range pitches {
		directions := []AccentDirection{convertBasePitch(currentPitch.IsHigh)}
		if previousPitch.Position == currentPitch.Position {
			if i == len(pitches)-1 && previousPitch.IsHigh != currentPitch.IsHigh {
				pitchShapes[i-1].Directions = append(pitchShapes[i-1].Directions, AccentDirectionRight)
			}
			continue
		}
		if previousPitch.IsHigh != currentPitch.IsHigh {
			directions = append(directions, AccentDirectionLeft)
		}
		pitchShapes = append(pitchShapes, PitchShape{
			Hiragana:   hiragana[previousPitch.Position:currentPitch.Position],
			Directions: directions,
		})
		previousPitch = currentPitch
	}
	if previousPitch.Position < len(hiragana)-1 {
		pitchShapes = append(pitchShapes, PitchShape{
			Hiragana: hiragana[previousPitch.Position:],
		})
	}
	return pitchShapes
}

type Pitch struct {
	Position int  `json:"Position,omitempty"`
	IsHigh   bool `json:"IsHigh,omitempty"`
}

type Furigana []FuriganaChar

type FuriganaChar struct {
	Kanji    string `json:"Kanji,omitempty"`
	Hiragana string `json:"Hiragana,omitempty"`
}

type WordSense struct {
	// Definition is slice of synonymous definitions in english
	Definition   []string `json:"Definition,omitempty"`
	PartOfSpeech []string `json:"PartOfSpeech,omitempty"`
	Tags         []string `json:"Tags,omitempty"`
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

func convertBasePitch(isHigh bool) AccentDirection {
	if isHigh {
		return AccentDirectionUp
	}
	return AccentDirectionDown
}
