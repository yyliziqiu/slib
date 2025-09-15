package sutil

import (
	"fmt"
	"testing"
)

func TestNext(t *testing.T) {
	m := NewRoundRobin()

	m.Add(1, 1)
	m.Add(2, 2)
	m.Add(3, 3)

	fmt.Println(m.swrr)

	for i := 0; i < 30; i++ {
		fmt.Print(m.MustNext(), " ")
	}

	fmt.Println()
}
