package ssnap

import (
	"testing"
	"time"
)

func TestDuplicate(t *testing.T) {
	data := map[string]string{
		"hello": "world",
	}

	err := Duplicate("/private/ws/self/slib/data/test.json", data, 10*time.Second)
	if err != nil {
		t.Error(err)
	}
}
