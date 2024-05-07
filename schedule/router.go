package schedule

import (
	"go-web/pkg/global"
	"go-web/pkg/scheduler"

	"github.com/gin-gonic/gin"
)

func InitRouter(public *gin.RouterGroup) {
	indexRouter := NewScheduleRouter(NewScheduleService())
	public.POST("/schedule/worker", indexRouter.RegisterWorker)
	public.GET("/schedule/worker", indexRouter.GetWorkerList)
	public.DELETE("/schedule/worker/:id", indexRouter.DelWorker)
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
		global.RequestError(c, global.NewEntity("", "Invalid register request", nil))
		return
	}

	_ = sr.ss.RegisterWorker(scheduler.WorkerId(cmd.Id), cmd.Addr)

	global.SuccessNoData(c)
}

func (sr *ScheduleRouter) GetWorkerList(c *gin.Context) {
	dtos := sr.ss.GetWorkerList()
	global.SuccessWithData(c, dtos)
}

func (sr *ScheduleRouter) DelWorker(c *gin.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		global.RequestError(c, global.NewEntity("", "Invalid id", nil))
		return
	}
	_ = sr.ss.DeRegisterWorker(id)
	global.SuccessNoData(c)
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
	global.SuccessNoData(c)
}
