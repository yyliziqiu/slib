package sutil

import (
	"sync"
	"sync/atomic"
)

// Rr Round Robin
type Rr struct {
	list []*rrNode  // 节点列表
	len  int        // 节点数量
	swrr bool       // 是否使用平滑加权轮询，如果所有节点的权重全部相同，则循环逐个遍历
	sum  int        // 所有节点的权重不同时，记录所有节点的权重总和
	seq  int32      // 所有节点的权重相同时，计算本次轮询的节点位置
	mu   sync.Mutex //
}

type rrNode struct {
	value  any // 值
	weight int // 权重
	status int // 状态值
}

type RrNode interface {
	GetWeight() int
}

// NewRr 创建轮询器
func NewRr() *Rr {
	return &Rr{list: make([]*rrNode, 0, 3)}
}

func (t *Rr) AddNodes(nodes []RrNode) {
	t.mu.Lock()
	for _, node := range nodes {
		t.add(node, node.GetWeight())
	}
	t.mu.Unlock()
}

// Add 添加一个节点
func (t *Rr) Add(value any, weight int) {
	t.mu.Lock()
	t.add(value, weight)
	t.mu.Unlock()
}

func (t *Rr) add(value any, weight int) {
	// 添加节点
	t.list = append(t.list, &rrNode{value: value, weight: weight})
	t.len = len(t.list)

	// 判断节点权重是否相同
	if t.len > 1 && t.list[t.len-1].weight != t.list[t.len-2].weight {
		t.swrr = true
	}

	// 计算权重总和
	t.sum += weight
}

// Next 轮询
func (t *Rr) Next() (v any) {
	if t.swrr {
		t.mu.Lock()
		v = t.nextBySwrr()
		t.mu.Unlock()
	} else {
		v = t.nextByIncr()
	}
	return
}

// 平滑加权轮询
func (t *Rr) nextBySwrr() any {
	var target *rrNode

	// 将所有节点的状态值加上该节点权重，并选出状态值最大的节点
	for _, node := range t.list {
		node.status += node.weight
		if target == nil || node.status > target.status {
			target = node
		}
	}

	if target == nil {
		return nil
	}

	// 将选中节点的状态值减去所有节点的权重总和
	target.status -= t.sum

	return target.value
}

// 循环逐个遍历
func (t *Rr) nextByIncr() any {
	if t.len == 0 {
		return nil
	}

	i := atomic.AddInt32(&t.seq, 1) & 0x7FFFFFFF

	return t.list[int(i)%t.len].value
}
