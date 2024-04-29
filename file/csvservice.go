package file

import (
	"go-web/pkg/config"
	"go-web/pkg/scheduler"
)

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
		scheduler: scheduler.GetScheduler(config.GetScheduler()),
	}
}

func (c *csvservice) SplitCsv(filepath string) error {
	return nil
}
