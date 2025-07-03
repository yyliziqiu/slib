package squeue

// EnableDebug 开启 debug 模式，该模式下会输出队列的操作日志
func (q *Queue) EnableDebug() *Queue {
	q.debug = true
	return q
}

// Get 获取指定下标的元素
func (q *Queue) Get(i int) (any, error) {
	q.mu.RLock()
	item, err := q.get(i)
	q.mu.RUnlock()

	return item, err
}

// HeadItem 获取头元素
func (q *Queue) HeadItem() (any, error) {
	q.mu.RLock()
	item, err := q.get(q.head)
	q.mu.RUnlock()

	return item, err
}

// TailItem 获取尾元素
func (q *Queue) TailItem() (any, error) {
	q.mu.RLock()
	item, err := q.get(q.tailPrev())
	q.mu.RUnlock()

	return item, err
}

// Status 获取队列状态
func (q *Queue) Status() string {
	q.mu.RLock()
	s := q.status()
	q.mu.RUnlock()

	return s
}

// Empty 判断队列是否为空
func (q *Queue) Empty() bool {
	q.mu.RLock()
	e := q.empty()
	q.mu.RUnlock()

	return e
}

// Cap 获取队列容量
func (q *Queue) Cap() int {
	q.mu.RLock()
	c := q.cap() - 1
	q.mu.RUnlock()

	return c
}

// Len 获取队列长度
func (q *Queue) Len() int {
	q.mu.RLock()
	l := q.len()
	q.mu.RUnlock()

	return l
}

// Push 从队列尾向队列中添加一个元素
func (q *Queue) Push(item any) {
	q.mu.Lock()
	q.push(item)
	q.mu.Unlock()
}

// Pop 从队列头弹出一个元素
func (q *Queue) Pop() (any, bool) {
	q.mu.Lock()
	item, ok := q.pop()
	q.mu.Unlock()

	return item, ok
}

// Filter 元素符合条件返回 true，否则返回 false
type Filter func(item any) bool

// Pops 从队列头开始弹出所有符合条件的元素，直到遇到第一个不符合条件的元素停止
func (q *Queue) Pops(filter Filter) []any {
	result := make([]any, 0, 4)

	q.mu.Lock()
	for q.head != q.tail {
		item := q.list[q.head]
		ok := filter(item)
		if !ok {
			break
		}
		result = append(result, item)
		q.list[q.head] = nil
		q.head = q.headNext()
	}
	q.mu.Unlock()

	return result
}

// SlideN 类似于滑动窗口，在队列尾添加一个元素，如果添加完元素队列长度大于 n，则删除前面的元素，最后只保留队列后 n 个元素
// 第一个返回值表示最后一个删除的元素
// 第二个返回值表示是窗口否发生了滑动
func (q *Queue) SlideN(item any, n int) (any, bool) {
	q.mu.Lock()

	q.push(item)

	slide := false
	for q.len() > n {
		slide = true
		item, _ = q.pop()
	}

	q.mu.Unlock()

	if slide && q.debug {
		q.print("slide")
	}

	return item, slide
}

// Remove 需要删除的元素返回 true，否则返回 false
type Remove func(item any) bool

// Slide  类似于滑动窗口，在队列尾添加一个元素，并从队列头开始直到第一个不需要删除的元素出现，该元素前面的元素全部删除
// 第一个返回值表示最后一个被删除的元素
// 第二个返回值表示被删除的元素个数
func (q *Queue) Slide(item any, remove Remove) (last any, n int) {
	q.mu.Lock()
	q.push(item)
	for !q.empty() && remove(q.list[q.head]) {
		n++
		last, _ = q.pop()
	}
	q.mu.Unlock()

	if n > 0 && q.debug {
		q.print("slide")
	}

	return
}

// Walk 遍历队列
// reverse false：从头到尾遍历，true：从尾到头遍历
func (q *Queue) Walk(f func(item any), reverse bool) {
	q.mu.RLock()
	if reverse {
		for i := q.tailPrev(); i != q.headPrev(); i = q.prev(i) {
			f(q.list[i])
		}
	} else {
		for i := q.head; i != q.tail; i = q.next(i) {
			f(q.list[i])
		}
	}
	q.mu.RUnlock()
}

// Find 遍历队列，返回第一个符合条件的元素
// reverse false：从头到尾遍历，true：从尾到头遍历
func (q *Queue) Find(filter Filter, reverse bool) (any, int) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if reverse {
		for i := q.tailPrev(); i != q.headPrev(); i = q.prev(i) {
			if item := q.list[i]; filter(item) {
				return item, i
			}
		}
	} else {
		for i := q.head; i != q.tail; i = q.next(i) {
			if item := q.list[i]; filter(item) {
				return item, i
			}
		}
	}

	return nil, 0
}

// FindAll 遍历队列，返回全部符合条件的元素
func (q *Queue) FindAll(f Filter) []any {
	all := make([]any, 0)

	q.mu.RLock()
	for i := q.head; i != q.tail; i = q.next(i) {
		if item := q.list[i]; f(item) {
			all = append(all, item)
		}
	}
	q.mu.RUnlock()

	return all
}

// TerminalN 获取队列前/后 n 个 item
func (q *Queue) TerminalN(n int, reverse bool) []any {
	items := make([]any, 0, n)

	q.mu.RLock()

	if n > q.len() {
		n = q.len()
	}

	if reverse {
		for i, j := 0, q.tailPrev(); i < n && j != q.headPrev(); i, j = i+1, q.prev(j) {
			items = append(items, q.list[j])
		}
	} else {
		for i, j := 0, q.head; i < n && j != q.tail; i, j = i+1, q.next(j) {
			items = append(items, q.list[j])
		}
	}

	q.mu.RUnlock()

	return items
}

// Terminal 获取队列前/后多个符合条件的 item，遇到第一个不符合条件的 item 停止遍历
func (q *Queue) Terminal(filter Filter, reverse bool) []any {
	items := make([]any, 0)

	q.mu.RLock()
	if reverse {
		for i := q.tailPrev(); i != q.headPrev(); i = q.prev(i) {
			item := q.list[i]
			if !filter(item) {
				break
			}
			items = append(items, item)
		}
	} else {
		for i := q.head; i != q.tail; i = q.next(i) {
			item := q.list[i]
			if !filter(item) {
				break
			}
			items = append(items, item)
		}
	}
	q.mu.RUnlock()

	return items
}

// Window
// 返回结果包含 bgn item，不包含 end item
func (q *Queue) Window(bgn Filter, end Filter) []any {
	start := false
	result := make([]any, 0)

	q.mu.RLock()
	for i := q.head; i != q.tail; i = q.next(i) {
		item := q.list[i]
		if !start && bgn(item) {
			start = true
		}
		if start && end(item) {
			break
		}
		if start {
			result = append(result, item)
		}
	}
	q.mu.RUnlock()

	return result
}

// Reset 重置队列
func (q *Queue) Reset(data []any) {
	q.mu.Lock()
	q.reset(data)
	q.mu.Unlock()
}

// CopyList 复制列表
func (q *Queue) CopyList() []any {
	q.mu.RLock()
	list := q.copyList()
	q.mu.RUnlock()

	return list
}
