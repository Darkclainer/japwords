package multidict

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"japwords/pkg/lemma"
)

type LemmaDictTest func(ctx context.Context, query string) ([]*lemma.Lemma, error)

func (f LemmaDictTest) Query(ctx context.Context, query string) ([]*lemma.Lemma, error) {
	return f(ctx, query)
}

type PitchDictTest func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error)

func (f PitchDictTest) Query(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
	return f(ctx, query)
}

func Test_Multidict_New(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	multidict, err := New(&Options{
		Workers: 0,
		LemmaDict: LemmaDictTest(func(ctx context.Context, query string) ([]*lemma.Lemma, error) {
			return getLemmasTest("query"), nil
		}),
		PitchDict: PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
			return getPitchedLemmasTest("query"), nil
		}),
	})
	require.NoError(t, err)
	multidict.Init()
	defer multidict.Close()
	lemmas, err := multidict.Query(ctx, "query")
	assert.NoError(t, err)
	assert.Equal(t, getResultLemmasTest("query"), lemmas)
}

func Test_Multidict_Query(t *testing.T) {
	const query = "query"
	testCases := []struct {
		Name string
		// Because multidict use single workerpool, some test can not pass
		// in principle with single worker. This can happen if first task to run
		// will wait second (that will not run ever, because there is no free worker).
		// In real environment it's hardly a problem, because request will be timeouted
		// eventualy.
		MinWorkers   int
		InitHandlers func(t *testing.T, ctx context.Context, cancel context.CancelFunc) (LemmaDictTest, PitchDictTest)
		ErrorAssert  assert.ErrorAssertionFunc
		Expected     func(t *testing.T, result []*lemma.Lemma)
	}{
		{
			Name: "OK",
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) (LemmaDictTest, PitchDictTest) {
				lemmaDict := LemmaDictTest(func(ctx context.Context, query string) ([]*lemma.Lemma, error) {
					return getLemmasTest(query), nil
				})
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					return getPitchedLemmasTest(query), nil
				})
				return lemmaDict, pitchDict
			},
			ErrorAssert: assert.NoError,
			Expected: func(t *testing.T, result []*lemma.Lemma) {
				assert.Equal(t, getResultLemmasTest(query), result)
			},
		},
		{
			Name: "OK lemma first",
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) (LemmaDictTest, PitchDictTest) {
				wait := make(chan struct{})
				lemmaDict := LemmaDictTest(func(ctx context.Context, query string) ([]*lemma.Lemma, error) {
					defer close(wait)
					return getLemmasTest(query), nil
				})
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					<-wait
					return getPitchedLemmasTest(query), nil
				})
				return lemmaDict, pitchDict
			},
			ErrorAssert: assert.NoError,
			Expected: func(t *testing.T, result []*lemma.Lemma) {
				assert.Equal(t, getResultLemmasTest(query), result)
			},
		},
		{
			Name:       "OK pitch first",
			MinWorkers: 2,
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) (LemmaDictTest, PitchDictTest) {
				wait := make(chan struct{})
				lemmaDict := LemmaDictTest(func(ctx context.Context, query string) ([]*lemma.Lemma, error) {
					<-wait
					return getLemmasTest(query), nil
				})
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					defer close(wait)
					return getPitchedLemmasTest(query), nil
				})
				return lemmaDict, pitchDict
			},
			ErrorAssert: assert.NoError,
			Expected: func(t *testing.T, result []*lemma.Lemma) {
				assert.Equal(t, getResultLemmasTest(query), result)
			},
		},
		{
			Name: "lemma fail first",
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) (LemmaDictTest, PitchDictTest) {
				lemmaDict := LemmaDictTest(func(ctx context.Context, query string) ([]*lemma.Lemma, error) {
					return nil, errors.New("lemma error")
				})
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					<-ctx.Done()
					return getPitchedLemmasTest(query), nil
				})
				return lemmaDict, pitchDict
			},
			ErrorAssert: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "lemma error")
			},
			Expected: func(t *testing.T, result []*lemma.Lemma) {
				assert.Nil(t, result)
			},
		},
		{
			Name:       "lemma fail second",
			MinWorkers: 2,
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) (LemmaDictTest, PitchDictTest) {
				wait := make(chan struct{})
				lemmaDict := LemmaDictTest(func(ctx context.Context, query string) ([]*lemma.Lemma, error) {
					<-wait
					return nil, errors.New("lemma error")
				})
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					defer close(wait)
					return getPitchedLemmasTest(query), nil
				})
				return lemmaDict, pitchDict
			},
			ErrorAssert: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "lemma error")
			},
			Expected: func(t *testing.T, result []*lemma.Lemma) {
				assert.Nil(t, result)
			},
		},
		{
			Name: "pitch fail second",
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) (LemmaDictTest, PitchDictTest) {
				wait := make(chan struct{})
				lemmaDict := LemmaDictTest(func(ctx context.Context, query string) ([]*lemma.Lemma, error) {
					defer close(wait)
					return getLemmasTest(query), nil
				})
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					<-wait
					return nil, errors.New("pitch error")
				})
				return lemmaDict, pitchDict
			},
			ErrorAssert: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "pitch error")
			},
			Expected: func(t *testing.T, result []*lemma.Lemma) {
				assert.Equal(t, getLemmasTest(query), result)
			},
		},
		{
			Name:       "pitch fail first",
			MinWorkers: 2,
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) (LemmaDictTest, PitchDictTest) {
				wait := make(chan struct{})
				lemmaDict := LemmaDictTest(func(ctx context.Context, query string) ([]*lemma.Lemma, error) {
					<-wait
					return getLemmasTest(query), nil
				})
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					defer close(wait)
					return nil, errors.New("pitch error")
				})
				return lemmaDict, pitchDict
			},
			ErrorAssert: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "pitch error")
			},
			Expected: func(t *testing.T, result []*lemma.Lemma) {
				assert.Equal(t, getLemmasTest(query), result)
			},
		},
		{
			Name:       "lemma partial fail first",
			MinWorkers: 2,
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) (LemmaDictTest, PitchDictTest) {
				wait := make(chan struct{})
				lemmaDict := LemmaDictTest(func(ctx context.Context, query string) ([]*lemma.Lemma, error) {
					defer close(wait)
					return getLemmasTest(query), errors.New("lemma error")
				})
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					<-wait
					return getPitchedLemmasTest(query), nil
				})
				return lemmaDict, pitchDict
			},
			ErrorAssert: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "lemma error")
			},
			Expected: func(t *testing.T, result []*lemma.Lemma) {
				assert.Equal(t, getResultLemmasTest(query), result)
			},
		},
		{
			Name:       "context cancelled",
			MinWorkers: 1,
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) (LemmaDictTest, PitchDictTest) {
				wait := make(chan struct{})
				lemmaDict := LemmaDictTest(func(fctx context.Context, query string) ([]*lemma.Lemma, error) {
					cancel()
					defer close(wait)
					return getLemmasTest(query), fctx.Err()
				})
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					<-wait
					return nil, nil
				})
				return lemmaDict, pitchDict
			},
			ErrorAssert: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(tt, err, "context canceled")
			},
			Expected: func(t *testing.T, result []*lemma.Lemma) {
				// we either would or would not get result, because context
				// cancel and sending result is concurrent
				if len(result) == 0 {
					return
				}
				assert.Equal(t, getLemmasTest(query), result)
			},
		},
	}
	for i := range testCases {
		workers := testCases[i].MinWorkers
		if workers == 0 {
			workers = 1
		}
		for ; workers < 3; workers++ {
			workers := workers
			tc := testCases[i]
			t.Run(fmt.Sprintf("name=%s,workers=%d", tc.Name, workers), func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				lemmaDict, pitchDict := tc.InitHandlers(t, ctx, cancel)
				multidict, err := New(&Options{
					Workers:   workers,
					LemmaDict: lemmaDict,
					PitchDict: pitchDict,
				})
				require.NoError(t, err)
				multidict.Init()
				defer multidict.Close()
				lemmas, err := multidict.Query(ctx, query)
				tc.ErrorAssert(t, err)
				tc.Expected(t, lemmas)
			})
		}
	}
}

