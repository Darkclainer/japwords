package fxapp

import (
	"go.uber.org/fx"

	"github.com/Darkclainer/japwords/graphql/gqlresolver"
	"github.com/Darkclainer/japwords/pkg/basicdict"
	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/httpserver"
	"github.com/Darkclainer/japwords/pkg/logger"
	"github.com/Darkclainer/japwords/pkg/multidict"
	"github.com/Darkclainer/japwords/ui"
)

func NewApp(configMgr *config.Manager) (*fx.App, error) {
	opts := []fx.Option{
		// util staff
		fx.Supply(configMgr),
		fx.Provide(
			logger.New,
		),
		// dictionary things
		fx.Provide(
			fx.Annotate(
				NewFetcher,
				fx.As(new(basicdict.Fetcher)),
			),
		),
		fx.Provide(
			NewBasicDict,
		),
		fx.Provide(
			fx.Annotate(
				NewJisho,
				fx.As(new(multidict.LemmaDict)),
			),
			fx.Annotate(
				NewWadoku,
				fx.As(new(multidict.PitchDict)),
			),
		),
		fx.Provide(NewMultidict),
		fx.Provide(NewAnki),
		// http/graphql staff
		fx.Provide(
			NewHttpServerConfig,
			httpserver.New,
		),
		fx.Provide(
			gqlresolver.New,
		),
		fx.Invoke(InvokeApp),
	}
	return fx.New(opts...), nil
}

func InvokeApp(
	server *httpserver.Server,
	resolver *gqlresolver.Resolver,
) {
	server.RegisterHandler("/api/query", resolver.Handler())
	server.RegisterHandler("/", ui.Handler("/"))
}
