package scheduler

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type TaskState string

const (
	TaskStateCreated   = "PENDING"
	TaskStateRunning   = "RUNNING"
	TaskStateDone      = "DONE"
	TaskStateFailure   = "FAILURE"
	TaskStateCancelled = "CANCELLED" // not support for now
)

type TaskType string
type SubTaskType string

type TaskPriority int

const (
	TaskPriorityLow TaskPriority = iota
	TaskPriorityMedium
	TaskPriorityHigh
)

type Task interface {
	GetId() string
	GetState() TaskState
	GetType() TaskType
	SetSubType() SubTaskType
	GetSubType() SubTaskType
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
	TaskId string
	Status TaskState
	Data   interface{}
}

type WorkerTaskResult struct {
	TaskId     string      `json:"task_id"`
	TaskStatus TaskState   `json:"task_status"`
	Data       interface{} `json:"data"`
}

type taskimpl struct {
	id           string
	workerTaskId string
	workerId     WorkerId
	state        TaskState
	taskType     TaskType
	subType      SubTaskType
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

func NewTask(t TaskType, sub SubTaskType, userdef interface{}) Task {
	return &taskimpl{
		id:        strings.ReplaceAll(uuid.New().String(), "-", ""),
		state:     TaskStateCreated,
		taskType:  t,
		subType:   sub,
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

func (t *taskimpl) GetId() string {
	return t.id
}

func (t *taskimpl) GetState() TaskState {
	return t.state
}

func (t *taskimpl) GetType() TaskType {
	return t.taskType
}

func (t *taskimpl) SetSubType() SubTaskType {
	return t.subType
}

func (t *taskimpl) GetSubType() SubTaskType {
	return t.subType
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

type TaskBuilder struct {
	task *taskimpl
}

func (b *TaskBuilder) SetID(id string) *TaskBuilder {
	b.task.id = id
	return b
}

func (b *TaskBuilder) SetWorkerTaskId(workerTaskId string) *TaskBuilder {
	b.task.workerTaskId = workerTaskId
	return b
}

func (b *TaskBuilder) SetWorkerId(workerId WorkerId) *TaskBuilder {
	b.task.workerId = workerId
	return b
}

func (b *TaskBuilder) SetState(state TaskState) *TaskBuilder {
	b.task.state = state
	return b
}

func (b *TaskBuilder) SetType(taskType TaskType) *TaskBuilder {
	b.task.taskType = taskType
	return b
}

func (b *TaskBuilder) SetSubType(subType SubTaskType) *TaskBuilder {
	b.task.subType = subType
	return b
}

func (b *TaskBuilder) SetPriority(priority TaskPriority) *TaskBuilder {
	b.task.priority = priority
	return b
}

func (b *TaskBuilder) SetUserDef(userDef interface{}) *TaskBuilder {
	b.task.userdef = userDef
	return b
}

func (b *TaskBuilder) SetCreatedAt(t time.Time) {
	b.task.createdAt = t
}

func (b *TaskBuilder) Build() Task {
	return b.task
}

func NewTaskBuilder() *TaskBuilder {
	return &TaskBuilder{
		task: &taskimpl{},
	}
}
