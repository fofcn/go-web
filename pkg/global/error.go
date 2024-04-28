package global

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomErr struct {
	Code string
	Msg  string
}

func (e *CustomErr) Error() string {
	return e.Msg
}

func GErrorHandler(c *gin.Context, err any) {
	customErr, ok := err.(*CustomErr)
	if ok {
		entity := &HttpEntity{
			ErrCode: customErr.Code,
			Msg:     customErr.Msg,
			Data:    nil,
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, entity)
		return
	}
}
