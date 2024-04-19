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
	selectedIndex := r.idx.Add(1)

	workers.Range(func(key, value interface{}) bool {
		workerCount++
		if workerCount == selectedIndex {
			selected = value.(Worker)
			return false // 停止遍历
		}
		return true // 继续遍历
	})

	if selected == nil && workerCount > 0 {
		// 如果 selectedIndex 超出范围，从头开始
		return r.Select(workers)
	}
	return selected
}
