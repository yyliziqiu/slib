package slog

import "time"

const (
	TextFormat = "text"
	JsonFormat = "json"
)

type Config struct {
	Console       bool
	Path          string
	Name          string
	Level         string
	ShowCaller    bool
	DataFormat    string
	DateFormat    string
	MaxAge        time.Duration
	RotationTime  time.Duration
	RotationLevel int
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
	if c.RotationTime == 0 {
		c.RotationTime = 24 * time.Hour
	}
	return c
}
