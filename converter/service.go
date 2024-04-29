package converter

import (
	"go-web/pkg/config"
	"go-web/pkg/scheduler"
)

type ConverterService interface {
	CreateConvertTask(cmd *CreationConvertCmd) (*CreationConvertDto, error)
	GetTaskStatus(cmd *ConverterStatusCmd) (*ConverterStatusDto, error)
}

type converterServiceImpl struct {
	scheduler *scheduler.Scheduler
}

func NewConverterService() ConverterService {
	return &converterServiceImpl{
		scheduler: scheduler.GetScheduler(config.GetScheduler()),
	}
}

func (c *converterServiceImpl) CreateConvertTask(cmd *CreationConvertCmd) (*CreationConvertDto, error) {
	task := scheduler.NewTask(cmd.Type, cmd.SubType, cmd.FileId)
	taskFuture, err := c.scheduler.Schedule(task)
	if err != nil {
		return nil, err
	}
	taskResult, err := taskFuture.Get()
	if err != nil {
	}
}

func (c *converterServiceImpl) GetTaskStatus(cmd *ConverterStatusCmd) (*ConverterStatusDto, error) {
	taskResult, err := c.scheduler.GetTaskStatus(cmd.TaskId)
	if err != nil {
		return nil, err
	}
	return &ConverterStatusDto{
		TaskId: taskResult.TaskId,
		Status: taskResult.Status,
		Type:   taskResult.Type,
	}
}
