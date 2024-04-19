package file

import "go-web/pkg/task"

type CsvService interface {
	SplitCsv(filepath string) error
}

type csvservice struct {
	scheduler task.Scheduler
}

type CsvSplitTask struct {
	FilePath string `json:"filepath"`
	TaskType string `json:"tasktype"`
}

func NewCsvService() CsvService {
	return &csvservice{
		scheduler: task.GetScheduler(),
	}
}

func (c *csvservice) SplitCsv(filepath string) error {
	task := task.NewTask(task.TaskTypeCSVSplitter, CsvSplitTask{FilePath: filepath, TaskType: "1"})
	_, err := c.scheduler.Schedule(task)
	return err
}
