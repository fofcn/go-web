package middleware

import (
	"errors"
	"fmt"
	"go-web/pkg/config"
	"go-web/pkg/global"
	"net/http"

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
	r.Use(MustAuth())
}

func MustAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieCfg := c.MustGet("CookieCfg").(config.Cookie)
		jwtCfg := c.MustGet("JwtCfg").(config.Jwt)
		tokenStr, err := c.Cookie(cookieCfg.Name)
		if len(tokenStr) == 0 || err != nil {
			// global.AuthError(c, global.NewEntity("", "Unauthorized operation", nil))
			// return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtCfg.Secret), nil
		})
		if err != nil {
			// global.AuthError(c, global.NewEntity("", "Unauthorized operation", nil))
			// return
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			global.AuthError(c, global.NewEntity("", "token is invalid", nil))
			return
		}
	}
}

func OptionalToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieCfg := c.MustGet("CookieCfg").(config.Cookie)
		jwtCfg := c.MustGet("JwtCfg").(config.Jwt)
		tokenStr, err := c.Cookie(cookieCfg.Name)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			global.AuthError(c, global.NewEntity("", "token is empty", nil))
			return
		}

		if len(tokenStr) > 0 {
			token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtCfg.Secret), nil
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
}
