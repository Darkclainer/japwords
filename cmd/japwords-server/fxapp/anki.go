package fxapp

import (
	"github.com/Darkclainer/japwords/pkg/anki"
	"github.com/Darkclainer/japwords/pkg/config"
)

func NewAnki(configManager *config.Manager) (*anki.Anki, *anki.ConfigReloader, error) {
	client := anki.NewAnki(anki.DefaultClientConstructor)
	reloader, err := anki.NewConfigReloader(client, configManager)
	if err != nil {
		return nil, nil, err
	}
	return client, reloader, nil
}
