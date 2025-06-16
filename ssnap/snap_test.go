package ssnap

import (
	"testing"
	"time"
)

func TestDuplicate(t *testing.T) {
	path := "/private/ws/self/slib/data/test.json"

	data := map[string]string{
		"hello": "world",
	}

	err := Save(path, data)
	if err != nil {
		t.Error(err)
	}

	err = Duplicate(path, data, 10*time.Second)
	if err != nil {
		t.Error(err)
	}
}
