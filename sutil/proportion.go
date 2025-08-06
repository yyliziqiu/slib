package sutil

import (
	"sync"
)

// Proportion 比例均衡器
type Proportion struct {
	scale int        // 此值越大 reset 周期越大
	ratio float64    // Next() 返回 true 的比例
	total int        // Next() 调用次数
	count int        // Next() 返回 true 的次数
	mu    sync.Mutex //
}

func NewProportion(percent float64) *Proportion {
	return &Proportion{
		scale: 1000000,
		ratio: percent / 100,
		total: 0,
		count: 0,
	}
}

func (t *Proportion) Next() bool {
	t.mu.Lock()
	ret := t.next()
	t.mu.Unlock()
	return ret
}

func (t *Proportion) next() bool {
	if t.total > t.scale {
		t.total = 0
		t.count = 0
	}

	t.total++

	if float64(t.count)/float64(t.total) >= t.ratio {
		return false
	}

	t.count++

	return true
}
