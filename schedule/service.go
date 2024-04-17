package schedule

import "go-web/pkg/task"

type ScheduleService interface {
	RegisterWorker(addr string) error
	GetWorkerList() []*WorkerListDto
	DeRegisterWorker(id int) error
}

type scheduleimpl struct {
	scheduler task.Scheduler
}

func NewScheduleService() ScheduleService {
	scheduler, _ := task.NewScheduler()
	return &scheduleimpl{
		scheduler: scheduler,
	}
}

func (s *scheduleimpl) RegisterWorker(addr string) error {
	_ = s.scheduler.RegisterWorker(task.NewWorker(addr))
	return nil
}

func (s *scheduleimpl) GetWorkerList() []*WorkerListDto {
	workers := s.scheduler.GetWorkers()
	workerdtos := make([]*WorkerListDto, len(workers))
	for _, worker := range workers {
		workerdtos = append(workerdtos, &WorkerListDto{
			Id:   int(worker.GetId()),
			Addr: worker.GetAddr(),
		})
	}

	return workerdtos
}

func (s *scheduleimpl) DeRegisterWorker(id int) error {
	_ = s.scheduler.DeRegisterWorker(task.WorkerId(id))
	return nil
}
