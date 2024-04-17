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
	return &workimpl{
		id:         WorkerId(workerIdAssigner.Add(1)),
		addr:       addr,
		httpclient: http.NewCustomHTTPClient(),
	}
}

func (w *workimpl) GetId() WorkerId {
	return w.id
}

func (w *workimpl) Exec(task Task) (TaskFuture, error) {
	taskapi, exists := w.taskapitable[task.GetType()]
	if exists {
		taskjson, err := json.Marshal(task)
		if err != nil {
			return nil, err
		}

		resp, status, err := w.httpclient.Post(taskapi, bytes.NewReader(taskjson), nil)
		if status == 200 {
			json.Unmarshal(resp, task)
			return NewTaskFuture(task), nil
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
