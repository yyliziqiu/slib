package sutil

import (
	"fmt"
	"testing"
)

func TestRtp(t *testing.T) {
	r := NewRtp(34)

	a, c := 0, 5235252
	for i := 0; i < c; i++ {
		ok := r.Do()
		if ok {
			a++
		}
	}

	fmt.Printf("a = %d, c = %f\n", a, float64(a)/float64(c))
}
