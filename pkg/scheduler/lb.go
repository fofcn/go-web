package scheduler

import (
	"errors"
	"sync"
	"sync/atomic"
)

type LoadBalancer interface {
	Select(workers *sync.Map) Worker
}

type rr struct {
	idx *atomic.Int64
}

var (
	lbRegistry = map[string]LoadBalancer{
		"rr": newrr(),
	}
)

func NewLB(name string) (LoadBalancer, error) {
	alg, ok := lbRegistry[name]
	if !ok {
		return nil, errors.New("unsupported load balancer algorithm")
	}

	return alg, nil
}

func newrr() *rr {
	return &rr{
		idx: new(atomic.Int64),
	}
}

func (r *rr) Select(workers *sync.Map) Worker {
	var selected Worker
	workerCount := int64(0)
	newIndex := r.idx.Add(1)

	// 获取 workers 数量
	workers.Range(func(key, value interface{}) bool {
		workerCount++
		return true
	})

	// 如果没有 workers 返回 nil
	if workerCount == 0 {
		return nil
	}

	// 确保 selectedIndex 在 workers 数量范围内
	selectedIndex := newIndex % workerCount

	workerCount = 0 // 重置计数器用于下一次遍历
	workers.Range(func(key, value interface{}) bool {
		if workerCount == selectedIndex {
			// 断言 Worker 类型
			selected = value.(Worker)
			// 已找到，停止遍历
			return false
		}
		workerCount++
		return true
	})

	return selected
}
