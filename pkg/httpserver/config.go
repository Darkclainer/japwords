package httpserver

type Config struct {
	Addr string
}

func (c *Config) Equal(other any) bool {
	oc, ok := other.(*Config)
	if !ok {
		return false
	}
	return c.Addr == oc.Addr
}
