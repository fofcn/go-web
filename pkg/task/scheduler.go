package task

import (
	"errors"
	"sync"
	"sync/atomic"
)

type Scheduler interface {
	Start() error
	Stop() error
	Schedule(task Task) (TaskFuture, error)
	RegisterWorker(worker Worker) error
	DeRegisterWorker(id WorkerId) error
	GetWorkers() []Worker
	HandleTaskCompletion(worker Worker, task Task, err error)
}

type defaultScheduler struct {
	wm      *WorkerManager
	started atomic.Bool
}

const (
	defaultLBAlg = "rr"
)

var (
	scheduler Scheduler
	once      sync.Once
)

func GetScheduler() Scheduler {
	once.Do(func() {
		scheduler, _ = NewScheduler(defaultLBAlg)
	})

	return scheduler
}

func NewScheduler(lbAlg string) (Scheduler, error) {
	return &defaultScheduler{
		wm: NewWorkerManager(lbAlg),
	}, nil
}

func (s *defaultScheduler) Start() error {
	if !s.started.CompareAndSwap(false, true) {
		return nil
	}
	return errors.New("scheduler has started")
}

func (s *defaultScheduler) Stop() error {
	if !s.started.CompareAndSwap(true, false) {
		return nil
	}
	return nil
}

func (s *defaultScheduler) Schedule(task Task) (TaskFuture, error) {
	w := s.wm.SelectWorker()
	if w == nil {
		return nil, errors.New("no worker available")
	}

	future, err := w.Exec(task)
	if err != nil {
		return nil, err
	}
	return future, nil
}

func (s *defaultScheduler) RegisterWorker(worker Worker) error {
	s.wm.AddWorker(worker)
	return nil
}

func (s *defaultScheduler) DeRegisterWorker(id WorkerId) error {
	s.wm.DelWorker(id)
	return nil
}

func (s *defaultScheduler) GetWorkers() []Worker {
	return s.wm.GetWorkers()
}

func (s *defaultScheduler) HandleTaskCompletion(worker Worker, task Task, err error) {

}
