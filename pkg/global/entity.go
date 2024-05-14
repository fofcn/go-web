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

func SuccessWithData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, NewEntity("0", "ok", data))
}

func SuccessNoData(c *gin.Context) {
	c.JSON(http.StatusOK, NewEntity("0", "ok", nil))
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

func InteralServerErrorWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, NewEntity("500", msg, nil))
	c.Abort()
}

func InternalServerError(c *gin.Context, entity *HttpEntity) {
	c.JSON(http.StatusInternalServerError, entity)
	c.Abort()
}
