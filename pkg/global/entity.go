package global

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	c.JSON(http.StatusOK, entity)
}

func Error(c *gin.Context, entity *HttpEntity) {
	c.JSON(http.StatusInternalServerError, entity)
	c.Abort()
}

func AuthError(c *gin.Context, entity *HttpEntity) {
	c.JSON(http.StatusUnauthorized, entity)
	c.Abort()
}

func RequestError(c *gin.Context, entity *HttpEntity) {
	c.JSON(http.StatusBadRequest, entity)
	c.Abort()
}

func InternalServerError(c *gin.Context, entity *HttpEntity) {
	c.JSON(http.StatusInternalServerError, entity)
	c.Abort()
}
