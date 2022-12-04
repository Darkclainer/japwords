package fxapp

import (
	"japwords/pkg/basicdict"
	"japwords/pkg/jisho"
	"japwords/pkg/wadoku"
)

func NewBasicDict(fetcher basicdict.Fetcher) (jisho.BasicDict, wadoku.BasicDict) {
	bd := basicdict.New(fetcher)
	return bd, bd
}
