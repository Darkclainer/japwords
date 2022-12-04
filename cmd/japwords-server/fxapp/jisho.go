package fxapp

import (
	"japwords/pkg/cachedict"
	"japwords/pkg/config"
	"japwords/pkg/jisho"
	"japwords/pkg/lemma"
)

func NewJisho(jishoClient jisho.BasicDict, uc *config.UserConfig) (*cachedict.CacheDict[[]*lemma.Lemma], error) {
	dict := jisho.New(jishoClient, uc.Dictionary.Jisho.URL)
	return cachedict.New[[]*lemma.Lemma](dict)
}
