package sutil

import (
	"fmt"
	"testing"
)

func TestNext(t *testing.T) {
	m := NewNext()

	m.add(1, 1)
	m.add(2, 2)
	m.add(3, 3)

	for i := 0; i < 30; i++ {
		fmt.Print(m.MustDo(), " ")
	}

	fmt.Println()
}
