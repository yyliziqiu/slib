package sutil

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrNoValidRoundRobinNode = errors.New("no valid round robin node")
)

// NextValue interface
type NextValue interface {
	GetWeight() int
}

// Next Round Robin
type Next struct {
	list []*nextNode // 节点列表
	len  int         // 节点数量
	swrr bool        // 是否使用平滑加权轮询，如果所有节点的权重全部相同，则循环遍历
	sum  int         // 所有节点的权重不同时，记录所有节点的权重总和
	seq  int32       // 所有节点的权重相同时，计算本次轮询的节点位置
	mu   sync.Mutex  //
}

type nextNode struct {
	value  any // 值
	weight int // 权重
	status int // 状态值
}

// NewNext 创建轮询器
func NewNext() *Next {
	return &Next{list: make([]*nextNode, 0, 3)}
}

// Add 添加一个节点
func (t *Next) Add(value any, weight int) {
	t.mu.Lock()
	t.add(value, weight)
	t.mu.Unlock()
}

func (t *Next) add(value any, weight int) {
	// 添加节点
	t.list = append(t.list, &nextNode{value: value, weight: weight})
	t.len = len(t.list)

	// 判断节点权重是否相同
	if t.len > 1 && t.list[t.len-1].weight != t.list[t.len-2].weight {
		t.swrr = true
	}

	// 计算权重总和
	t.sum += weight
}

// AddValue 添加一个节点
func (t *Next) AddValue(v NextValue) {
	t.mu.Lock()
	t.add(v, v.GetWeight())
	t.mu.Unlock()
}

// Do 轮询
func (t *Next) Do() (v any, err error) {
	if t.swrr {
		t.mu.Lock()
		v, err = t.doSwrr()
		t.mu.Unlock()
	} else {
		v, err = t.doLoop()
	}
	return
}

// 平滑加权轮询
func (t *Next) doSwrr() (any, error) {
	var target *nextNode

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
func (t *Next) doLoop() (any, error) {
	if t.len == 0 {
		return nil, ErrNoValidRoundRobinNode
	}

	i := atomic.AddInt32(&t.seq, 1) & 0x7FFFFFFF

	return t.list[int(i)%t.len].value, nil
}

// MustDo 轮询
func (t *Next) MustDo() any {
	v, _ := t.Do()
	return v
}
