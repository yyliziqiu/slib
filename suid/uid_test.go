package uid

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUid_hex(t *testing.T) {
	for i := 1; i <= 16; i++ {
		assert.True(t, len(hex(1, i)) == i)
	}
}

func TestGetOrFail(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			_, _ = GetOrFail()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			_, _ = GetOrFail()
		}
	}()

	wg.Wait()

	t.Log("Completed")
}
