package sutil

import (
	"sync"
)

// PercentBalancer 百分比均衡器
type PercentBalancer struct {
	scale int        // 此值越大 reset 周期越大
	ratio float64    // Take() 返回 true 的比例
	total int        // Take() 调用次数
	count int        // Take() 返回 true 的次数
	mu    sync.Mutex //
}

func NewPercentBalancer(target float64) *PercentBalancer {
	return &PercentBalancer{
		scale: 1000000,
		ratio: target / 100,
		total: 0,
		count: 0,
	}
}

func (t *PercentBalancer) Reset() {
	t.mu.Lock()
	t.reset()
	t.mu.Unlock()
}

func (t *PercentBalancer) reset() {
	t.total = 0
	t.count = 0
}

func (t *PercentBalancer) Take() bool {
	t.mu.Lock()
	ret := t.take()
	t.mu.Unlock()
	return ret
}

func (t *PercentBalancer) take() bool {
	if t.total > t.scale {
		t.reset()
	}

	t.total++

	if float64(t.count)/float64(t.total) >= t.ratio {
		return false
	}

	t.count++

	return true
}
