package slog

import (
	"testing"
)

func TestNew(t *testing.T) {
	config := Config{
		Console:       false,
		Path:          "/private/ws/self/slib/logs",
		RotationLevel: 0,
	}

	logger, err := New(config.Default())
	if err != nil {
		t.Fatal(err)
	}

	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
	// logger.Fatal("fatal")
	// logger.Panic("panic")
}
