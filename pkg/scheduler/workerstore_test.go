package scheduler_test

import (
	"go-web/pkg/scheduler"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cfg = &scheduler.RedisConfig{
		Addrs: []string{":6379"},
	}
)

func TestAddWorker_ShouldAddedSuccess_WhenConnectRedisSuccess(t *testing.T) {
	store := scheduler.NewRedisWorkerStore(cfg)

	err := store.AddWorker(scheduler.NewWorker("1", "127.0.0.1:8080"))

	assert.True(t, err == nil)
}
