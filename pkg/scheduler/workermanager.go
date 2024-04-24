package scheduler

import (
	"errors"
	"time"
)

type WorkerManager struct {
	workers     WorkerStore
	lb          LoadBalancer
	timer       *time.Ticker
	healthTimer *time.Ticker
	done        chan bool
}

func NewWorkerManager(lb LoadBalancer) *WorkerManager {
	ww := &WorkerManager{
		lb:          lb,
		workers:     NewRedisWorkerStore(nil),
		timer:       time.NewTicker(5 * time.Second),
		healthTimer: time.NewTicker(5 * time.Second),
		done:        make(chan bool),
	}

	go func() {
		ww.evictWorker()
	}()

	return ww
}

func (ww *WorkerManager) AddWorker(worker Worker) {
	ww.workers.Store(worker.GetId(), worker)
}

func (ww *WorkerManager) DelWorker(id WorkerId) {
	ww.workers.Delete(id)
}

func (ww *WorkerManager) GetWorkers() []Worker {
	var workers []Worker
	ww.workers.Range(func(key, value any) bool {
		workers = append(workers, value.(Worker))
		return true
	})

	return workers
}

func (ww *WorkerManager) SelectWorker() Worker {
	return ww.lb.Select(ww.workers)
}

func (ww *WorkerManager) GetWorker(workerId WorkerId) (Worker, error) {
	if worker, ok := ww.workers.Load(workerId); ok {
		return worker.(Worker), nil
	}
	return nil, errors.New("no such worker registerd")
}

func (ww *WorkerManager) Close() {
	ww.timer.Stop()
	ww.healthTimer.Stop()
	ww.done <- true
}
func (ww *WorkerManager) evictWorker() {
	for {
		select {
		case <-ww.done:
			return
		case <-ww.healthTimer.C:
			println("health check timer")
			ww.workers.Range(func(key, value any) bool {
				worker := value.(Worker)
				err := worker.CheckStatus()
				if err != nil {
					if worker.IncrErrorCounter() >= 3 {
						ww.DelWorker(worker.GetId())
					}
				} else {
					println(worker.Status().String())
				}
				return true
			})
		}
	}
}

func (ww *WorkerManager) Stop() error {
	ww.done <- true
	defer ww.timer.Stop()
	defer ww.healthTimer.Stop()
	return nil
}
