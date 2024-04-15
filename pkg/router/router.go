package router

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func errorHandler(c *gin.Context, err any) {
	// goErr := errors.Wrap(err, 2)
	c.AbortWithStatusJSON(500, "Hello Buddy")
}

func CreateEngine() *gin.Engine {
	r := gin.Default()
	configRouter(r)
	return r
}

func configRouter(r *gin.Engine) {
	r.Use(gin.CustomRecovery(errorHandler))
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
}
