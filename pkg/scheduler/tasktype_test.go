package scheduler_test

import (
	"go-web/pkg/scheduler"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	SubTaskTypeNonExistent scheduler.SubTaskType = "non-existent"
)

func TestIsValidTask(t *testing.T) {
	tests := []struct {
		taskType      scheduler.TaskType
		subTaskType   scheduler.SubTaskType
		expectedValid bool
	}{
		{scheduler.TaskTypePdf, scheduler.SubTaskTypePdf2Csv, true},
		{scheduler.TaskTypeCsv, scheduler.SubTaskTypePdf2Csv, true},
		{scheduler.TaskTypePdf, scheduler.SubTaskTypePdf2Img, true},
		{scheduler.TaskTypePdf, scheduler.SubTaskTypePdfSplitter, true},
		{scheduler.TaskTypeCsv, scheduler.SubTaskTypePdf2Img, false},
		{scheduler.TaskTypePdf, SubTaskTypeNonExistent, false},
	}
	for _, test := range tests {
		actualValid := scheduler.IsValidTask(test.taskType, test.subTaskType)
		assert.Equal(t, test.expectedValid, actualValid)
	}
}
