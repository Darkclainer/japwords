package fxapp

import (
	"context"

	"go.uber.org/fx"

	"github.com/Darkclainer/japwords/pkg/anki"
	"github.com/Darkclainer/japwords/pkg/config"
)

func NewAnki(configManager *config.Manager, LC fx.Lifecycle) (*anki.Anki, *anki.ConfigReloader, error) {
	client := anki.NewAnki(anki.DefaultStatefullClientConstructor)
	reloader, err := anki.NewConfigReloader(client, configManager)
	if err != nil {
		return nil, nil, err
	}
	LC.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			client.Stop()
			return nil
		},
	})
	return client, reloader, nil
}
