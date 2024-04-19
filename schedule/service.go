package schedule

import (
	"go-web/pkg/scheduler"
)

type ScheduleService interface {
	RegisterWorker(workerId scheduler.WorkerId, addr string) error
	GetWorkerList() []*WorkerListDto
	DeRegisterWorker(id string) error
}

type scheduleimpl struct {
	scheduler *scheduler.Scheduler
}

func NewScheduleService() ScheduleService {
	scheduler := scheduler.GetScheduler("rr")
	return &scheduleimpl{
		scheduler: scheduler,
	}
}

func (s *scheduleimpl) RegisterWorker(workerId scheduler.WorkerId, addr string) error {
	_ = s.scheduler.RegisterWorker(scheduler.NewWorker(workerId, addr))
	return nil
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
