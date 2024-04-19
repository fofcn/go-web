package schedule

type RegisterWorkerCmd struct {
	Id   string `json:"id"`
	Addr string `json:"addr"`
}
