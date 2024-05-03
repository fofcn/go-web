package schedule

import (
	"go-web/pkg/scheduler"

	"github.com/gin-gonic/gin"
)

func InitRouter(public *gin.RouterGroup) {
	indexRouter := NewScheduleRouter(NewScheduleService())
	public.POST("/schedule", indexRouter.RegisterWorker)
	public.GET("/schedule", indexRouter.GetWorkerList)
	public.DELETE("/schedule/:id", indexRouter.DelWorker)
	public.PUT("/schedule/task/:id", indexRouter.UpdateTaskState)
}

type ScheduleRouter struct {
	ss ScheduleService
}

func NewScheduleRouter(scheduleservice ScheduleService) *ScheduleRouter {
	return &ScheduleRouter{
		ss: scheduleservice,
	}
}

func (sr *ScheduleRouter) RegisterWorker(c *gin.Context) {
	var cmd RegisterWorkerCmd
	err := c.ShouldBindJSON(&cmd)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid register request",
		})
		return
	}

	_ = sr.ss.RegisterWorker(scheduler.WorkerId(cmd.Id), cmd.Addr)

	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func (sr *ScheduleRouter) GetWorkerList(c *gin.Context) {
	dtos := sr.ss.GetWorkerList()
	c.JSON(200, gin.H{
		"workers": dtos,
	})
}

func (sr *ScheduleRouter) DelWorker(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.JSON(400, gin.H{
			"message": "Invalid id",
		})
		return
	}
	_ = sr.ss.DeRegisterWorker(id)
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func (sr *ScheduleRouter) UpdateTaskState(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.JSON(400, gin.H{
			"message": "Invalid id",
		})
		return
	}

	var cmd TaskUpdateCmd
	err := c.ShouldBindJSON(&cmd)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid cmd",
		})
		return
	}

	sr.ss.UpdateTaskState(id, cmd.TaskState)
}
