package schedule

import "go-web/pkg/task"

type ScheduleService interface {
	RegisterWorker(workerId task.WorkerId, addr string) error
	GetWorkerList() []*WorkerListDto
	DeRegisterWorker(id string) error
}

type scheduleimpl struct {
	scheduler task.Scheduler
}

func NewScheduleService() ScheduleService {
	scheduler := task.GetScheduler()
	return &scheduleimpl{
		scheduler: scheduler,
	}
}

func (s *scheduleimpl) RegisterWorker(workerId task.WorkerId, addr string) error {
	_ = s.scheduler.RegisterWorker(task.NewWorker(workerId, addr))
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
	_ = s.scheduler.DeRegisterWorker(task.WorkerId(id))
	return nil
}
