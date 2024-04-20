package scheduler_test

import (
	"encoding/json"
	"go-web/pkg/scheduler"
	"testing"
)

type userdeftask struct {
	path string
	name string
}

func TestNewTask_ShouldSuccess_WhenGivenValidInput(t *testing.T) {
	task := scheduler.NewTask(scheduler.TaskTypeCSVSplitter, &userdeftask{
		path: "/test",
		name: "test",
	})
	bytes, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	println(string(bytes))
}
