package scheduler

type TaskManagerConfig struct {
	StoreType   string
	RedisConfig RedisConfig
}

type TaskManager struct {
	store TaskStore
}

func NewTaskManager(tmCfg TaskManagerConfig) *TaskManager {
	var store TaskStore
	if tmCfg.StoreType == "redis" {
		store = NewRedisTaskStore(&tmCfg.RedisConfig)
	} else {
		store = NewInMemStore()
	}
	return &TaskManager{
		store: store,
	}
}

func (tm *TaskManager) AddTask(task Task) error {
	return tm.store.AddTask(task)
}

func (tm *TaskManager) GetTask(id int) (Task, error) {
	return tm.store.GetTask(id)
}

func (tm *TaskManager) DelTask(id int) error {
	return tm.store.DelTask(id)
}
