package fetcher

import "golang.org/x/exp/maps"

type Config struct {
	Headers map[string]string
}

func (c *Config) Equal(other any) bool {
	otherConfig, ok := other.(*Config)
	if !ok {
		return false
	}
	return maps.Equal(c.Headers, otherConfig.Headers)
}
