package schedule

import (
	"go-web/pkg/config"
	"go-web/pkg/scheduler"
	"net/url"
	"strconv"
)

type ScheduleService interface {
	RegisterWorker(workerId scheduler.WorkerId, addr string) error
	GetWorkerList() []*WorkerListDto
	DeRegisterWorker(id string) error
	GetTaskStaus(taskid string) (*TaskResultDto, error)
	UpdateTaskState(taskid string, state string) error
}

type scheduleimpl struct {
	scheduler *scheduler.Scheduler
}

func NewScheduleService() ScheduleService {
	scheduler := scheduler.GetScheduler(config.GetScheduler())
	return &scheduleimpl{
		scheduler: scheduler,
	}
}

func (s *scheduleimpl) RegisterWorker(workerId scheduler.WorkerId, addr string) error {
	if !s.isValidUrl(addr) {
		return ErrInvalidWorkerAddress
	}
	_ = s.scheduler.RegisterWorker(scheduler.NewWorker(workerId, addr))
	return nil
}

func (s *scheduleimpl) isValidUrl(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	host := u.Hostname()
	port := u.Port()
	if host == "" || port == "" {
		return false
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return false
	}

	return p >= 0 && p <= 65535
}

func (s *scheduleimpl) GetWorkerList() []*WorkerListDto {
	workers := s.scheduler.GetWorkers()
	workerdtos := make([]*WorkerListDto, len(workers))
	for i, worker := range workers {
		workerdtos[i] = &WorkerListDto{
			Id:   string(worker.GetId()),
			Addr: worker.GetAddr(),
		}
	}

	return workerdtos
}

func (s *scheduleimpl) DeRegisterWorker(id string) error {
	_ = s.scheduler.DeRegisterWorker(scheduler.WorkerId(id))
	return nil
}

func (s *scheduleimpl) GetTaskStaus(taskId string) (*TaskResultDto, error) {
	taskResult, err := s.scheduler.GetTaskStatus(taskId)
	if err != nil {
		return nil, err
	}

	return &TaskResultDto{taskResult.TaskId, taskResult.Data}, nil
}

func (s *scheduleimpl) UpdateTaskState(taskId string, state string) error {
	return s.scheduler.UpdateTaskState(taskId, state)
}
