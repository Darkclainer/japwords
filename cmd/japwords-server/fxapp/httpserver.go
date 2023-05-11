package fxapp

import (
	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/httpserver"
)

func NewHttpServerConfig(configMgr *config.Manager) *httpserver.Config {
	part, err := configMgr.Register(config.ConsumerFunc(func(uc *config.UserConfig) (config.Part, error) {
		return &httpserver.Config{
			Addr: uc.Addr,
		}, nil
	}))
	if err != nil {
		return nil
	}
	conf := part.(*httpserver.Config)
	return conf
}
