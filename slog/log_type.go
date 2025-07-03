package slog

import (
	"time"

	"github.com/sirupsen/logrus"
)

const (
	TextFormat = "text"
	JsonFormat = "json"
)

type LevelDispatch map[string][]logrus.Level

type Config struct {
	Console        bool
	Path           string
	Name           string
	Level          string
	ShowCaller     bool
	DataFormat     string
	DateFormat     string
	MaxAge         time.Duration
	RotateTime     time.Duration
	RotateLevel    int
	RotateTimezone string
}

func (c Config) Default() Config {
	if c.Path == "" {
		c.Console = true
	}
	if c.Name == "" {
		c.Name = "app"
	}
	if c.Level == "" {
		c.Level = "debug"
	}
	if c.DataFormat == "" {
		c.DataFormat = TextFormat
	}
	if c.DateFormat == "" {
		c.DateFormat = time.DateTime
	}
	if c.MaxAge == 0 {
		c.MaxAge = 7 * 24 * time.Hour
	}
	if c.RotateTime == 0 {
		c.RotateTime = 24 * time.Hour
	}
	if c.RotateLevel == 0 {
		c.RotateLevel = 2
	}
	if c.RotateTimezone == "" {
		c.RotateTimezone = "Asia/Shanghai"
	}
	return c
}
