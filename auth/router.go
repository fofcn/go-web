package auth

import (
	"go-web/pkg/config"
	"go-web/pkg/global"
	"go-web/pkg/middleware"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func InitRouter(public *gin.RouterGroup) {
	public.GET("/login", Login)
}

func Login(c *gin.Context) {
	_, exists := c.Get("claims")
	if exists {
		global.Success(c, global.NewEntity("0", "success", nil))
		return
	}
	token, err := genToken(c, strings.ReplaceAll(uuid.New().String(), "-", ""))
	if err != nil {
		global.Error(c, global.NewEntity("1", "error", nil))
		return
	}
	setCookieToken(c, token)

	global.Success(c, global.NewEntity("0", "success", nil))
}

func setCookieToken(c *gin.Context, token string) {
	cookieCfg := c.MustGet("CookieCfg").(config.Cookie)
	c.SetCookie(cookieCfg.Name, token, cookieCfg.MaxAge, cookieCfg.Path, cookieCfg.Domain, cookieCfg.Security, cookieCfg.HttpOnly)
}

func genToken(c *gin.Context, userId string) (string, error) {
	jwtCfg := c.MustGet("JwtCfg").(config.Jwt)

	claims := middleware.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtCfg.ExpiresIn) * time.Second)),
			Issuer:    jwtCfg.Issuer,
		},
		UserId: userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtCfg.Secret))
}
