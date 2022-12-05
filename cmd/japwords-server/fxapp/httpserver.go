package fxapp

import (
	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/httpserver"
)

func NewHttpServerConfig(uc *config.UserConfig) *httpserver.Config {
	return &httpserver.Config{
		Addr: uc.Addr,
	}
}
