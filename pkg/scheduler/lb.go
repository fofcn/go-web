package scheduler

import (
	"errors"
	"sync/atomic"
)

type LoadBalancer interface {
	Select(workers []WorkerId) WorkerId
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

func (r *rr) Select(workers []WorkerId) WorkerId {
	workerCount := int64(len(workers))
	newIndex := r.idx.Add(1)

	// 如果没有 workers 返回 nil
	if workerCount == 0 {
		return ""
	}

	// 确保 selectedIndex 在 workers 数量范围内
	return workers[newIndex%workerCount]
}
