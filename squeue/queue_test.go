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

func TestLen(t *testing.T) {
	t.Log("Q1 len: ", q1.len())
	t.Log("Q2 len: ", q2.len())
}

func TestPushAndPop(t *testing.T) {
	items := []int{7, 8, 9, 10}
	for _, i := range items {
		q1.Push(i)
	}
	fmt.Println(q1.list)

	for !q1.empty() {
		q1.pop()
	}
	fmt.Println(q1.list)
}

func TestPushAndPop2(t *testing.T) {
	items := []int{7, 8, 9, 10}
	for _, i := range items {
		q2.Push(i)
	}
	fmt.Println(q2.list)

	for !q2.empty() {
		q2.pop()
	}
	fmt.Println(q2.list)

	items = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, i := range items {
		q2.Push(i)
	}
	fmt.Println(q2.list)
}
