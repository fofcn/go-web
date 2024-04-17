package task

import (
	"errors"
	"sync/atomic"
)

type Scheduler interface {
	Start() error
	Stop() error
	Schedule(task Task) (TaskFuture, error)
	RegisterWorker(worker Worker) error
	DeRegisterWorker(worker Worker) error
	HandleTaskCompletion(worker Worker, task Task, err error)
}

type defaultScheduler struct {
	workerlb LoadBalancer
	started  atomic.Bool
}

func NewScheduler() (Scheduler, error) {
	workerlb, err := NewLB("rr")
	if err != nil {
		return nil, err
	}

	return &defaultScheduler{
		workerlb: workerlb,
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
	w := s.workerlb.Select()
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
	s.workerlb.AddWorker(worker)
	return nil
}

func (s *defaultScheduler) DeRegisterWorker(worker Worker) error {
	s.workerlb.DelWorker(worker)
	return nil
}

func (s *defaultScheduler) HandleTaskCompletion(worker Worker, task Task, err error) {

}
