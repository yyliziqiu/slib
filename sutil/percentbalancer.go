package sutil

import (
	"sync"
)

// PercentBalancer 百分比均衡器
type PercentBalancer struct {
	scale int        // 此值越大 reset 周期越大
	bound float64    // Next() 返回 true 的比例
	total int        // Next() 调用次数
	count int        // Next() 返回 true 的次数
	mu    sync.Mutex //
}

func NewPercentBalancer(p float64) *PercentBalancer {
	return &PercentBalancer{
		scale: 1e6,
		bound: p / 100,
		total: 0,
		count: 0,
	}
}

func (t *PercentBalancer) Next() bool {
	t.mu.Lock()
	ret := t.next()
	t.mu.Unlock()

	return ret
}

func (t *PercentBalancer) next() bool {
	if t.total > t.scale {
		t.total = 0
		t.count = 0
	}

	t.total++
	if float64(t.count)/float64(t.total) >= t.bound {
		return false
	}

	t.count++

	return true
}
