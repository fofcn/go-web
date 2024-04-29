package scheduler_test

import (
	"go-web/pkg/scheduler"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cfg = &scheduler.RedisConfig{
		ClusterMode: scheduler.RedisClusterModeStandalone,
		Addrs:       []string{":6379"},
	}
)

func TestAddWorker_ShouldAddedSuccess_WhenConnectRedisSuccess(t *testing.T) {
	store, err := scheduler.NewRedisWorkerStore(cfg)
	if err != nil {
		t.Fatal(err)
	}

	err = store.AddWorker(scheduler.NewWorker("1", "127.0.0.1:8080"))

	assert.True(t, err == nil)
}

func TestAddWorker_ShouldUpdateExpiration_WhenTheWorkerExistsed(t *testing.T) {
	store, err := scheduler.NewRedisWorkerStore(cfg)
	store.AddWorker(scheduler.NewWorker("1", "127.0.0.1:8080"))

	err = store.AddWorker(scheduler.NewWorker("1", "127.0.0.1:8080"))
	assert.True(t, err == nil)
}

func TestDelWorker_ShouldSuccess_WhenTheWorkerExistsed(t *testing.T) {
	store, err := scheduler.NewRedisWorkerStore(cfg)
	store.AddWorker(scheduler.NewWorker("1", "127.0.0.1:8080"))

	err = store.DelWorker("1")
	assert.True(t, err == nil)
}

func TestGetWokrerIds_ShouldReturnWorkerIdList_WhenWorkerExisting(t *testing.T) {
	store, err := scheduler.NewRedisWorkerStore(cfg)
	store.AddWorker(scheduler.NewWorker("1", "127.0.0.1:8080"))
	ids, err := store.GetWorkerIds()
	assert.True(t, err == nil)
	assert.True(t, len(ids) == 1)
}

func TestHeartbeat_ShouldSuccess_WhenWorkerExisting(t *testing.T) {
	store, err := scheduler.NewRedisWorkerStore(cfg)
	store.AddWorker(scheduler.NewWorker("1", "127.0.0.1:8080"))
	err = store.Heartbeat(scheduler.NewWorker("1", "127.0.0.1"))
	assert.True(t, err == nil)
}
