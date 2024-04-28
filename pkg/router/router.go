package router

import (
	"fmt"
	"go-web/pkg/global"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CreateEngine() *gin.Engine {
	r := gin.Default()
	configRouter(r)
	return r
}

func configRouter(r *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}                                                               // 允许的源
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}                // 允许的方法
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "User-Agent"} // 允许的头部
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true // 是否允许发送cookie
	config.MaxAge = time.Hour      // 设置预请求（preflight）的缓存时间

	r.Use(cors.New(config))

	r.Use(gin.CustomRecovery(global.GErrorHandler))
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
