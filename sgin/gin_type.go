package sgin

type Config struct {
	Listen           string
	DisableAccessLog bool
}

func (c Config) Default() Config {
	if c.Listen == "" {
		c.Listen = ":80"
	}
	return c
}
