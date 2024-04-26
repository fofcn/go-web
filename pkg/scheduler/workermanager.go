package scheduler

import (
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
	ww.workers.AddWorker(worker)
}

func (ww *WorkerManager) DelWorker(id WorkerId) {
	ww.workers.DelWorker(id)
}

func (ww *WorkerManager) GetWorkers() []Worker {
	workerIds, err := ww.workers.GetWorkerIds()
	if err != nil {
		return nil
	}

	workers := make([]Worker, 0, len(workerIds))

	for idx, workerId := range workerIds {
		worker, err := ww.workers.GetWorker(workerId)
		if err != nil {
			continue
		}
		workers[idx] = worker
	}

	return workers
}

func (ww *WorkerManager) SelectWorker() Worker {
	workerIds, err := ww.workers.GetWorkerIds()
	if err != nil {
		return nil
	}

	workerId := ww.lb.Select(workerIds)
	worker, err := ww.workers.GetWorker(workerId)
	if err != nil {
		return nil
	}

	return worker
}

func (ww *WorkerManager) GetWorker(workerId WorkerId) (Worker, error) {
	return ww.workers.GetWorker(workerId)
}

func (ww *WorkerManager) Heartbeat(worker Worker) error {
	return ww.workers.Heartbeat(worker)
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
			// get worker id list from worker store
			if workerIds, err := ww.workers.GetWorkerIds(); err == nil {
				for _, workerId := range workerIds {
					if worker, err := ww.workers.GetWorker(workerId); err == nil {
						if lastHeartbeatTime := worker.GetLastHeartbeat(); time.Since(lastHeartbeatTime) > time.Second*300 {
							println("worker", workerId, "is not healthy, evict it")
							if err := ww.workers.DelWorker(workerId); err != nil {
								println("failed to delete worker", workerId, "from worker store")
							}
						}
					}
				}
			}
			// get worker last ping time from worker store
			// if timeout then remove worker from worker store

		}
	}
}

func (ww *WorkerManager) Stop() error {
	ww.done <- true
	defer ww.timer.Stop()
	defer ww.healthTimer.Stop()
	return nil
}
