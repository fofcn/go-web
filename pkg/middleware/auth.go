package middleware

import (
	"fmt"
	"go-web/pkg/global"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func (c CustomClaims) Valid() error {
	if len(c.UserId) == 0 {
		return jwt.ErrInvalidKey
	}

	return nil
}

func InitMiddleware(r *gin.Engine) {
	r.Use(Auth())
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("Token")
		if len(tokenStr) == 0 || err != nil {
			global.AuthError(c, global.NewEntity("", "token is empty", nil))
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})
		if err != nil {
			global.AuthError(c, global.NewEntity("", "token is invalid", nil))
			return
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			global.AuthError(c, global.NewEntity("", "token is invalid", nil))
			return
		}
	}
}
