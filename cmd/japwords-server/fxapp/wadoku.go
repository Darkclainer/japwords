package fxapp

import (
	"japwords/pkg/cachedict"
	"japwords/pkg/config"
	"japwords/pkg/lemma"
	"japwords/pkg/wadoku"
)

func NewWadoku(wadokuClient wadoku.BasicDict, uc *config.UserConfig) (*cachedict.CacheDict[[]*lemma.PitchedLemma], error) {
	dict := wadoku.New(wadokuClient, uc.Dictionary.Wadoku.URL)
	return cachedict.New[[]*lemma.PitchedLemma](dict)
}
