package scheduler

import (
	"context"
	"encoding/json"
	"log"
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

const (
	WORKER_LIST_KEY = "ktools:worker:list"
	WORKER_INFO_KEY = "ktools:worker:info:"
)

func NewRedisWorkerStore(cfg *RedisConfig) *RedisWorkerStore {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{":6379"},
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
	ctx := context.Background()
	pushed, err := rws.client.RPush(ctx, WORKER_LIST_KEY, worker.GetId()).Result()
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
	wokerInfoKey := WORKER_INFO_KEY + string(worker.GetId())
	err = rws.client.Set(ctx, wokerInfoKey, workerJson, 300).Err()
	if err != nil {
		log.Fatalf("error on setting worker info: %v", err)
		return err
	}

	return nil
}

func (rws *RedisWorkerStore) DelWorker(id WorkerId) error {
	return nil
}

func (rws *RedisWorkerStore) GetWorker(id WorkerId) (Worker, error) {
	return nil, nil
}

func (rws *RedisWorkerStore) GetWorkerIds() ([]WorkerId, error) {
	return nil, nil
}
