package scheduler

import (
	"sync"
	"sync/atomic"
)

const (
	defaultLBAlg = "rr"
)

type Scheduler struct {
	wm       *WorkerManager
	started  atomic.Bool
	executor Executor
}

var (
	scheduler *Scheduler
	once      sync.Once
)

func GetScheduler(lbAlg string) *Scheduler {
	once.Do(func() {
		var err error
		if len(lbAlg) == 0 {
			lbAlg = defaultLBAlg
		}
		scheduler, err = NewScheduler(lbAlg)
		if err != nil {
			panic(err)
		}
	})
	return scheduler
}

func NewScheduler(lbAlg string) (*Scheduler, error) {
	lb, err := NewLB(lbAlg)
	if err != nil {
		return nil, err
	}
	wm := NewWorkerManager(lb)

	executor := NewExecutor(wm)

	return &Scheduler{
		started:  atomic.Bool{},
		wm:       wm,
		executor: executor,
	}, nil
}

func (s *Scheduler) RegisterWorker(worker Worker) error {
	s.wm.AddWorker(worker)
	return nil
}

func (s *Scheduler) DeRegisterWorker(id WorkerId) error {
	s.wm.DelWorker(id)
	return nil
}

func (s *Scheduler) Start() error {
	return nil
}

func (s *Scheduler) Schedule(task Task) (TaskFuture, error) {
	return s.executor.Execute(task)
}

func (s *Scheduler) Stop() error {
	err := s.executor.Stop()
	if err != nil {
		return err
	}

	return s.wm.Stop()
}

func (s *Scheduler) GetWorkers() []Worker {
	return s.wm.GetWorkers()
}
