package router

import (
	"credibooktest/config"
	controller "credibooktest/controllers"

	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	config.SetMiddlewares(e)
	config.SetJwtMiddlewares(e)

	// set main routes
	controller.Routes(e)

	return e
}
