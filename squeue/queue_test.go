package squeue

import (
	"fmt"
	"testing"
)

var q1 = &Queue{
	step:  10,
	path:  "",
	debug: true,
	list:  []any{0, 1, 2, 3, 4, 5, 6, 7, 0, 0},
	head:  1,
	tail:  8,
}

var q2 = &Queue{
	step:  10,
	path:  "",
	debug: true,
	list:  []any{10, 0, 0, 0, 0, 0, 0, 0, 8, 9},
	head:  8,
	tail:  1,
}

func TestQueue(t *testing.T) {
	t.Log("Q1 len: ", q1.len())
	t.Log("Q2 len: ", q2.len())

	q1.push(8)
	q1.push(9)
	q1.push(10)
	fmt.Println(q1.list)
}
