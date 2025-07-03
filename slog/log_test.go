package slog

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	path := "/private/ws/self/slib/data/logs"

	_ = os.RemoveAll(path)

	config := Config{
		Console:     false,
		Path:        path,
		RotateLevel: 5,
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
