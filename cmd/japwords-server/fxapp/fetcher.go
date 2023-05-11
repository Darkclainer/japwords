package fxapp

import (
	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/fetcher"
)

func NewFetcher(configMgr *config.Manager) (*fetcher.Fetcher, error) {
	part, err := configMgr.Register(config.ConsumerFunc(func(uc *config.UserConfig) (config.Part, error) {
		var conf fetcher.Config
		conf.Headers = uc.Dictionary.Headers
		if conf.Headers == nil {
			conf.Headers = make(map[string]string)
		}
		conf.Headers["User-Agent"] = uc.Dictionary.UserAgent
		return &conf, nil
	}))
	if err != nil {
		return nil, err
	}
	fetcherConf := part.(*fetcher.Config)
	fetcherClient, err := fetcher.New(fetcherConf)
	return fetcherClient, err
}
