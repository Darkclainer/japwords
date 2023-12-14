package anki

import "github.com/Darkclainer/japwords/pkg/anki/ankiconnect"

// defaultCreateModelRequest returns request for create default note type.
func defaultCreateModelRequest() *ankiconnect.CreateModelRequest {
	return &ankiconnect.CreateModelRequest{
		Fields: []string{
			"Sort",
			"Kanji",
			"Furigana",
			"Kana",
			"PoS",
			"English",
			"SenseTags",
			"Notes",
			"Audio",
			"Example",
		},
		CSS: `.card {
  --color-text-main: black;
  --color-text-secondary: #808080;
  --color-background: white;
  --color-kana-line: #808080;
  --color-pos-background: #6af3ff;
  --color-tag-background: #ff9288;
}
.card.night_mode {
  --color-text-main: white;
  --color-text-secondary: #808080;
  --color-background: black;
  --color-kana-line: #808080;
  --color-pos-background: #007b86;
  --color-tag-background: #860b00;
}
.card {
  font-family: "Inter", sans-serif;
  font-size: 1rem;
  padding: 0 1rem;
  margin: 0 0;
  text-align: center;
  color: var(--color-text-main);
  background-color: var(--color-background);
  word-wrap: break-word;
}
.card.night_mode {
  color: var(--color-text-main);
  background-color: var(--color-background);
}

.divider {
  margin: 1rem 0;
  height: 1px;
  background-color: var(--color-text-secondary);
}

.kanji-block {
  margin: 2.5rem 0 0.5rem;
  font-size: 3rem;
}
.kana-block {
  margin-bottom: 2.5rem;
  font-size: 1.125rem;
  color: var(--color-text-secondary);
}
.kana-block span {
  border-color: var(--color-text-secondary);
  border-width: 0px;
  border-style: dashed;
}
.kana-block .border-u {
  border-top-width: 1px;
}
.kana-block .border-d {
  border-bottom-width: 1px;
}
.kana-block .border-r {
  border-right-width: 1px;
}
.kana-block .border-l {
  border-left-width: 1px;
}

.pill-block {
  display: flex;
  flex-direction: row;
  margin: 2.5rem 0;
  justify-content: center;
  flex-wrap: wrap;
  font-size: 0.85rem;
}
.pill-block span {
  margin: 0.25rem 0.25rem;
  padding: 0.25rem 0.85rem;
  border-radius: 15px;
}
.pill-block .pos {
  background: var(--color-pos-background);
}
.pill-block .sensetag {
  background: var(--color-tag-background);
}

.english-block {
  font-size: 1.75rem;
  margin: 2.5rem 0;
  line-height: 1rem;
}
.english-block span {
}
.english-block span:first-child {
  display: block;
  font-weight: bold;
  font-size: 1.75rem;
  line-height: 2rem;
  margin-bottom: 1rem;
  width: 100%;
}
.english-block :not(span:first-child) {
  font-size: 1rem;
}
.english-block span:nth-child(2):before {
  content: "(";
}
.english-block :not(span:first-child):after {
  content: ";";
}
.english-block span:last-child:not(:first-child):after {
  content: ")";
}

.example-block {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  counter-reset: list-number;
  gap: 0.75rem;
}
.example-block .example .japanese:before {
  counter-increment: list-number;
  content: counter(list-number) ". ";
}
.example-block .example {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}
.example-block .example .japanese {
}
.example-block .example .english {
  padding-left: 1.25rem;
  font-size: 0.85rem;
}
.highlight {
  font-weight: bold;
}
`,
		CardTemplates: []ankiconnect.CreateModelCardTemplate{
			{
				Name: "Recognition",
				Front: `<div class="kanji-block">
{{Kanji}}
</div>`,
				Back: `<div class="kanji-block">
{{furigana:Furigana}}
</div>
<div class="kana-block">
{{Kana}}
</div>
<div class="divider"></div>
<div class="english-block">
{{English}}
</div>
<div class="pill-block">
{{PoS}}
{{SenseTags}}
</div>
<div class="example-block">
{{furigana:Example}}
</div>
{{Audio}}`,
			},
			{
				Name: "Recall",
				Front: `<div class="english-block">
{{English}}
</div>`,
				Back: `<div class="english-block">
{{English}}
</div>
<div class="divider"></div>
<div class="kanji-block">
{{furigana:Furigana}}
</div>
<div class="kana-block">
{{Kana}}
</div>
<div class="pill-block">
{{PoS}}
{{SenseTags}}
</div>
<div class="example-block">
{{furigana:Example}}
{{Audio}}`,
			},
		},
	}
}
