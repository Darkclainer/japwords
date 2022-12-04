package fxapp

import (
	"japwords/pkg/config"
	"japwords/pkg/wadoku"
)

func NewWadoku(wadokuClient wadoku.BasicDict, uc *config.UserConfig) (*wadoku.Wadoku, error) {
	return wadoku.New(wadokuClient, uc.Dictionary.Wadoku.URL), nil
}
