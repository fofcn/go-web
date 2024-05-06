package converter

import (
	"go-web/pkg/global"

	"github.com/gin-gonic/gin"
)

func InitRouter(private *gin.RouterGroup) {
	cr := newConverterRouter()
	private.POST("/convert", cr.createConvertTask)
	private.GET("/convert/:id", cr.getTaskStatus)
}

type ConverterRouter struct {
	converterService ConverterService
}

func newConverterRouter() *ConverterRouter {
	return &ConverterRouter{
		converterService: NewConverterService(),
	}
}

func (cr *ConverterRouter) createConvertTask(c *gin.Context) {
	var cmd CreationConvertCmd
	err := c.ShouldBindJSON(&cmd)
	if err != nil {
		global.RequestError(c, global.NewEntity("convert task parameter error", err.Error(), nil))
		return
	}

	dto, err := cr.converterService.CreateConvertTask(&cmd)
	if err != nil {
		global.InternalServerError(c, global.NewEntity("create convert task error", err.Error(), nil))
		return
	}

	global.Success(c, global.NewEntity("success", "", dto))
}

func (cr *ConverterRouter) getTaskStatus(c *gin.Context) {
	taskId := c.Param("id")
	if len(taskId) == 0 {
		global.RequestError(c, global.NewEntity("task id is empty", "", nil))
		return
	}

	dto, err := cr.converterService.GetTaskStatus(&ConverterStatusCmd{
		TaskId: taskId,
	})
	if err != nil {
		global.InternalServerError(c, global.NewEntity("get task status error", err.Error(), nil))
		return
	}

	global.Success(c, global.NewEntity("success", "", dto))
}
