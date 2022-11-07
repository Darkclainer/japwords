package wadoku

import (
	"strings"
	"testing"

	"japwords/pkg/htmltest"

	"github.com/stretchr/testify/assert"
)

func Test_parseHTML(t *testing.T) {
	testCases := []struct {
		Name        string
		HTML        string
		Expected    []*Lemma
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "simple",
			HTML: `
	<html><body>
		<section id="content"></section>
	</body></html>`,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "no content section",
			HTML: `
	<html><body>
		<table id="resulttable"></table>
	</body></html>`,
			ErrorAssert: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "content section")
			},
		},
		{
			Name: "one lemma",
			HTML: `
	<html><body><section id="content">
		<table id="resulttable"><tbody>
			<tr class="resultline"> 
				<td class="resultdetail">
					<div class="japanese">
						<span class="orth">hello</span>
					</div>
					<div class="reading">
						<span class="pron accent">
							<span class="t">
								world
							</span>
						</span>
					</div>
				</td>
			</tr>
		</tbody></table>
	</section></body></html>`,
			Expected: []*Lemma{
				{
					Slug: "hello",
					Reading: Reading{
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 5,
								IsHigh:   true,
							},
						},
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			lemmas, err := parseHTML(strings.NewReader(tc.HTML))
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, lemmas)
		})
	}
}

