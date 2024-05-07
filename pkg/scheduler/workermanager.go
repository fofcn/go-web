package scheduler

import (
	"log"
	"time"
)

type WorkerManager struct {
	workers     WorkerStore
	lb          LoadBalancer
	timer       *time.Ticker
	healthTimer *time.Ticker
	done        chan bool
}

type WorkerManagerCfg struct {
	RedisConfig
}

func NewWorkerManager(lb LoadBalancer, wmCfg WorkerManagerCfg) *WorkerManager {
	store, err := NewRedisWorkerStore(&wmCfg.RedisConfig)
	if err != nil {
		log.Fatal("create worker store error.", err)
	}
	ww := &WorkerManager{
		lb:          lb,
		workers:     store,
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
	err := ww.workers.AddWorker(worker)
	if err != nil {
		log.Printf("add worker error %v.", err)
	}
	err = ww.workers.Heartbeat(worker)
	if err != nil {
		log.Printf("add worker error %v.", err)
	}
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

	for _, workerId := range workerIds {
		worker, err := ww.workers.GetWorker(workerId)
		if err != nil {
			continue
		}
		workers = append(workers, worker)
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
						if lastHeartbeatTime := worker.GetLastHeartbeat(); time.Duration(time.Since(lastHeartbeatTime).Seconds()) > time.Second*300 {
							log.Println("worker", workerId, "is not healthy, evict it")
							if err := ww.workers.DelWorker(workerId); err != nil {
								log.Println("failed to delete worker", workerId, "from worker store")
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
