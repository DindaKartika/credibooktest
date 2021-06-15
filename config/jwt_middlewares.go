package config

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

type JwtCustomClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
	IsAdmin bool `json:"is_admin"`
}

var JwtConfig middleware.JWTConfig

func SetJwtMiddlewares(e *echo.Echo) {
	JwtConfig = middleware.JWTConfig{
		Claims:     &JwtCustomClaims{},
		SigningKey: []byte(viper.GetString("jwtSign")),
	}
}

func ACL(isAdmin string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if isAdmin == "admin" {
				user := c.Get("user").(*jwt.Token)
				claims := user.Claims.(*JwtCustomClaims)

				if claims.IsAdmin {
					return next(c)
				}
			}
			return c.JSON(http.StatusForbidden, "You don't have permission to access this resource")
		}
	}
}
