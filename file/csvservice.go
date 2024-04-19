package file

import "go-web/pkg/scheduler"

type CsvService interface {
	SplitCsv(filepath string) error
}

type csvservice struct {
	scheduler *scheduler.Scheduler
}

type CsvSplitTask struct {
	FilePath string `json:"filepath"`
	TaskType string `json:"tasktype"`
}

func NewCsvService() CsvService {
	return &csvservice{
		scheduler: scheduler.GetScheduler("rr"),
	}
}

func (c *csvservice) SplitCsv(filepath string) error {
	task := scheduler.NewTask(scheduler.TaskTypeCSVSplitter, CsvSplitTask{FilePath: filepath, TaskType: "1"})
	_, err := c.scheduler.Schedule(task)
	return err
}
