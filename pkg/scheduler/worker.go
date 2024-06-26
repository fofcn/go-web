package scheduler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-web/pkg/http"
	"sync/atomic"
	"time"
)

type WorkerStatus struct {
	IsHealthy      bool
	ActiveTasks    int `json:"active_tasks"`
	CompletedTasks int `json:"completed_tasks"`
	FailedTasks    int `json:"failed_tasks"`
	CancelledTasks int `json:"cancelled_tasks"`
}

func (ws WorkerStatus) String() string {
	return fmt.Sprintf("WorkerStatus{ , ActiveTasks: %d, CompletedTasks: %d , FailedTasks: %d , CancelledTasks: %d}", ws.ActiveTasks, ws.CompletedTasks, ws.FailedTasks, ws.CancelledTasks)
}

type WorkerId string
type Worker interface {
	GetId() WorkerId
	GetAddr() string
	Exec(task Task) (TaskFuture, error)
	CheckStatus() error
	GetTaskStatus(taskId string) (*WorkerTaskResult, error)
	GetLastHeartbeat() time.Time
	Heartbeat(time.Time)
	Status() WorkerStatus
	IncrErrorCounter() int32
	ResetErrorCounter()
	// not support now
	SetAbilities([]*WorkerAbility)
	// not support now
	GetAbilities() []*WorkerAbility
}

type WorkerAbility struct {
}

type workimpl struct {
	addr          string
	id            WorkerId
	httpclient    http.HTTPClient
	taskapitable  map[TaskType]string
	errCounter    atomic.Int32
	isHealty      atomic.Bool
	heartbeattime time.Time
	workerStatus  *WorkerStatus
}

func NewWorker(id WorkerId, addr string) Worker {
	table := map[TaskType]string{}

	worker := &workimpl{
		id:           id,
		addr:         addr,
		httpclient:   http.NewCustomHTTPClient(),
		taskapitable: table,
		isHealty:     atomic.Bool{},
		workerStatus: &WorkerStatus{},
	}
	worker.isHealty.Store(true)
	return worker
}

func (w *workimpl) GetId() WorkerId {
	return w.id
}

type TaskSubmitDto struct {
	TaskId string `json:"task_id"`
}

type taskdispath struct {
	TaskId      string      `json:"task_id"`
	TaskType    string      `json:"task_type"`
	TaskSubType string      `json:"task_sub_type"`
	TaskDetail  interface{} `json:"task_detail"`
}

func (w *workimpl) Exec(task Task) (TaskFuture, error) {
	task_dispath := &taskdispath{TaskId: task.GetId(),
		TaskType:    string(task.GetType()),
		TaskSubType: string(task.GetSubType()),
		TaskDetail:  task.GetUserDef(),
	}
	taskjson, err := json.Marshal(task_dispath)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	url := fmt.Sprintf("http://%s/task", w.addr)
	resp, status, err := w.httpclient.Post(url, bytes.NewReader(taskjson), headers)
	if status == 200 && err == nil {
		taskSubmitResponse := &TaskSubmitDto{}
		err := json.Unmarshal(resp, taskSubmitResponse)
		if err != nil {
			return nil, err
		}
		task.SetWorkerTaskId(taskSubmitResponse.TaskId)
		return NewTaskFuture(task), nil
	} else {
		return nil, errors.New("dispath task error")
	}
	return nil, errors.New("no api found")
}

func (w *workimpl) Status() WorkerStatus {
	return *w.workerStatus
}

func (w *workimpl) IncrErrorCounter() int32 {
	return w.errCounter.Add(1)
}
func (w *workimpl) ResetErrorCounter() {
	w.errCounter.Store(0)
}

type ExecutorStatusDto struct {
	Data StatusCheckDto `json:"data"`
}

type StatusCheckDto struct {
	IsHealthy      bool `json:"is_healthy"`
	ActiveTasks    int  `json:"active_tasks"`
	CompletedTasks int  `json:"completed_tasks"`
	FailedTasks    int  `json:"failed_tasks"`
	CancelledTasks int  `json:"cancelled_tasks"`
}

func (w *workimpl) CheckStatus() error {
	resp, status, err := w.httpclient.Get("http://"+w.addr+"/executor/status", nil)
	if status == 200 && err == nil {
		executorStatus := &ExecutorStatusDto{}
		err := json.Unmarshal(resp, executorStatus)
		if err != nil {
			return err
		}

		w.workerStatus.ActiveTasks = executorStatus.Data.ActiveTasks
		w.workerStatus.CompletedTasks = executorStatus.Data.CompletedTasks
		w.workerStatus.FailedTasks = executorStatus.Data.FailedTasks
		w.workerStatus.CancelledTasks = executorStatus.Data.CancelledTasks

		return nil
	}
	return errors.New("check status error")
}

func (w *workimpl) GetAddr() string {
	return w.addr
}

func (w *workimpl) GetTaskStatus(taskId string) (*WorkerTaskResult, error) {
	resp, status, err := w.httpclient.Get("http://"+w.addr+"/task/"+taskId, nil)
	if status == 200 && err == nil {
		taskStatus := &WorkerTaskResult{}
		err := json.Unmarshal(resp, taskStatus)
		if err != nil {
			return nil, err
		}

		return taskStatus, nil
	}
	return nil, errors.New("check task status error")
}

func (w *workimpl) Heartbeat(ht time.Time) {
	w.heartbeattime = ht
}

func (w *workimpl) GetLastHeartbeat() time.Time {
	return w.heartbeattime
}

// not support now
func (w *workimpl) SetAbilities([]*WorkerAbility) {

}

// not support now
func (w *workimpl) GetAbilities() []*WorkerAbility {
	return nil
}
