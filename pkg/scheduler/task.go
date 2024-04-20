package scheduler

import (
	"sync/atomic"
	"time"
)

type TaskState int

const (
	TaskStateCreated TaskState = iota
	TaskStateRunning
	TaskStateDone
	TaskStateFailure
	TaskStateCancelled
)

type TaskType int

const (
	TaskTypeCSVSplitter TaskType = iota
	TaskTypeB
	TaskTypeC
)

type TaskPriority int

const (
	TaskPriorityLow TaskPriority = iota
	TaskPriorityMedium
	TaskPriorityHigh
)

type Task interface {
	GetId() int
	GetState() TaskState
	GetType() TaskType
	GetCreatedAt() time.Time
	GetPriority() TaskPriority
	GetUserDef() interface{}
	GetWorkerId() WorkerId
	SetWorkerId(workerId WorkerId)
	GetWorkerTaskId() string
	SetWorkerTaskId(string)
}

type TaskFuture interface {
	GetTask() Task
	Wait(timeout time.Duration) error
	Done() bool
	Cancel() bool
}

type TaskResult struct {
	TaskId int
	Status TaskState
	Data   interface{}
}

type WorkerTaskResult struct {
	TaskId string
	Status TaskState
	Data   interface{}
}

type taskimpl struct {
	id           int
	workerTaskId string
	workerId     WorkerId
	state        TaskState
	taskType     TaskType
	createdAt    time.Time
	priority     TaskPriority
	userdef      interface{}
}

type taskfutureimpl struct {
	task   Task
	done   bool
	cancel bool
	err    error
}

var (
	taskIdAssigner = atomic.Int32{}
)

func NewTask(t TaskType, userdef interface{}) Task {
	return &taskimpl{
		id:        int(taskIdAssigner.Add(1)),
		state:     TaskStateCreated,
		taskType:  t,
		createdAt: time.Now(),
		priority:  TaskPriorityLow,
		userdef:   userdef,
	}
}

func NewTaskFuture(task Task) TaskFuture {
	return &taskfutureimpl{
		task:   task,
		done:   false,
		cancel: false,
		err:    nil,
	}
}

func (t *taskimpl) GetWorkerTaskId() string {
	return t.workerTaskId
}

func (t *taskimpl) SetWorkerTaskId(id string) {
	t.workerTaskId = id
}

func (t *taskimpl) GetId() int {
	return t.id
}

func (t *taskimpl) GetState() TaskState {
	return t.state
}

func (t *taskimpl) GetType() TaskType {
	return t.taskType
}

func (t *taskimpl) GetCreatedAt() time.Time {
	return t.createdAt
}

func (t *taskimpl) GetPriority() TaskPriority {
	return t.priority
}

func (t *taskimpl) GetUserDef() interface{} {
	return t.userdef
}

func (t *taskimpl) SetWorkerId(workerId WorkerId) {
	t.workerId = workerId
}

func (t *taskimpl) GetWorkerId() WorkerId {
	return t.workerId
}

func (t *taskfutureimpl) GetTask() Task {
	return nil
}

func (t *taskfutureimpl) Wait(timeout time.Duration) error {
	return nil
}

func (t *taskfutureimpl) Done() bool {
	return false
}

func (t *taskfutureimpl) Cancel() bool {
	return false
}
