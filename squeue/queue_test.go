package squeue

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/ssnap"
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
	path:  "/private/ws/self/slib/data/q2",
	debug: true,
	list:  []any{3, nil, nil, nil, nil, nil, nil, nil, nil, 1, 2},
	head:  9,
	tail:  1,
}

var q3 = &Queue{
	step:  10,
	path:  "/private/ws/self/slib/data/q3",
	debug: true,
	list:  []any{3, 4, 5, 6, 7, 8, nil, nil, nil, 1, 2},
	head:  9,
	tail:  6,
}

func echo(a ...any) {
	fmt.Println(a...)
}

func echo2(a ...any) {
	fmt.Print(a...)
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
	echo(q1.list) // [1 2 3 4 5 6 7 8 9 10 11 12 <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]

	for !q1.empty() {
		q1.pop()
	}
	echo(q1.list) // [<nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
}

func TestPushAndPop2(t *testing.T) {
	items := []int{4, 5, 6}
	for _, i := range items {
		q2.Push(i)
	}
	echo(q2.list) // [3 4 5 6 <nil> <nil> <nil> <nil> <nil> 1 2]

	for !q2.empty() {
		q2.pop()
	}
	echo(q2.list) // [<nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]

	items = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, i := range items {
		q2.Push(i)
	}
	echo(q2.list) // [8 9 <nil> <nil> 1 2 3 4 5 6 7]

	items = []int{10, 11, 12, 13}
	for _, i := range items {
		q2.Push(i)
	}
	echo(q2.list) // [1 2 3 4 5 6 7 8 9 10 11 12 13 <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
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
	echo(q2.Empty()) // false
	for !q2.empty() {
		q2.pop()
	}
	echo(q2.list)    // [<nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
	echo(q2.Empty()) // true
}

func TestPops(t *testing.T) {
	result := q1.Pops(func(item any) bool {
		n := item.(int)
		return n <= 4
	})
	echo(q1.list) // [<nil> <nil> <nil> <nil> <nil> 5 6 7 <nil> <nil> <nil>]
	echo(result)  // [1 2 3 4]

	result = q1.Pops(func(item any) bool {
		n := item.(int)
		return n <= 100
	})
	echo(q1.list) // [<nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
	echo(result)  // [5 6 7]

	echo("\n=========================================\n")

	result = q3.Pops(func(item any) bool {
		n := item.(int)
		return n <= 4
	})
	echo(q3.list) // [<nil> <nil> 5 6 7 8 <nil> <nil> <nil> <nil> <nil>]
	echo(result)  // [1 2 3 4]

	result = q3.Pops(func(item any) bool {
		n := item.(int)
		return n <= 100
	})
	echo(q3.list) // [<nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
	echo(result)  // [5 6 7 8]
}

func TestSlideN(t *testing.T) {
	items := []int{8, 9, 10, 11, 12}
	for _, i := range items {
		q1.SlideN(i, 3)
	}
	echo(q1.list) // [11 12 <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> 10]

	items = []int{4, 5, 6, 7, 8, 9}
	for _, i := range items {
		q2.SlideN(i, 5)
	}
	echo(q2.list) // [<nil> <nil> 5 6 7 8 9 <nil> <nil> <nil> <nil>]
}

func TestSlide(t *testing.T) {
	last, n := q1.Slide(8, func(item any) bool {
		nn := item.(int)
		return nn <= 4
	})
	echo(q1.list) // [<nil> <nil> <nil> <nil> <nil> 5 6 7 8 <nil> <nil>]
	echo(last, n) // 4 4

	last, n = q2.Slide(8, func(item any) bool {
		nn := item.(int)
		return nn <= 4
	})
	echo(q2.list) // [<nil> 8 <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
	echo(last, n) // 3 3

	last, n = q2.Slide(8, func(item any) bool {
		nn := item.(int)
		return nn <= 4
	})
	echo(q2.list) // [<nil> 8 8 <nil> <nil> <nil> <nil> <nil> <nil> <nil> <nil>]
	echo(last, n) // <nil> 0
}

