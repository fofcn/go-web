package task

type WorkerStatus struct {
	IsHealthy      bool
	ActiveTasks    int
	CompletedTasks int
	FailedTasks    int
}

type WorkerId int
type Worker interface {
	GetId() WorkerId
	Exec(task Task) (TaskFuture, error)
	Status() WorkerStatus
}

type workimpl struct {
}

func (w *workimpl) Exec(task Task) TaskFuture {
	return nil
}

func (w *workimpl) Status() WorkerStatus {
	return WorkerStatus{}
}
