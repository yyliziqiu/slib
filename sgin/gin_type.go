package sgin

type Config struct {
	Listen           string
	Tls              bool
	KeyFile          string
	CertFile         string
	DisableAccessLog bool
}

func (c Config) Default() Config {
	if c.Listen == "" {
		if c.Tls {
			c.Listen = ":443"
		} else {
			c.Listen = ":80"
		}
	}
	return c
}
