package fxapp

import (
	"japwords/pkg/config"
	"japwords/pkg/jisho"
)

func NewJisho(jishoClient jisho.BasicDict, uc *config.UserConfig) (*jisho.Jisho, error) {
	return jisho.New(jishoClient, uc.Dictionary.Jisho.URL), nil
}
