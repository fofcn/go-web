package global

import "github.com/gin-gonic/gin"

type HttpEntity struct {
	ErrCode string `json:"err_code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
}

func NewEntity(code, msg string, data any) *HttpEntity {
	return &HttpEntity{
		ErrCode: code,
		Msg:     msg,
		Data:    data,
	}
}

func Success(c *gin.Context, entity *HttpEntity) {
	c.JSON(200, entity)
}

func Error(c *gin.Context, entity *HttpEntity) {
	c.JSON(500, entity)
	c.Abort()
}

func AuthError(c *gin.Context, entity *HttpEntity) {
	c.JSON(401, entity)
	c.Abort()
}

func RequestError(c *gin.Context, entity *HttpEntity) {
	c.JSON(400, entity)
	c.Abort()
}
