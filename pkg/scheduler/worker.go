package scheduler

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-web/pkg/http"
	"sync/atomic"
)

type WorkerStatus struct {
	IsHealthy      bool
	ActiveTasks    int `json:"active_tasks"`
	CompletedTasks int `json:"completed_tasks"`
	FailedTasks    int `json:"failed_tasks"`
	CancelledTasks int `json:"cancelled_tasks"`
}

func (ws WorkerStatus) String() string {
	return "WorkerStatus{" +
		", ActiveTasks:" + string(ws.ActiveTasks) +
		", CompletedTasks:" + string(ws.CompletedTasks) +
		", FailedTasks:" + string(ws.FailedTasks) +
		", CancelledTasks:" + string(ws.CancelledTasks) +
		"}"
}

type WorkerId string
type Worker interface {
	GetId() WorkerId
	GetAddr() string
	Exec(task Task) (TaskFuture, error)
	CheckStatus() error
	Status() WorkerStatus
	IncrErrorCounter() int32
	ResetErrorCounter()
}

type workimpl struct {
	addr         string
	id           WorkerId
	httpclient   http.HTTPClient
	taskapitable map[TaskType]string
	errCounter   atomic.Int32
	isHealty     atomic.Bool
}

func NewWorker(id WorkerId, addr string) Worker {
	table := map[TaskType]string{
		TaskTypeCSVSplitter: "/task",
	}

	worker := &workimpl{
		id:           id,
		addr:         addr,
		httpclient:   http.NewCustomHTTPClient(),
		taskapitable: table,
		isHealty:     atomic.Bool{},
	}
	worker.isHealty.Store(true)
	return worker
}

func (w *workimpl) GetId() WorkerId {
	return w.id
}

func (w *workimpl) Exec(task Task) (TaskFuture, error) {
	taskapi, exists := w.taskapitable[task.GetType()]
	if exists {
		taskjson, err := json.Marshal(task.GetUserDef())
		if err != nil {
			return nil, err
		}

		println(w.addr + taskapi)
		headers := map[string]string{
			"Content-Type": "application/json",
		}
		resp, status, err := w.httpclient.Post("http://"+w.addr+taskapi, bytes.NewReader(taskjson), headers)
		if status == 200 || err == nil {
			println(string(resp))
			json.Unmarshal(resp, task)
			return NewTaskFuture(task), nil
		} else {
			return nil, err
		}
	}
	return nil, errors.New("no api found")
}

func (w *workimpl) Status() WorkerStatus {
	return WorkerStatus{}
}

func (w *workimpl) IncrErrorCounter() int32 {
	return w.errCounter.Add(1)
}
func (w *workimpl) ResetErrorCounter() {
	w.errCounter.Store(0)
}

type statuscheckdto struct {
	IsHealthy      bool `json:"is_healthy"`
	ActiveTasks    int  `json:"active_tasks"`
	CompletedTasks int  `json:"completed_tasks"`
	FailedTasks    int  `json:"failed_tasks"`
}

func (w *workimpl) CheckStatus() error {
	resp, status, err := w.httpclient.Get("http://"+w.addr+"/executor/status", nil)
	if status == 200 && err == nil {
		statusdto := &statuscheckdto{}
		err := json.Unmarshal(resp, statusdto)
		if err != nil {
			return err
		}

		return nil
	}
	return errors.New("check status error")
}

func (w *workimpl) GetAddr() string {
	return w.addr
}