func Test_parseContentSection(t *testing.T) {
	testCases := []struct {
		Name        string
		HTML        string
		Expected    []*Lemma
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "empty",
			HTML: `
		<div id="root"> 
		</div>`,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "single",
			HTML: `
	<div id="root"><table id="resulttable"><tbody>
		<tr class="resultline"> 
			<td class="resultdetail">
				<div class="japanese">
					<span class="orth">hello</span>
				</div>
				<div class="reading">
					<span class="pron accent">
						<span class="t">
							world
						</span>
					</span>
				</div>
			</td>
		</tr>
	</tbody></table></div>`,
			Expected: []*Lemma{
				{
					Slug: "hello",
					Reading: Reading{
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 5,
								IsHigh:   true,
							},
						},
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "two resultlines",
			HTML: `
	<div id="root"><table id="resulttable"><tbody>
		<tr class="resultline"> 
			<td class="resultdetail">
				<div class="japanese">
					<span class="orth">hello</span>
				</div>
				<div class="reading">
					<span class="pron accent">
						<span class="t">
							world
						</span>
					</span>
				</div>
			</td>
		</tr>
		<tr class="resultline"> 
			<td class="resultdetail">
				<div class="japanese">
					<span class="orth">greetings</span>
				</div>
				<div class="reading">
					<span class="pron accent">
						<span class="t">
							world
						</span>
					</span>
				</div>
			</td>
		</tr>
	</tbody></table></div>`,
			Expected: []*Lemma{
				{
					Slug: "hello",
					Reading: Reading{
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 5,
								IsHigh:   true,
							},
						},
					},
				},
				{
					Slug: "greetings",
					Reading: Reading{
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 5,
								IsHigh:   true,
							},
						},
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "result and error seperate",
			HTML: `
	<div id="root"><table id="resulttable"><tbody>
		<tr class="resultline"> 
			<td class="resultdetail">
				<div class="japanese">
					<span class="orth">hello</span>
				</div>
				<div class="reading">
					<span class="pron accent">
						<span class="t">
							world
						</span>
					</span>
				</div>
			</td>
		</tr>
		<tr class="resultline"> 
			<td class="resultdetail">
				<div class="japanese">
					<span class="orth"></span>
				</div>
				<div class="reading">
					<span class="pron accent">
						<span class="t">
							world
						</span>
					</span>
				</div>
			</td>
		</tr>
	</tbody></table></div>`,
			Expected: []*Lemma{
				{
					Slug: "hello",
					Reading: Reading{
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 5,
								IsHigh:   true,
							},
						},
					},
				},
			},
			ErrorAssert: func(tt assert.TestingT, err error, _ ...interface{}) bool {
				var batchErr *LemmaBatchError
				if !assert.ErrorAs(tt, err, &batchErr) {
					return false
				}
				if !assert.Len(tt, batchErr.Errs, 1) {
					return false
				}
				lemmaErr := batchErr.Errs[0]
				if !assert.Equal(tt, 1, lemmaErr.ID) {
					return false
				}
				return assert.ErrorContains(tt, lemmaErr, "no japanese slug")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			root := htmltest.MustRootSelection(t, tc.HTML)
			result, err := parseContentSection(root)
			tc.ErrorAssert(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func Test_parseRowResult(t *testing.T) {
	testCases := []struct {
		Name        string
		HTML        string
		Expected    []*Lemma
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "simple",
			HTML: `
	<table><tbody><tr id="root"> 
		<td class="resultdetail">
			<div class="japanese">
				<span class="orth">hello</span>
			</div>
			<div class="reading">
				<span class="pron accent">
					<span class="t">
						world
					</span>
				</span>
			</div>
		</td>
	</tr></tbody></table>`,
			Expected: []*Lemma{
				{
					Slug: "hello",
					Reading: Reading{
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 5,
								IsHigh:   true,
							},
						},
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "two slug",
			HTML: `
	<table><tbody><tr id="root"> 
		<td class="resultdetail">
			<div class="japanese">
				<span class="orth">
					hello
					<span class="divider">;</span>
					nothello
				</span>
			</div>
			<div class="reading">
				<span class="pron accent">
					<span class="t">
						world
					</span>
				</span>
			</div>
		</td>
	</tr></tbody></table>`,
			Expected: []*Lemma{
				{
					Slug: "hello",
					Reading: Reading{
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 5,
								IsHigh:   true,
							},
						},
					},
				},
				{
					Slug: "nothello",
					Reading: Reading{
						Hiragana: "world",
						Pitches: []Pitch{
							{
								Position: 5,
								IsHigh:   true,
							},
						},
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "filtered slug",
			HTML: `
	<table><tbody><tr id="root"> 
		<td class="resultdetail">
			<div class="japanese">
				<span class="orth">…hello</span>
			</div>
			<div class="reading">
				<span class="pron accent">
					<span class="t">
						world
					</span>
				</span>
			</div>
		</td>
	</tr></tbody></table>`,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "filtered reading",
			HTML: `
	<table><tbody><tr id="root"> 
		<td class="resultdetail">
			<div class="japanese">
				<span class="orth">hello</span>
			</div>
			<div class="reading">
			</div>
		</td>
	</tr></tbody></table>`,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "no slug error",
			HTML: `
	<table><tbody><tr id="root"> 
		<td class="resultdetail">
			<div class="japanese">
			</div>
			<div class="reading">
			</div>
		</td>
	</tr></tbody></table>`,
			ErrorAssert: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "no japanese slug")
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			root := htmltest.MustRootSelection(t, tc.HTML)
			result, err := parseRowResult(root)
			tc.ErrorAssert(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func Test_parseJapanese(t *testing.T) {
	testCases := []struct {
		Name        string
		HTML        string
		Expected    []string
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "simple",
			HTML: `
		<div id="root"> 
			<span class="orth">hello</span>
		</div>`,
			Expected:    []string{"hello"},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "one divider",
			HTML: `
		<div id="root"> 
			<span class="orth">  
				hello
				<span class="divider">;</span>
				world 
			</span>
		</div>`,
			Expected:    []string{"hello", "world"},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "extra divider",
			HTML: `
		<div id="root"> 
			<span class="orth">  
				hello
				<span class="divider">;</span>
				world 
				<span class="divider">;</span>
			</span>
		</div>`,
			Expected:    []string{"hello", "world"},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "text inside span",
			HTML: `
		<div id="root"> 
			<span class="orth">  
				hello
				<span class="njok">world</span>
			</span>
		</div>`,
			Expected:    []string{"helloworld"},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "prefix filtered",
			HTML: `
		<div id="root"> 
			<span class="orth">  
				hello
				<span class="divider">;</span>
				…world 
			</span>
		</div>`,
			Expected:    []string{"hello"},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "suffix filtered",
			HTML: `
		<div id="root"> 
			<span class="orth">  
				hello
				<span class="divider">;</span>
				world…
			</span>
		</div>`,
			Expected:    []string{"hello"},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "everything filtered",
			HTML: `
		<div id="root"> 
			<span class="orth">  
				…hello
			</span>
		</div>`,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "nothing found",
			HTML: `
		<div id="root"> 
		</div>`,
			ErrorAssert: assert.Error,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			root := htmltest.MustRootSelection(t, tc.HTML)
			variants, err := parseJapanese(root)
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, variants)
		})
	}
}

func Test_parseReading(t *testing.T) {
	testCases := []struct {
		Name        string
		HTML        string
		Expected    *Reading
		ErrorAssert assert.ErrorAssertionFunc
	}{
		{
			Name: "contains no reading",
			HTML: `
		<div id="root"> 
			<span class="pron accent">
			</span>
		</div>`,
			ErrorAssert: assert.NoError,
		},
		{
			Name: "contains no reading found",
			HTML: `
		<div id="root"> 
			<span class="pron accent">
				<span class="t">
				</span>
			</span>
		</div>`,
			ErrorAssert: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "no reading")
			},
		},
		{
			Name: "constant up",
			HTML: `
		<div id="root"> 
			<span class="pron accent">
				<span class="t">
					hello
				</span>
			</span>
		</div>`,
			Expected: &Reading{
				Hiragana: "hello",
				Pitches: []Pitch{
					{
						Position: 5,
						IsHigh:   true,
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "constant down",
			HTML: `
		<div id="root"> 
			<span class="pron accent">
				<span class="b">
					worlds
				</span>
			</span>
		</div>`,
			Expected: &Reading{
				Hiragana: "worlds",
				Pitches: []Pitch{
					{
						Position: 6,
						IsHigh:   false,
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "up down",
			HTML: `
		<div id="root"> 
			<span class="pron accent">
				<span class="t">
					hel
				</span>
				<span class="b">
					lo
				</span>
			</span>
		</div>`,
			Expected: &Reading{
				Hiragana: "hello",
				Pitches: []Pitch{
					{
						Position: 3,
						IsHigh:   true,
					},
					{
						Position: 5,
						IsHigh:   false,
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "down up",
			HTML: `
		<div id="root"> 
			<span class="pron accent">
				<span class="b">
					hel
				</span>
				<span class="t">
					lo
				</span>
			</span>
		</div>`,
			Expected: &Reading{
				Hiragana: "hello",
				Pitches: []Pitch{
					{
						Position: 3,
						IsHigh:   false,
					},
					{
						Position: 5,
						IsHigh:   true,
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "down up down (end)",
			HTML: `
		<div id="root"> 
			<span class="pron accent">
				<span class="b">
					hel
				</span>
				<span class="t r">
					lo
				</span>
			</span>
		</div>`,
			Expected: &Reading{
				Hiragana: "hello",
				Pitches: []Pitch{
					{
						Position: 3,
						IsHigh:   false,
					},
					{
						Position: 5,
						IsHigh:   true,
					},
					{
						Position: 5,
						IsHigh:   false,
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{
			Name: "up down up (end)",
			HTML: `
		<div id="root"> 
			<span class="pron accent">
				<span class="t">
					hel
				</span>
				<span class="b r">
					lo
				</span>
			</span>
		</div>`,
			Expected: &Reading{
				Hiragana: "hello",
				Pitches: []Pitch{
					{
						Position: 3,
						IsHigh:   true,
					},
					{
						Position: 5,
						IsHigh:   false,
					},
					{
						Position: 5,
						IsHigh:   true,
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
		{ // real example
			Name: "inu neko",
			HTML: `
		<div id="root"> 
			<span class="pron accent" data-accent-id="1">
				<span class="b r">い</span>
				<span class="t r">ぬ･ね</span>
				<span class="b">こ</span>
			</span>
		</div>`,
			Expected: &Reading{
				Hiragana: "いぬねこ",
				// len(い) = len(ぬ) = len(ね) = len(こ) = 3
				Pitches: []Pitch{
					{
						Position: 3,
						IsHigh:   false,
					},
					{
						Position: 9,
						IsHigh:   true,
					},
					{
						Position: 12,
						IsHigh:   false,
					},
				},
			},
			ErrorAssert: assert.NoError,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			root := htmltest.MustRootSelection(t, tc.HTML)
			variants, err := parseReading(root)
			tc.ErrorAssert(t, err)
			assert.Equal(t, tc.Expected, variants)
		})
	}
}

func Test_extractReading(t *testing.T) {
	testCases := []struct {
		Name     string
		HTML     string
		Expected string
	}{
		{
			Name: "empty",
			HTML: `
		<div id="root"> 
		</div>`,
			Expected: "",
		},
		{
			Name: "simple",
			HTML: `
		<div id="root"> 
			hello
		</div>`,
			Expected: "hello",
		},
		{
			Name: "simple japanese",
			HTML: `
		<div id="root"> 
			三人	
		</div>`,
			Expected: "三人",
		},
		{
			Name: "dot",
			HTML: `
		<div id="root"> 
			hel･lo
		</div>`,
			Expected: "hello",
		},
		{
			Name: "divider",
			HTML: `
		<div id="root"> 
			hel
			<span class="divider">|</span>
			lo
		</div>`,
			Expected: "hello",
		},
		{
			Name: "dot divider",
			HTML: `
		<div id="root"> 
			h･el
			<span class="divider">|</span>
			lo
		</div>`,
			Expected: "hello",
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			root := htmltest.MustRootSelection(t, tc.HTML)
			variants := extractReading(root)
			assert.Equal(t, tc.Expected, variants)
		})
	}
}
