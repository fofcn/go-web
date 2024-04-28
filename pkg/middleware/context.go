package middleware

import (
	"go-web/pkg/config"

	"github.com/gin-gonic/gin"
)

func ConfigContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("CookieCfg", config.GetAuthConfig().Cookie)
		c.Set("JwtCfg", config.GetAuthConfig().Jwt)
	}
}
