package gqlresolver

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"go.uber.org/fx"

	"github.com/Darkclainer/japwords/graphql/gqlgenerated"
	"github.com/Darkclainer/japwords/pkg/anki"
	"github.com/Darkclainer/japwords/pkg/config"
	"github.com/Darkclainer/japwords/pkg/multidict"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	configManager *config.Manager
	multiDict     *multidict.MultiDict
	ankiClient    *anki.Anki
	ankiConfig    *anki.ConfigReloader
}

type In struct {
	fx.In

	ConfigManager *config.Manager
	MultiDict     *multidict.MultiDict
	AnkiClient    *anki.Anki
	AnkiConfig    *anki.ConfigReloader
}

func New(in In) (*Resolver, error) {
	return &Resolver{
		configManager: in.ConfigManager,
		multiDict:     in.MultiDict,
		ankiClient:    in.AnkiClient,
		ankiConfig:    in.AnkiConfig,
	}, nil
}

func (r *Resolver) Handler() http.Handler {
	h := handler.New(gqlgenerated.NewExecutableSchema(gqlgenerated.Config{
		Resolvers: r,
	}))
	h.AddTransport(transport.POST{})
	const (
		queryCacheSize              = 1000
		autoPersistedQueryCacheSize = 100
	)
	h.SetQueryCache(lru.New(queryCacheSize))
	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(autoPersistedQueryCacheSize),
	})
	return h
}
