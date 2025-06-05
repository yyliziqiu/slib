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

var q3 = &Queue{
	step:  10,
	path:  "",
	debug: true,
	list:  []any{3, 4, 5, 6, 7, 8, nil, nil, nil, 1, 2},
	head:  9,
	tail:  6,
}

func echo(a ...any) {
	fmt.Println(a...)
}

func TestLen(t *testing.T) {
	echo("Q1 len: ", q1.len()) // 7
	echo("Q2 len: ", q2.len()) // 3
}

func TestPushAndPop(t *testing.T) {
	items := []int{8, 9, 10, 11, 12}
	for _, i := range items {
		q1.Push(i)
	}
	echo(q1.list)

	for !q1.empty() {
		q1.pop()
	}
	echo(q1.list)
}

func TestPushAndPop2(t *testing.T) {
	items := []int{4, 5, 6}
	for _, i := range items {
		q2.Push(i)
	}
	echo(q2.list)

	for !q2.empty() {
		q2.pop()
	}
	echo(q2.list)

	items = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, i := range items {
		q2.Push(i)
	}
	echo(q2.list)

	items = []int{10, 11, 12, 13}
	for _, i := range items {
		q2.Push(i)
	}
	echo(q2.list)
}

func TestGet(t *testing.T) {
	echo(q1.Get(0))  // err
	echo(q1.Get(1))  // 1
	echo(q1.Get(4))  // 4
	echo(q1.Get(8))  // err
	echo(q1.Get(10)) // err
}

func TestHeadItem(t *testing.T) {
	echo(q1.HeadItem()) // 1
	echo(q2.HeadItem()) // 1
}

func TestTailItem(t *testing.T) {
	echo(q1.TailItem()) // 7
	echo(q2.TailItem()) // 3
}

func TestEmpty(t *testing.T) {
	echo(q2.Empty())
	for !q2.empty() {
		q2.pop()
	}
	echo(q2.list)
	echo(q2.Empty())
}

func TestPops(t *testing.T) {
	result := q1.Pops(func(item any) bool {
		n := item.(int)
		return n <= 4
	})
	echo(q1.list)
	echo(result)

	result = q1.Pops(func(item any) bool {
		n := item.(int)
		return n <= 100
	})
	echo(q1.list)
	echo(result)

	echo("\n=========================================\n")

	result = q3.Pops(func(item any) bool {
		n := item.(int)
		return n <= 4
	})
	echo(q3.list)
	echo(result)

	result = q3.Pops(func(item any) bool {
		n := item.(int)
		return n <= 100
	})
	echo(q3.list)
	echo(result)
}

func TestSlideN(t *testing.T) {
	items := []int{8, 9, 10, 11, 12}
	for _, i := range items {
		q1.SlideN(i, 3)
	}
	echo(q1.list)

	items = []int{4, 5, 6, 7, 8, 9}
	for _, i := range items {
		q2.SlideN(i, 5)
	}
	echo(q2.list)
}

func TestSlide(t *testing.T) {
	n, last := q1.Slide(8, func(item any) bool {
		nn := item.(int)
		return nn <= 4
	})
	echo(q1.list)
	echo(n, last)

	n, last = q2.Slide(8, func(item any) bool {
		nn := item.(int)
		return nn <= 4
	})
	echo(q2.list)
	echo(n, last)

	n, last = q2.Slide(8, func(item any) bool {
		nn := item.(int)
		return nn <= 4
	})
	echo(q2.list)
	echo(n, last)
}
