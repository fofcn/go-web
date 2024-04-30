package converter

import (
	"go-web/pkg/config"
	"go-web/pkg/scheduler"
	"time"
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
	taskUserDef := make(map[string]string)
	taskUserDef["file_id"] = cmd.FileId
	taskBuilder := scheduler.NewTaskBuilder()
	task := taskBuilder.
		SetType(scheduler.TaskType(cmd.Type)).
		SetSubType(scheduler.SubTaskType(cmd.SubType)).
		SetCreatedAt(time.Now()).
		SetUserDef(taskUserDef).
		Build()

	_, err := c.scheduler.Schedule(task)
	if err != nil {
		return nil, err
	}

	return &CreationConvertDto{
		TaskId: task.GetId(),
	}, nil

}

func (c *converterServiceImpl) GetTaskStatus(cmd *ConverterStatusCmd) (*ConverterStatusDto, error) {
	taskResult, err := c.scheduler.GetTaskStatus(cmd.TaskId)
	if err != nil {
		return nil, err
	}
	return &ConverterStatusDto{
		TaskId:  taskResult.TaskId,
		Status:  string(taskResult.Status),
		Type:    taskResult.Type,
		SubType: taskResult.SubType,
	}, nil
}
