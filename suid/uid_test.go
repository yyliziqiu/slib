package uid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUid_hex(t *testing.T) {
	for i := 1; i <= 16; i++ {
		assert.True(t, len(hex(1, i)) == i)
	}
}
