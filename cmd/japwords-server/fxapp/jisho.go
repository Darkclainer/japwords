package fxapp

import (
	"github.com/Darkclainer/japwords/pkg/cachedict"
	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/jisho"
	"github.com/Darkclainer/japwords/pkg/lemma"
)

type JishoConfig struct {
	URL string
}

func (c *JishoConfig) Equal(o any) bool {
	oc, ok := o.(*JishoConfig)
	if !ok {
		return false
	}
	return c.URL == oc.URL
}

func NewJisho(jishoClient jisho.BasicDict, configMgr *config.Manager) (*cachedict.CacheDict[[]*lemma.Lemma], error) {
	part, _, err := configMgr.Register(config.ConsumerFunc(func(uc *config.UserConfig) (config.Part, error) {
		return &JishoConfig{
			URL: uc.Dictionary.Jisho.URL,
		}, nil
	}))
	if err != nil {
		return nil, err
	}
	jishoConfig := part.(*JishoConfig)
	dict := jisho.New(jishoClient, jishoConfig.URL)
	return cachedict.New[[]*lemma.Lemma](dict)
}
