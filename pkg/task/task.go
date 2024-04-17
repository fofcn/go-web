package task

import (
	"sync/atomic"
	"time"
)

type TaskState int

const (
	TaskStateCreated TaskState = iota
	TaskStateRunning
	TaskStateDone
	TaskStateCancelled
	TaskStateFailure
)

type TaskType int

const (
	TaskTypeA TaskType = iota
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
}

type TaskFuture interface {
	GetTask() Task
	Wait(timeout time.Duration) error
	Done() bool
	Cancel() bool
}

type taskimpl struct {
	id        int
	state     TaskState
	taskType  TaskType
	createdAt time.Time
	priority  TaskPriority
	userdef   interface{}
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
