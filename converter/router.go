package converter

import "github.com/gin-gonic/gin"

func InitRouter(private *gin.RouterGroup) {
	cr := newConverterRouter()
	private.POST("/convert", cr.createConvertTask)
	private.GET("/convert", cr.getTaskStatus)
}

type ConverterRouter struct {
	converterService ConverterService
}

func newConverterRouter() *ConverterRouter {
	return &ConverterRouter{}
}

func (cr *ConverterRouter) createConvertTask(c *gin.Context) {

}

func (cr *ConverterRouter) getTaskStatus(c *gin.Context) {

}
