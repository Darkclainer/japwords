package fxapp

import (
	"context"

	"go.uber.org/fx"

	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/multidict"
)

type MultiDictIn struct {
	fx.In

	LC         fx.Lifecycle
	UserConfig *config.UserConfig
	LemmaDict  multidict.LemmaDict
	PitchDict  multidict.PitchDict
}

func NewMultidict(in MultiDictIn) (*multidict.MultiDict, error) {
	uc := in.UserConfig

	dict, err := multidict.New(&multidict.Options{
		Workers:   uc.Dictionary.Workers,
		LemmaDict: in.LemmaDict,
		PitchDict: in.PitchDict,
	})
	if err != nil {
		return nil, err
	}
	in.LC.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			dict.Init()
			return nil
		},
		OnStop: func(_ context.Context) error {
			dict.Close()
			return nil
		},
	})
	return dict, nil
}
