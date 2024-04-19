package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-web/pkg/http"
	"sync"
	"time"
)

type WorkerStatus struct {
	IsHealthy      bool
	ActiveTasks    int
	CompletedTasks int
	FailedTasks    int
}

type WorkerId string
type Worker interface {
	GetId() WorkerId
	GetAddr() string
	Exec(task Task) (TaskFuture, error)
	CheckStatus() error
	Status() WorkerStatus
}

type workimpl struct {
	addr         string
	id           WorkerId
	httpclient   http.HTTPClient
	taskapitable map[TaskType]string
}

func NewWorker(id WorkerId, addr string) Worker {
	table := map[TaskType]string{
		TaskTypeB:           "/task/b",
		TaskTypeCSVSplitter: "/task",
	}
	return &workimpl{
		id:           id,
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

type statuscheckdto struct {
	IsHealthy      bool `json:"is_healthy"`
	ActiveTasks    int  `json:"active_tasks"`
	CompletedTasks int  `json:"completed_tasks"`
	FailedTasks    int  `json:"failed_tasks"`
}

func (w *workimpl) CheckStatus() error {
	resp, status, err := w.httpclient.Get("http://"+w.addr+"/status", nil)
	if status == 200 || err == nil {
		statusdto := &statuscheckdto{}
		err := json.Unmarshal(resp, statusdto)
		if err != nil {
			return err
		}

		return nil
	}
	return err
}

func (w *workimpl) GetAddr() string {
	return w.addr
}

type WorkerManager struct {
	workers     *sync.Map
	lb          LoadBalancer
	timer       *time.Ticker
	healthTimer *time.Ticker
	done        chan bool
}

func NewWorkerManager(lbAlg string) *WorkerManager {
	lb, _ := NewLB(lbAlg)
	ww := &WorkerManager{
		lb:          lb,
		workers:     &sync.Map{},
		timer:       time.NewTicker(5 * time.Second),
		healthTimer: time.NewTicker(5 * time.Second),
		done:        make(chan bool),
	}

	go func() {
		ww.healthCheck()
		ww.evictWorker()
	}()

	return ww
}

func (ww *WorkerManager) AddWorker(worker Worker) {
	ww.workers.Store(worker.GetId(), worker)
}

func (ww *WorkerManager) DelWorker(id WorkerId) {
	ww.workers.Delete(id)
}

func (ww *WorkerManager) GetWorkers() []Worker {
	var workers []Worker
	ww.workers.Range(func(key, value any) bool {
		workers = append(workers, value.(Worker))
		return true
	})

	return workers
}

func (ww *WorkerManager) SelectWorker() Worker {
	return ww.lb.Select(ww.workers)
}

func (ww *WorkerManager) Close() {
	ww.timer.Stop()
	ww.healthTimer.Stop()
	ww.done <- true
}

func (ww *WorkerManager) healthCheck() {
	for {
		select {
		case <-ww.done:
			return
		case <-ww.healthTimer.C:
			ww.workers.Range(func(key, value any) bool {
				worker := value.(Worker)
				worker.CheckStatus()
				return true
			})
		}
	}
}

func (ww *WorkerManager) evictWorker() {
	for {
		select {
		case <-ww.done:
			return
		case <-ww.timer.C:
			ww.workers.Range(func(key, value any) bool {
				worker := value.(Worker)
				if worker.Status().IsHealthy == false {
					ww.DelWorker(worker.GetId())
				}
				return true
			})
		}
	}
}