func Test_Multidict_Stop_Query(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	multidict, err := New(&Options{
		Workers:   1,
		LemmaDict: nil,
		PitchDict: nil,
	})
	require.NoError(t, err)
	multidict.Init()
	multidict.Close()
	lemmas, err := multidict.Query(ctx, "hello")
	assert.Error(t, err)
	assert.Nil(t, lemmas)
}

func Test_Multidict_Query_Stop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	multidict, err := New(&Options{
		Workers:   1,
		LemmaDict: nil,
		PitchDict: nil,
	})
	require.NoError(t, err)
	multidict.lemmaDict = LemmaDictTest(func(ctx context.Context, query string) ([]*lemma.Lemma, error) {
		multidict.Close()
		<-ctx.Done()
		return nil, ctx.Err()
	})
	multidict.pitchDict = PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
		return nil, nil
	})
	multidict.Init()

	lemmas, err := multidict.Query(ctx, "hello")
	assert.Error(t, err)
	assert.Nil(t, lemmas)
}

func Test_Multidict_QueryPitch(t *testing.T) {
	testCases := []struct {
		Name         string
		Slug         string
		Hiragana     string
		InitHandlers func(t *testing.T, ctx context.Context, cancel context.CancelFunc) PitchDictTest
		ErrorAssert  assert.ErrorAssertionFunc
		Expected     []lemma.Pitch
	}{
		{
			Name:     "first result",
			Slug:     "hello",
			Hiragana: "world",
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) PitchDictTest {
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					return []*lemma.PitchedLemma{
						{
							Slug:     "hello",
							Hiragana: "world",
							Pitches: []lemma.Pitch{
								{
									Position: 9,
								},
							},
						},
						{
							Slug:     "hi",
							Hiragana: "world",
							Pitches: []lemma.Pitch{
								{
									Position: 4,
								},
							},
						},
					}, nil
				})
				return pitchDict
			},
			ErrorAssert: assert.NoError,
			Expected: []lemma.Pitch{
				{
					Position: 9,
				},
			},
		},
		{
			Name:     "second result",
			Slug:     "hi",
			Hiragana: "world",
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) PitchDictTest {
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					return []*lemma.PitchedLemma{
						{
							Slug:     "hello",
							Hiragana: "world",
							Pitches: []lemma.Pitch{
								{
									Position: 9,
								},
							},
						},
						{
							Slug:     "hi",
							Hiragana: "world",
							Pitches: []lemma.Pitch{
								{
									Position: 4,
								},
							},
						},
					}, nil
				})
				return pitchDict
			},
			ErrorAssert: assert.NoError,
			Expected: []lemma.Pitch{
				{
					Position: 4,
				},
			},
		},
		{
			Name:     "second result and error",
			Slug:     "hi",
			Hiragana: "world",
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) PitchDictTest {
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					return []*lemma.PitchedLemma{
						{
							Slug:     "hello",
							Hiragana: "world",
							Pitches: []lemma.Pitch{
								{
									Position: 9,
								},
							},
						},
						{
							Slug:     "hi",
							Hiragana: "world",
							Pitches: []lemma.Pitch{
								{
									Position: 4,
								},
							},
						},
					}, errors.New("pitch error")
				})
				return pitchDict
			},
			ErrorAssert: assert.Error,
			Expected: []lemma.Pitch{
				{
					Position: 4,
				},
			},
		},
		{
			Name:     "only error",
			Slug:     "hi",
			Hiragana: "world",
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) PitchDictTest {
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					return nil, errors.New("pitch error")
				})
				return pitchDict
			},
			ErrorAssert: assert.Error,
			Expected:    nil,
		},
		{
			Name:     "cancelled context",
			Slug:     "hi",
			Hiragana: "world",
			InitHandlers: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) PitchDictTest {
				pitchDict := PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
					cancel()
					return nil, errors.New("pitch error")
				})
				return pitchDict
			},
			ErrorAssert: assert.Error,
			Expected:    nil,
		},
	}
	for i := range testCases {
		for workers := 1; workers < 3; workers++ {
			workers := workers
			tc := testCases[i]
			t.Run(fmt.Sprintf("name=%s,workers=%d", tc.Name, workers), func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				pitchDict := tc.InitHandlers(t, ctx, cancel)
				multidict, err := New(&Options{
					Workers:   workers,
					PitchDict: pitchDict,
				})
				require.NoError(t, err)
				multidict.Init()
				defer multidict.Close()
				pitches, err := multidict.QueryPitch(ctx, tc.Slug, tc.Hiragana)
				tc.ErrorAssert(t, err)
				assert.Equal(t, tc.Expected, pitches)
			})
		}
	}
}

