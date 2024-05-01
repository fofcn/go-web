package scheduler

import (
	"go-web/pkg/config"
	"sync"
	"sync/atomic"
	"time"
)

const ()

type Scheduler struct {
	wm       *WorkerManager
	started  atomic.Bool
	executor Executor
	store    TaskStore
}

var (
	scheduler *Scheduler
	once      sync.Once
)

func GetScheduler(cfg *config.Scheduler) *Scheduler {
	once.Do(func() {
		var err error
		scheduler, err = NewScheduler(cfg)
		if err != nil {
			panic(err)
		}
	})
	return scheduler
}

func NewScheduler(cfg *config.Scheduler) (*Scheduler, error) {
	lb, err := NewLB(cfg.WorkerConfig.LoadBalancer)
	if err != nil {
		return nil, err
	}

	redisCfg := RedisConfig{
		ClusterMode:  cfg.Redis.ClusterMode,
		Addrs:        cfg.Redis.Addrs,
		DB:           cfg.Redis.DB,
		ClientName:   cfg.Redis.ClientName,
		MinIdleConns: cfg.Redis.Pool.MaxIdle,
		MaxIdleConns: cfg.Redis.Pool.MaxActive,
		MaxRetries:   cfg.Redis.MaxRetries,
		PoolSize:     cfg.Redis.Pool.Size,
		PoolTimeout:  time.Duration(cfg.Redis.Pool.IdleTimeout) * time.Second,
		Password:     cfg.Redis.Password,
	}
	wmCfg := WorkerManagerCfg{RedisConfig: redisCfg}

	wm := NewWorkerManager(lb, wmCfg)

	executor := NewExecutor(wm)

	return &Scheduler{
		started:  atomic.Bool{},
		wm:       wm,
		executor: executor,
		store:    NewInMemStore(),
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
	err := s.store.AddTask(task)
	if err != nil {
		return nil, err
	}

	future, err := s.executor.Execute(task)
	if err != nil {
		return nil, err
	}

	return future, nil
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

func (s *Scheduler) GetTaskStatus(taskId string) (*TaskResult, error) {
	task, err := s.store.GetTask(taskId)
	if err != nil {
		return nil, err
	}

	workerTaskResult, err := s.executor.GetTaskStatus(task.GetWorkerTaskId(), task.GetWorkerId())
	if err != nil {
		return nil, err
	}

	// need to delete task from task store if task has reached its end state
	if workerTaskResult.TaskStatus == TaskStateDone ||
		workerTaskResult.TaskStatus == TaskStateFailure ||
		workerTaskResult.TaskStatus == TaskStateCancelled {
		s.store.DelTask(taskId)
	}

	return &TaskResult{
		TaskId: taskId,
		Status: workerTaskResult.TaskStatus,
		Data:   workerTaskResult.Data,
	}, nil
}
