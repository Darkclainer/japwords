package fxapp

import (
	"github.com/Darkclainer/japwords/pkg/cachedict"
	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/lemma"
	"github.com/Darkclainer/japwords/pkg/wadoku"
)

func NewWadoku(wadokuClient wadoku.BasicDict, uc *config.UserConfig) (*cachedict.CacheDict[[]*lemma.PitchedLemma], error) {
	dict := wadoku.New(wadokuClient, uc.Dictionary.Wadoku.URL)
	return cachedict.New[[]*lemma.PitchedLemma](dict)
}
