package schedule

type WorkerListDto struct {
	Id   string `json:"id"`
	Addr string `json:"addr"`
}

type TaskResultDto struct {
	TaskId string
	Data   interface{}
}
