package controllers

import (
	"credibooktest/config"
	"credibooktest/controllers/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func acl(isAdmin string) echo.MiddlewareFunc {
	return config.ACL(isAdmin)
}

func Routes(e *echo.Echo) {
	e.POST("/login", handlers.Login)
	e.POST("/user", handlers.AddUser)

	users := e.Group("users", middleware.JWTWithConfig(config.JwtConfig))
	users.Use(acl("admin"))
	users.GET("", handlers.GetAllUsers)

	transaction := e.Group("transaction", middleware.JWTWithConfig(config.JwtConfig))
	transaction.GET("", handlers.GetAllTransaction)
	transaction.POST("", handlers.AddTransaction)
	transaction.PUT("/:id", handlers.UpdateTransaction)
	transaction.DELETE("/:id", handlers.DeleteTransaction)
}
