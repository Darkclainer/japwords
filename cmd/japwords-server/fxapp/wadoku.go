package fxapp

import (
	"github.com/Darkclainer/japwords/pkg/cachedict"
	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/lemma"
	"github.com/Darkclainer/japwords/pkg/wadoku"
)

type WadokuConfig struct {
	URL string
}

func (c *WadokuConfig) Equal(o any) bool {
	oc, ok := o.(*WadokuConfig)
	if !ok {
		return false
	}
	return c.URL == oc.URL
}

func NewWadoku(wadokuClient wadoku.BasicDict, configMgr *config.Manager) (*cachedict.CacheDict[[]*lemma.PitchedLemma], error) {
	part, err := configMgr.Register(config.ConsumerFunc(func(uc *config.UserConfig) (config.Part, error) {
		return &WadokuConfig{
			URL: uc.Dictionary.Wadoku.URL,
		}, nil
	}))
	if err != nil {
		return nil, err
	}
	wadokuConfig := part.(*WadokuConfig)
	dict := wadoku.New(wadokuClient, wadokuConfig.URL)
	return cachedict.New[[]*lemma.PitchedLemma](dict)
}
