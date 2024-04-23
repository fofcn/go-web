package scheduler

import "fmt"

type TaskStore interface {
	AddTask(task Task) error
	GetTask(id int) (Task, error)
	DelTask(id int) error
}

type InMemStore struct {
	tasks map[int]Task
}

// 与InMemStore类似，只是存储在redis中，需要引入redis包
// 目前暂时不接入Redis,所以与InMemStore实现一样
type RedisStore struct {
	tasks map[int]Task
}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		tasks: make(map[int]Task),
	}
}

func (s *InMemStore) AddTask(task Task) error {
	s.tasks[task.GetId()] = task
	return nil
}

func (s *InMemStore) GetTask(id int) (Task, error) {
	if task, exists := s.tasks[id]; exists {
		return task, nil
	}

	return nil, fmt.Errorf("task not found, task id: %d", id)
}

func (s *InMemStore) DelTask(id int) error {
	delete(s.tasks, id)
	return nil
}
func NewRedisStore() *RedisStore {
	return &RedisStore{
		tasks: make(map[int]Task),
	}
}

func (s *RedisStore) AddTask(task Task) error {
	s.tasks[task.GetId()] = task
	return nil
}

func (s *RedisStore) GetTask(id int) (Task, error) {
	return s.tasks[id], nil
}

func (s *RedisStore) DelTask(id int) error {
	delete(s.tasks, id)
	return nil
}
