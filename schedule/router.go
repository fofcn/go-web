package schedule

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func InitRouter(public *gin.RouterGroup) {
	indexRouter := NewScheduleRouter(NewScheduleService())
	public.POST("/schedule", indexRouter.RegisterWorker)
	public.GET("/schedule", indexRouter.GetWorkerList)
	public.DELETE("/schedule/:id", indexRouter.DelWorker)
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

	_ = sr.ss.RegisterWorker(cmd.Addr)

	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func (sr *ScheduleRouter) GetWorkerList(c *gin.Context) {
	c.JSON(200, gin.H{
		"workers": sr.ss.GetWorkerList(),
	})
}

func (sr *ScheduleRouter) DelWorker(c *gin.Context) {
	id := c.Param("id")
	idi, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid id",
		})
		return
	}
	_ = sr.ss.DeRegisterWorker(idi)
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
