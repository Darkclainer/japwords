package fxapp

import (
	"github.com/Darkclainer/japwords/pkg/basicdict"
	"github.com/Darkclainer/japwords/pkg/jisho"
	"github.com/Darkclainer/japwords/pkg/wadoku"
)

func NewBasicDict(fetcher basicdict.Fetcher) (jisho.BasicDict, wadoku.BasicDict) {
	bd := basicdict.New(fetcher)
	return bd, bd
}