func TestWalk(t *testing.T) {
	q1.Walk(func(item any) {
		echo2(item.(int), " ") // 1 2 3 4 5 6 7
	}, false)
	echo()

	q1.Walk(func(item any) {
		echo2(item.(int), " ") // 7 6 5 4 3 2 1
	}, true)
	echo()

	q3.Walk(func(item any) {
		echo2(item.(int), " ") // 1 2 3 4 5 6 7 8
	}, false)
	echo()

	q3.Walk(func(item any) {
		echo2(item.(int), " ") // 8 7 6 5 4 3 2 1
	}, true)
	echo()
}

func TestFind(t *testing.T) {
	item, _ := q1.Find(func(item any) bool {
		n := item.(int)
		echo2(n, " ")
		return n == 3
	}, false)
	echo()
	echo(item) // 3

	item, _ = q1.Find(func(item any) bool {
		n := item.(int)
		echo2(n, " ")
		return n == 100
	}, false)
	echo()
	echo(item) // <nil>

	item, _ = q1.Find(func(item any) bool {
		n := item.(int)
		echo2(n, " ")
		return n == 3
	}, true)
	echo()
	echo(item) // 3

	item, _ = q1.Find(func(item any) bool {
		n := item.(int)
		echo2(n, " ")
		return n == 100
	}, true)
	echo()
	echo(item) // <nil>
}

func TestFindAll(t *testing.T) {
	result := q1.FindAll(func(item any) bool {
		return item.(int) < 5
	})
	echo(result) // [1 2 3 4]
}

func TestTerminalN(t *testing.T) {
	result := q1.TerminalN(3, false)
	echo(result) // [1 2 3]
	result = q1.TerminalN(3, true)
	echo(result) // [7 6 5]

	result = q2.TerminalN(5, false)
	echo(result) // [1 2 3]
	result = q2.TerminalN(5, true)
	echo(result) // [3 2 1]
}

func TestTerminal(t *testing.T) {
	result := q1.Terminal(func(item any) bool {
		return item.(int) <= 3
	}, false)
	echo(result) // [1 2 3]
	result = q1.Terminal(func(item any) bool {
		return item.(int) <= 3
	}, true)
	echo(result) // []

	result = q3.Terminal(func(item any) bool {
		return item.(int) <= 3
	}, false)
	echo(result) // [1 2 3]
	result = q3.Terminal(func(item any) bool {
		return item.(int) >= 5
	}, true)
	echo(result) // [8 7 6 5]
}

func TestWindow(t *testing.T) {
	result := q1.Window(func(item any) bool {
		return item.(int) == 2
	}, func(item any) bool {
		return item.(int) == 4
	})
	echo(result) // [2 3]

	result = q1.Window(func(item any) bool {
		return item.(int) == 20
	}, func(item any) bool {
		return item.(int) == 4
	})
	echo(result) // []

	result = q3.Window(func(item any) bool {
		return item.(int) == 2
	}, func(item any) bool {
		return item.(int) >= 4
	})
	echo(result) // [2 3]

	result = q3.Window(func(item any) bool {
		return item.(int) >= 2
	}, func(item any) bool {
		return item.(int) == 20
	})
	echo(result) // [2 3 4 5 6 7 8]
}

func TestSaveAndLoad(t *testing.T) {
	_ = q3.Save()

	_ = q3.Load(1)

	echo(q3.list)
}

func TestWatcher(t *testing.T) {
	slog.Default, _ = slog.New(slog.Config{Console: true})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watchers := Watchers([]WatcherConfig{
		{
			Queue: q1,
			Item:  0,
			Path:  "/private/ws/self/slib/data/q1",
		},
		{
			Queue: q2,
			Item:  0,
			Poll:  5 * time.Second,
		},
	}...)

	err := ssnap.Watches(ctx, watchers)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(q1.list)
	fmt.Println(q2.list)

	time.Sleep(30 * time.Second)

	cancel()

	time.Sleep(time.Second)
}
