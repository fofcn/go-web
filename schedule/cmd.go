package schedule

type RegisterWorkerCmd struct {
	Id   string `json:"id"`
	Addr string `json:"addr"`
}

type TaskUpdateCmd struct {
	TaskState  string `json:"task_state"`
	UpdateTime string `json:"update_time"`
}
