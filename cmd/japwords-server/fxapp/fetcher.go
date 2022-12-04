package fxapp

import (
	"japwords/pkg/config"
	"japwords/pkg/fetcher"
)

func NewFetcher(uc *config.UserConfig) (*fetcher.Fetcher, error) {
	headers := uc.Dictionary.Headers
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["User-Agent"] = uc.Dictionary.UserAgent
	fetcherClient, err := fetcher.New(&fetcher.Config{
		Headers: headers,
	})
	if err != nil {
		return nil, err
	}
	return fetcherClient, nil
}
