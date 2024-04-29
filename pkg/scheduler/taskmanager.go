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
	var err error
	if tmCfg.StoreType == "redis" {
		store, err = NewRedisTaskStore(&tmCfg.RedisConfig)
		if err != nil {
			panic(err)
		}
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

func (tm *TaskManager) GetTask(id string) (Task, error) {
	return tm.store.GetTask(id)
}

func (tm *TaskManager) DelTask(id string) error {
	return tm.store.DelTask(id)
}