func Test_Multidict_Stop_QueryPitch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	multidict, err := New(&Options{
		Workers:   1,
		PitchDict: nil,
	})
	require.NoError(t, err)
	multidict.Init()
	multidict.Close()
	lemmas, err := multidict.QueryPitch(ctx, "hello", "world")
	assert.Error(t, err)
	assert.Nil(t, lemmas)
}

func Test_Multidict_QueryPitch_Stop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	multidict, err := New(&Options{
		Workers: 1,
	})
	require.NoError(t, err)
	wait := make(chan struct{})
	multidict.pitchDict = PitchDictTest(func(ctx context.Context, query string) ([]*lemma.PitchedLemma, error) {
		close(wait)
		<-ctx.Done()
		return nil, ctx.Err()
	})
	multidict.Init()
	go func() {
		<-wait
		multidict.Close()
	}()
	lemmas, err := multidict.QueryPitch(ctx, "hello", "world")
	assert.Error(t, err)
	assert.Nil(t, lemmas)
}

func getLemmasTest(query string) []*lemma.Lemma {
	return []*lemma.Lemma{
		{
			Slug: lemma.Word{
				Word:     query,
				Hiragana: query,
			},
			Tags: []string{"just test"},
		},
	}
}

func getPitchedLemmasTest(query string) []*lemma.PitchedLemma {
	return []*lemma.PitchedLemma{
		{
			Slug:     query,
			Hiragana: query,
			Pitches: []lemma.Pitch{
				{
					Position: len(query),
					IsHigh:   true,
				},
			},
		},
	}
}

func getResultLemmasTest(query string) []*lemma.Lemma {
	return []*lemma.Lemma{
		{
			Slug: lemma.Word{
				Word:     query,
				Hiragana: query,
				Pitches: []lemma.Pitch{
					{
						Position: len(query),
						IsHigh:   true,
					},
				},
			},
			Tags: []string{"just test"},
		},
	}
}
