package main

import (
	"go.uber.org/fx"

	"japwords/pkg/config"
	"japwords/pkg/httpserver"
)

// ConvertConfig converts user config to fx.Option's to
// pass them to fx providers.
func ConvertConfig(uc *config.UserConfig) []fx.Option {
	var opts []fx.Option

	opts = append(opts, fx.Supply(
		&httpserver.Config{
			Addr: uc.Addr,
		},
	))

	return opts
}
