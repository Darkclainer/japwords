package multidict

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/multierr"

	"github.com/Darkclainer/japwords/pkg/lemma"
	"github.com/Darkclainer/japwords/pkg/workerpool"
)

// MultiDict
type MultiDict struct {
	lemmaDict LemmaDict
	pitchDict PitchDict

	workerPool *workerpool.WorkerPool
}

type LemmaDict interface {
	Query(ctx context.Context, query string) ([]*lemma.Lemma, error)
}

type PitchDict interface {
	Query(ctx context.Context, query string) ([]*lemma.PitchedLemma, error)
}

func New(opts *Options) (*MultiDict, error) {
	workers := opts.Workers
	if workers == 0 {
		workers = 4
	}
	wp, err := workerpool.New(workers)
	if err != nil {
		return nil, fmt.Errorf("new multidict failed: %w", err)
	}

	return &MultiDict{
		lemmaDict:  opts.LemmaDict,
		pitchDict:  opts.PitchDict,
		workerPool: wp,
	}, nil
}

func (m *MultiDict) Init() {
	m.workerPool.Init()
}

func (m *MultiDict) Close() {
	m.workerPool.Stop()
}

func (m *MultiDict) Query(ctx context.Context, query string) ([]*lemma.Lemma, error) {
	ctx, cancel := m.defaultContext(ctx)
	defer cancel()

	lemmasChan, err := m.queryAsync(ctx, query)
	if err != nil {
		return nil, err
	}
	pitchChan, err := m.queryPitchAsync(ctx, query)
	if err != nil {
		return nil, err
	}
	var (
		lemmas      []*lemma.Lemma
		pitches     []*lemma.PitchedLemma
		combinedErr error
	)
	// we want collect result from both dictionary if possible,
	// but pitch dictionary more or less optional.
	for lemmasChan != nil || pitchChan != nil {
		select {
		case r := <-lemmasChan:
			if r.Err != nil {
				combinedErr = multierr.Append(
					combinedErr,
					fmt.Errorf("lemma dict request failed: %w", r.Err),
				)
			}
			lemmasChan = nil
			lemmas = r.Value
			// if we didn't get any lemmas then no need to wait pitches
			if len(lemmas) == 0 {
				return nil, combinedErr
			}
		case r := <-pitchChan:
			if r.Err != nil {
				combinedErr = multierr.Append(
					combinedErr,
					fmt.Errorf("pitch dict request failed: %w", r.Err),
				)
			}
			pitchChan = nil
			pitches = r.Value
		case <-ctx.Done():
			return lemmas, ctx.Err()
		}
	}
	lemma.Enrich(lemmas, pitches)
	return lemmas, combinedErr
}

func (m *MultiDict) QueryPitch(ctx context.Context, slug string, hiragana string) ([]lemma.Pitch, error) {
	ctx, cancel := m.defaultContext(ctx)
	defer cancel()
	pitchedLemmasChan, err := m.queryPitchAsync(ctx, slug)
	if err != nil {
		return nil, err
	}
	var pitchedLemmas []*lemma.PitchedLemma
	var pitchedErr error
	select {
	case r := <-pitchedLemmasChan:
		pitchedErr = r.Err
		pitchedLemmas = r.Value
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	for _, reading := range pitchedLemmas {
		if reading.Slug == slug && reading.Hiragana == hiragana {
			return reading.Pitches, pitchedErr
		}
	}
	return nil, pitchedErr
}

func (m *MultiDict) queryAsync(ctx context.Context, query string) (<-chan result[[]*lemma.Lemma], error) {
	// we need boffered channel to not block workerpool
	lemmaResult := make(chan result[[]*lemma.Lemma], 1)
	fn := func(fctx context.Context) {
		lemmas, err := m.lemmaDict.Query(fctx, query)
		lemmaResult <- newResult(lemmas, err)
	}
	err := m.workerPool.Add(ctx, fn)
	if err != nil {
		return nil, err
	}
	return lemmaResult, nil
}

func (m *MultiDict) queryPitchAsync(ctx context.Context, query string) (<-chan result[[]*lemma.PitchedLemma], error) {
	// we need boffered channel to not block workerpool
	pitchResult := make(chan result[[]*lemma.PitchedLemma], 1)
	fn := func(fctx context.Context) {
		pitches, err := m.pitchDict.Query(fctx, query)
		pitchResult <- newResult(pitches, err)
	}
	err := m.workerPool.Add(ctx, fn)
	if err != nil {
		return nil, err
	}
	return pitchResult, nil
}

func (m *MultiDict) defaultContext(ctx context.Context) (context.Context, context.CancelFunc) {
	// some default timeout for safety
	ctx, cancel := context.WithTimeout(ctx, time.Second*45)
	return ctx, cancel
}

type result[T any] struct {
	Value T
	Err   error
}

func newResult[T any](v T, err error) result[T] {
	return result[T]{
		Value: v,
		Err:   err,
	}
}
