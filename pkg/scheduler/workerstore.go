package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	RedisClusterModeStandalone = "standalone"
	RedisClusterModeCluster    = "cluster"
	RedisClusterModeSentinel   = "sentinel"
	RedisClusterModeFailover   = "failover"
)

type WorkerStore interface {
	AddWorker(worker Worker) error
	DelWorker(id WorkerId) error
	GetWorker(id WorkerId) (Worker, error)
	GetWorkerIds() ([]WorkerId, error)
	Heartbeat(worker Worker) error
	Close() error
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

func NewRedisWorkerStore(cfg *RedisConfig) (WorkerStore, error) {
	var client redis.UniversalClient
	if cfg.ClusterMode == RedisClusterModeStandalone {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Addrs[0],
			Password: cfg.Password,
			DB:       cfg.DB,
		})
	} else if cfg.ClusterMode == RedisClusterModeCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    cfg.Addrs,
			Password: cfg.Password,
		})
	} else if cfg.ClusterMode == RedisClusterModeSentinel {
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    cfg.Addrs[0],
			SentinelAddrs: cfg.Addrs[1:],
			Password:      cfg.Password,
			DB:            cfg.DB,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisWorkerStore{
		client: client,
	}, nil
}

type RedisWorkerInfo struct {
	Id           WorkerId  `json:"id"`
	Addr         string    `json:"addr"`
	LastPingTime time.Time `json:"last_ping_time"`
	// not supported now
	Abilities []string
}

func (rws *RedisWorkerStore) AddWorker(worker Worker) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rws.client.Ping(ctx).Result()
	if err != nil {
		log.Printf("error on checking redis connection: %v", err)
		return err
	}

	wokerInfoKey := WORKER_INFO_KEY + string(worker.GetId())
	exists, err := rws.client.Exists(ctx, wokerInfoKey).Result()
	if err != nil {
		log.Printf("error on checking if worker info exists: %v", err)
		return err
	}
	if exists == 0 {
		err = rws.doAddWorker(ctx, worker, wokerInfoKey)
		if err != nil {
			log.Printf("error on adding worker info: %v", err)
			return err
		}
	} else {
		log.Printf("worker info already exists， %v", worker.GetId())
		setted, err := rws.client.Expire(ctx, wokerInfoKey, 300*time.Second).Result()
		if err != nil {
			log.Printf("error on setting worker info: %v", err)
			return err
		}

		if !setted {
			log.Printf("worker info already set")
			return errors.New("worker info already set")
		}
	}

	return nil
}

func (rws *RedisWorkerStore) hashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func (rws *RedisWorkerStore) doAddWorker(ctx context.Context, worker Worker, wokerInfoKey string) error {
	var idHash uint64 = rws.hashString(string(worker.GetId()))
	z := redis.Z{
		Score:  float64(idHash),
		Member: string(worker.GetId()),
	}
	pushed, err := rws.client.ZAdd(ctx, WORKER_LIST_KEY, z).Result()
	if err != nil {
		log.Printf("error on pushing elements to the list: %v", err)
		return err
	}

	log.Printf("Number of elements after pushed:%v \n", pushed)

	workerInfo := &RedisWorkerInfo{
		Id:           worker.GetId(),
		Addr:         worker.GetAddr(),
		LastPingTime: time.Now(),
	}
	workerJson, err := json.Marshal(workerInfo)
	if err != nil {
		log.Printf("error on marshaling worker info: %v", err)
		return err
	}
	log.Printf("worker info: %v", string(workerJson))

	err = rws.client.Set(ctx, wokerInfoKey, workerJson, 300*time.Second).Err()
	if err != nil {
		log.Printf("error on setting worker info: %v", err)
		return err
	}

	return nil
}

func (rws *RedisWorkerStore) DelWorker(id WorkerId) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rws.client.ZRem(ctx, WORKER_LIST_KEY, id)
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
			log.Printf("error on unmarshaling worker info: %v", err)
			return nil, err
		}
		return NewWorker(workerInfo.Id, workerInfo.Addr), nil
	}
	return nil, stringCmd.Err()
}

func (rws *RedisWorkerStore) GetWorkerIds() ([]WorkerId, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := rws.client.ZRange(ctx, WORKER_LIST_KEY, 0, -1)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	workerIds, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	var retWorkerIds = make([]WorkerId, len(workerIds))
	for i, workerId := range workerIds {
		retWorkerIds[i] = WorkerId(workerId)
	}

	return retWorkerIds, nil
}

func (rws *RedisWorkerStore) Heartbeat(worker Worker) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	workerKey := fmt.Sprintf("%s%s", WORKER_INFO_KEY, string(worker.GetId()))
	existed, err := rws.client.Exists(ctx, workerKey).Result()
	if err != nil {
		return errors.New("worker cannot be found")
	}

	if existed == 0 {
		err = rws.doAddWorker(ctx, worker, workerKey)
		if err != nil {
			return err
		}
	}

	return rws.client.Expire(ctx, workerKey, 300*time.Second).Err()
}

func (rws *RedisWorkerStore) Close() error {
	return rws.client.Close()
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

func (rws *InMemWorkerStore) Heartbeat(worker Worker) error {
	return nil
}

func (rws *InMemWorkerStore) Close() error {
	return nil
}
