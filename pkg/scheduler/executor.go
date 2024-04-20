package scheduler

import (
	"errors"
	"sync/atomic"
)

type Executor interface {
	Start() error
	Stop() error
	Execute(task Task) (TaskFuture, error)
	GetTaskStatus(taskId string, workerId WorkerId) (*WorkerTaskResult, error)
	HandleTaskCompletion(worker Worker, task Task, err error)
}

type defaultScheduler struct {
	wm      *WorkerManager
	started atomic.Bool
}

func NewExecutor(vm *WorkerManager) Executor {
	return &defaultScheduler{
		wm: vm,
	}
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

func (s *defaultScheduler) Execute(task Task) (TaskFuture, error) {
	w := s.wm.SelectWorker()
	if w == nil {
		return nil, errors.New("no worker available")
	}

	future, err := w.Exec(task)
	if err != nil {
		return nil, err
	}
	task.SetWorkerId(w.GetId())
	return future, nil
}

func (s *defaultScheduler) HandleTaskCompletion(worker Worker, task Task, err error) {

}

func (s *defaultScheduler) GetTaskStatus(taskId string, workerId WorkerId) (*WorkerTaskResult, error) {
	worker, err := s.wm.GetWorker(workerId)
	if err != nil {
		return nil, err
	}

	return worker.GetTaskStatus(taskId)
}
