package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-web/pkg/http"
	"sync/atomic"
)

type WorkerStatus struct {
	IsHealthy      bool
	ActiveTasks    int
	CompletedTasks int
	FailedTasks    int
}

type WorkerId int
type Worker interface {
	GetId() WorkerId
	GetAddr() string
	Exec(task Task) (TaskFuture, error)
	Status() WorkerStatus
}

var (
	workerIdAssigner = atomic.Int32{}
)

type workimpl struct {
	addr         string
	id           WorkerId
	httpclient   http.HTTPClient
	taskapitable map[TaskType]string
}

func NewWorker(addr string) Worker {
	table := map[TaskType]string{
		TaskTypeB:           "/task/b",
		TaskTypeCSVSplitter: "/task",
	}
	return &workimpl{
		id:           WorkerId(workerIdAssigner.Add(1)),
		addr:         addr,
		httpclient:   http.NewCustomHTTPClient(),
		taskapitable: table,
	}
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

func (w *workimpl) GetAddr() string {
	return w.addr
}
