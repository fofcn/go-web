package task

import (
	"errors"
	"sync"
)

type LoadBalancer interface {
	Select() Worker
	AddWorker(worker Worker)
	DelWorker(id WorkerId)
	GetWorkers() []Worker
}

type rr struct {
	idx     int
	workers []Worker
	mux     sync.Mutex
}

type workerweight struct {
	worker Worker
	weight int
}

type weight struct {
	idx     int
	workers []workerweight
	mux     sync.Mutex
}

var (
	lbRegistry = map[string]LoadBalancer{
		"rr":  &rr{},
		"wrr": &weight{},
	}
)

func NewLB(name string) (LoadBalancer, error) {
	alg, ok := lbRegistry[name]
	if !ok {
		return nil, errors.New("unsupported load balancer algorithm")
	}

	return alg, nil
}

func (r *rr) Select() Worker {
	r.mux.Lock()
	defer r.mux.Unlock()

	if len(r.workers) == 0 {
		return nil // 没有可用工作者
	}

	worker := r.workers[r.idx]
	r.idx = (r.idx + 1) % len(r.workers)
	return worker
}

func (r *rr) AddWorker(worker Worker) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, w := range r.workers {
		if w.GetAddr() == worker.GetAddr() {
			return
		}
	}
	r.workers = append(r.workers, worker)
}

func (r *rr) DelWorker(id WorkerId) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for i, w := range r.workers {
		if w.GetId() == id {
			r.workers = append(r.workers[:i], r.workers[i+1:]...)
			break
		}
	}
}

func (r *rr) GetWorkers() []Worker {
	return r.workers
}

func (w *weight) Select() Worker {
	w.mux.Lock()
	defer w.mux.Unlock()

	if len(w.workers) == 0 {
		return nil // 没有可用工作者
	}

	// 这里只是简单的选择下一个工作者，而没有按权重进行选择
	// 后面再改
	worker := w.workers[w.idx].worker
	w.idx = (w.idx + 1) % len(w.workers)
	return worker
}

func (w *weight) AddWorker(worker Worker) {
	w.mux.Lock()
	defer w.mux.Unlock()

	w.workers = append(w.workers, workerweight{worker: worker, weight: 1})
}

func (w *weight) DelWorker(id WorkerId) {
	w.mux.Lock()
	defer w.mux.Unlock()

	for i, ww := range w.workers {
		if ww.worker.GetId() == id {
			w.workers = append(w.workers[:i], w.workers[i+1:]...)
			break
		}
	}
}

func (w *weight) GetWorkers() []Worker {
	w.mux.Lock()
	defer w.mux.Unlock()

	workers := make([]Worker, len(w.workers))
	for _, ww := range w.workers {
		workers = append(workers, ww.worker)
	}

	return workers
}
