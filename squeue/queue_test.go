package squeue

import (
	"fmt"
	"testing"
)

var q1 = &Queue{
	step:  10,
	path:  "",
	debug: true,
	list:  []any{nil, 1, 2, 3, 4, 5, 6, 7, nil, nil, nil},
	head:  1,
	tail:  8,
}

var q2 = &Queue{
	step:  10,
	path:  "",
	debug: true,
	list:  []any{3, nil, nil, nil, nil, nil, nil, nil, nil, 1, 2},
	head:  9,
	tail:  1,
}

func TestLen(t *testing.T) {
	t.Log("Q1 len: ", q1.len()) // 7
	t.Log("Q2 len: ", q2.len()) // 3
}

func TestPushAndPop(t *testing.T) {
	items := []int{8, 9, 10, 11, 12}
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
	items := []int{4, 5, 6}
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
