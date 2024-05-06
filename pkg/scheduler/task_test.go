package scheduler_test

import (
	"encoding/json"
	"go-web/pkg/scheduler"
	"testing"
)

func TestNewTask_ShouldSuccess_WhenGivenValidInput(t *testing.T) {
	task := scheduler.NewTaskBuilder().Build()
	bytes, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	println(string(bytes))
}
