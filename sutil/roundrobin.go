package sutil

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrNoValidRoundRobinNode = errors.New("no valid round robin node")
)

// RoundRobin 平滑加权轮询器
type RoundRobin struct {
	list []*RoundRobinNode // 节点列表
	len  int               // 节点数量
	swrr bool              // 是否使用平滑加权轮询，如果所有节点的权重全部相同，则循环遍历
	sum  int               // 所有节点的权重不同时，记录所有节点的权重总和
	seq  int32             // 所有节点的权重相同时，计算本次轮询的节点位置
	mu   sync.Mutex        //
}

type RoundRobinNode struct {
	value  any
	weight int
	status int
}

type RoundRobinValue interface {
	GetWeight() int
}

// NewRoundRobin 创建轮询器
func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		list: make([]*RoundRobinNode, 0, 3),
	}
}

// Add 添加一个节点
func (t *RoundRobin) Add(value any, weight int) {
	t.mu.Lock()
	t.add(value, weight)
	t.mu.Unlock()
}

func (t *RoundRobin) add(value any, weight int) {
	// 创建节点
	node := &RoundRobinNode{
		value:  value,
		weight: weight,
	}

	// 判断节点权重是否相同
	if t.len > 1 && node.weight != t.list[t.len-1].weight {
		t.swrr = true
	}

	// 添加节点
	t.list = append(t.list, node)
	t.len = len(t.list)
	t.sum += weight
}

// AddValue 添加一个节点
func (t *RoundRobin) AddValue(v RoundRobinValue) {
	t.mu.Lock()
	t.add(v, v.GetWeight())
	t.mu.Unlock()
}

// Next 轮询
func (t *RoundRobin) Next() (v any, err error) {
	if t.swrr {
		t.mu.Lock()
		v, err = t.nextSwrr()
		t.mu.Unlock()
	} else {
		v, err = t.nextLoop()
	}
	return
}

// 平滑加权轮询
func (t *RoundRobin) nextSwrr() (any, error) {
	var target *RoundRobinNode

	// 将所有节点的状态值加上该节点权重，并选出状态值最大的节点
	for _, node := range t.list {
		node.status += node.weight
		if target == nil || node.status > target.status {
			target = node
		}
	}

	// 无有效节点
	if target == nil {
		return nil, ErrNoValidRoundRobinNode
	}

	// 将选中节点的状态值减去所有节点的权重总和
	target.status -= t.sum

	return target.value, nil
}

// 循环遍历
func (t *RoundRobin) nextLoop() (any, error) {
	if t.len == 0 {
		return nil, ErrNoValidRoundRobinNode
	}

	i := atomic.AddInt32(&t.seq, 1) & 0x7FFFFFFF

	return t.list[int(i)%t.len].value, nil
}

// MustNext 轮询
func (t *RoundRobin) MustNext() any {
	v, _ := t.Next()
	return v
}
