package fxapp

import (
	"japwords/pkg/config"
	"japwords/pkg/httpserver"
)

func NewHttpServerConfig(uc *config.UserConfig) *httpserver.Config {
	return &httpserver.Config{
		Addr: uc.Addr,
	}
}
