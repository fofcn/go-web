package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type WorkerStore interface {
	AddWorker(worker Worker) error
	DelWorker(id WorkerId) error
	GetWorker(id WorkerId) (Worker, error)
	GetWorkerIds() ([]WorkerId, error)
	Ping(id WorkerId) error
}

type RedisConfig struct {
	ClusterMode string
	Addrs       []string
	ClientName  string
	Password    string
	DB          int
	MaxRetries  int
	Timeout     time.Duration

	PoolSize        int
	PoolTimeout     time.Duration
	MaxIdleConns    int
	MinIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

type RedisWorkerStore struct {
	redisConfig *RedisConfig
	client      redis.UniversalClient
}

type InMemWorkerStore struct {
	workers *sync.Map
}

const (
	WORKER_LIST_KEY = "ktools:worker:set"
	WORKER_INFO_KEY = "ktools:worker:info:"
)

func NewInMemWorkerStore() WorkerStore {
	return &InMemWorkerStore{
		workers: &sync.Map{},
	}
}

func NewRedisWorkerStore(cfg *RedisConfig) WorkerStore {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{cfg.Addrs[0]},
	})
	return &RedisWorkerStore{
		client: client,
	}
}

type RedisWorkerInfo struct {
	Id   WorkerId `json:"id"`
	Addr string   `json:"addr"`
}

func (rws *RedisWorkerStore) AddWorker(worker Worker) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wokerInfoKey := WORKER_INFO_KEY + string(worker.GetId())
	exists, err := rws.client.Exists(ctx, wokerInfoKey).Result()
	if err != nil {
		log.Fatalf("error on checking if worker info exists: %v", err)
		return err
	}
	if exists == 0 {
		err = rws.doAddWorker(ctx, worker, wokerInfoKey)
		if err != nil {
			log.Fatalf("error on adding worker info: %v", err)
			return err
		}
	} else {
		log.Printf("worker info already exists")
		setted, err := rws.client.Expire(ctx, wokerInfoKey, 300*time.Second).Result()
		if err != nil {
			log.Fatalf("error on setting worker info: %v", err)
			return err
		}

		if !setted {
			log.Printf("worker info already set")
			return errors.New("worker info already set")
		}
	}

	return nil
}

func (rws *RedisWorkerStore) doAddWorker(ctx context.Context, worker Worker, wokerInfoKey string) error {
	pushed, err := rws.client.RPush(ctx, WORKER_LIST_KEY, string(worker.GetId())).Result()
	if err != nil {
		log.Fatalf("error on pushing elements to the list: %v", err)
		return err
	}

	log.Printf("Number of elements after pushed:%v \n", pushed)

	workerInfo := &RedisWorkerInfo{
		Id:   worker.GetId(),
		Addr: worker.GetAddr(),
	}
	workerJson, err := json.Marshal(workerInfo)
	if err != nil {
		log.Fatalf("error on marshaling worker info: %v", err)
		return err
	}
	log.Printf("worker info: %v", string(workerJson))

	err = rws.client.Set(ctx, wokerInfoKey, workerJson, 300*time.Second).Err()
	if err != nil {
		log.Fatalf("error on setting worker info: %v", err)
		return err
	}

	return nil
}

func (rws *RedisWorkerStore) DelWorker(id WorkerId) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rws.client.LRem(ctx, WORKER_LIST_KEY, 0, id)
	rws.client.Del(ctx, WORKER_INFO_KEY+string(id))
	return nil
}

func (rws *RedisWorkerStore) GetWorker(id WorkerId) (Worker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stringCmd := rws.client.Get(ctx, WORKER_INFO_KEY+string(id))
	if stringCmd.Err() == nil {
		var workerInfo RedisWorkerInfo
		err := json.Unmarshal([]byte(stringCmd.Val()), &workerInfo)
		if err != nil {
			log.Fatalf("error on unmarshaling worker info: %v", err)
			return nil, err
		}
		return NewWorker(workerInfo.Id, workerInfo.Addr), nil
	}
	return nil, stringCmd.Err()
}

func (rws *RedisWorkerStore) GetWorkerIds() ([]WorkerId, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if rws.client.Exists(ctx, WORKER_LIST_KEY).Val() == 0 {
		return nil, nil
	} else if rws.client.Exists(ctx, WORKER_LIST_KEY).Val() == 1 {
		return []WorkerId{WorkerId(rws.client.LIndex(ctx, WORKER_LIST_KEY, 0).Val())}, nil
	}

	return nil, errors.New("no workers found")
}

func (rws *RedisWorkerStore) Ping(id WorkerId) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return rws.client.Expire(ctx, fmt.Sprintf("%s%s", WORKER_INFO_KEY, id), 300*time.Second).Err()
}

func (rws *InMemWorkerStore) AddWorker(worker Worker) error {

	return nil
}

func (rws *InMemWorkerStore) DelWorker(id WorkerId) error {

	return nil
}

func (rws *InMemWorkerStore) GetWorker(id WorkerId) (Worker, error) {
	return nil, nil
}

func (rws *InMemWorkerStore) GetWorkerIds() ([]WorkerId, error) {
	return nil, nil
}

func (rws *InMemWorkerStore) Ping(id WorkerId) error {
	return nil
}
