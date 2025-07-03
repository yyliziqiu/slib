package slog

import (
	"testing"
)

func TestNew(t *testing.T) {
	config := Config{
		Console:     false,
		Path:        "/private/ws/self/slib/data/logs",
		RotateLevel: 0,
	}

	logger, err := New(config.Default())
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
	logger.Debug("debug")
}
