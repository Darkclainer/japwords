package fxapp

import (
	"context"

	"go.uber.org/fx"

	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/multidict"
)

type MultiDictConfig struct {
	Workers int
}

func (c *MultiDictConfig) Equal(o any) bool {
	oc, ok := o.(*MultiDictConfig)
	if !ok {
		return false
	}
	return c.Workers == oc.Workers
}

type MultiDictIn struct {
	fx.In

	LC        fx.Lifecycle
	ConfigMgr *config.Manager
	LemmaDict multidict.LemmaDict
	PitchDict multidict.PitchDict
}

func NewMultidict(in MultiDictIn) (*multidict.MultiDict, error) {
	part, err := in.ConfigMgr.Register(config.ConsumerFunc(func(uc *config.UserConfig) (config.Part, error) {
		return &MultiDictConfig{
			Workers: uc.Dictionary.Workers,
		}, nil
	}))
	if err != nil {
		return nil, err
	}
	mdConfig := part.(*MultiDictConfig)

	dict, err := multidict.New(&multidict.Options{
		Workers:   mdConfig.Workers,
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
