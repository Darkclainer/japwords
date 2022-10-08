package main

import (
	"go.uber.org/fx"

	"japwords/graphql/gqlresolver"
	"japwords/pkg/httpserver"
	"japwords/pkg/logger"
)

func NewApp(opts []fx.Option) (*fx.App, error) {
	opts = append(
		opts,
		fx.Provide(
			logger.New,

			httpserver.New,
			gqlresolver.New,
		),
		fx.Invoke(InvokeApp),
	)
	return fx.New(opts...), nil
}

func InvokeApp(
	server *httpserver.Server,
	resolver *gqlresolver.Resolver,
) {
	server.RegisterHandler("/query", resolver.Handler())
}
