package auth

import (
	"go-web/pkg/global"
	"go-web/pkg/middleware"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func InitRouter(public *gin.RouterGroup) {
	public.POST("/login", Login)
}

func Login(c *gin.Context) {
	_, exists := c.Get("claims")
	if exists {
		global.Success(c, global.NewEntity("0", "success", nil))
		return
	}
	token, err := genToken(strings.ReplaceAll(uuid.New().String(), "-", ""))
	if err != nil {
		global.Error(c, global.NewEntity("1", "error", nil))
		return
	}
	setCookieToken(c, token)

	global.Success(c, global.NewEntity("0", "success", nil))
}

func setCookieToken(c *gin.Context, token string) {
	_ = c.MustGet("AuthConfig")
	c.SetCookie("Token", token, 3600, "/", "localhost", false, true)
}

func genToken(userid string) (string, error) {
	claims := middleware.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 设置token过期时间
			Issuer:    "go-web",                                           // 设置token发行人
		},
		UserId: userid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret")) //  todo : secret
}
