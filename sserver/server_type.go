package sserver

import (
	"github.com/sirupsen/logrus"
)

var (
	_errorLogger  *logrus.Logger
	_accessLogger *logrus.Logger
)

type Config struct {
	Listen           string
	ErrorLog         string
	AccessLog        string
	DisableAccessLog bool
}

func (c Config) Default() Config {
	if c.Listen == "" {
		c.Listen = ":80"
	}
	if c.ErrorLog == "" {
		c.ErrorLog = "web-error"
	}
	if c.AccessLog == "" {
		c.AccessLog = "web-access"
	}
	return c
}
