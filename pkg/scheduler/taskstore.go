package scheduler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type TaskStore interface {
	AddTask(task Task) error
	GetTask(id string) (Task, error)
	DelTask(id string) error
}

type InMemStore struct {
	tasks map[string]Task
}

type RedisTaskStore struct {
	client redis.UniversalClient
}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		tasks: make(map[string]Task),
	}
}

func (s *InMemStore) AddTask(task Task) error {
	s.tasks[task.GetId()] = task
	return nil
}

func (s *InMemStore) GetTask(id string) (Task, error) {
	if task, exists := s.tasks[id]; exists {
		return task, nil
	}

	return nil, fmt.Errorf("task not found, task id: %d", id)
}

func (s *InMemStore) DelTask(id string) error {
	delete(s.tasks, id)
	return nil
}
func NewRedisTaskStore(redisCfg *RedisConfig) (*RedisTaskStore, error) {
	var client redis.UniversalClient
	if redisCfg.ClusterMode == RedisClusterModeStandalone {
		client = redis.NewClient(&redis.Options{
			Addr:     redisCfg.Addrs[0],
			Password: redisCfg.Password,
			DB:       redisCfg.DB,
		})
	} else if redisCfg.ClusterMode == RedisClusterModeCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisCfg.Addrs,
			Password: redisCfg.Password,
		})
	} else if redisCfg.ClusterMode == RedisClusterModeSentinel {
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    redisCfg.Addrs[0],
			SentinelAddrs: redisCfg.Addrs[1:],
			Password:      redisCfg.Password,
			DB:            redisCfg.DB,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &RedisTaskStore{
		client: client,
	}, nil
}

type TaskRedisDto struct {
	State     string      `json:"task_state"`
	Type      string      `json:"type"`
	SubType   string      `json:"sub_type"`
	Priority  int         `json:"priority"` // not supported for now
	WorkerId  string      `json:"worker_id"`
	CreatedAt time.Time   `json:"created_at"`
	UserDef   interface{} `json:"user_def"`
}

func (s *RedisTaskStore) AddTask(task Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	taskDto := &TaskRedisDto{
		State:     string(task.GetState()),
		Type:      string(task.GetType()),
		SubType:   string(task.GetSubType()),
		Priority:  int(task.GetPriority()),
		WorkerId:  string(task.GetWorkerId()),
		UserDef:   task.GetUserDef(),
		CreatedAt: task.GetCreatedAt(),
	}

	pushed, err := s.client.HSet(ctx, task.GetId(), taskDto).Result()
	if err != nil || pushed == 0 {
		return err
	}

	return nil
}

func (s *RedisTaskStore) GetTask(id string) (Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := s.client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	taskDetail, err := s.client.HGetAll(ctx, id).Result()
	if err != nil {
		return nil, err
	}

	priority, err := strconv.Atoi(taskDetail["priority"])
	if err != nil {
		return nil, err
	}
	taskBuilder := NewTaskBuilder()
	task := taskBuilder.
		SetID(id).
		SetWorkerTaskId(taskDetail["worker_task_id"]).
		SetState(TaskState(taskDetail["task_state"])).
		SetType(TaskType(taskDetail["type"])).
		SetSubType(SubTaskType(taskDetail["sub_type"])).
		SetPriority(TaskPriority(priority)).
		SetWorkerId(WorkerId(taskDetail["worker_id"])).
		SetUserDef(taskDetail["user_def"]).
		Build()

	return task, nil
}

func (s *RedisTaskStore) DelTask(id string) error {
	return nil
}
