package fxapp

import (
	"github.com/Darkclainer/japwords/pkg/cachedict"
	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/jisho"
	"github.com/Darkclainer/japwords/pkg/lemma"
)

func NewJisho(jishoClient jisho.BasicDict, uc *config.UserConfig) (*cachedict.CacheDict[[]*lemma.Lemma], error) {
	dict := jisho.New(jishoClient, uc.Dictionary.Jisho.URL)
	return cachedict.New[[]*lemma.Lemma](dict)
}
