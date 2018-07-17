package main

import (
	"github.com/labstack/echo"
	"server/controller"
	"github.com/labstack/echo/middleware"
	"fmt"
	"server/config"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	controller.GithubController{}.Init(e)
	e.Start(fmt.Sprintf(":%d", config.Port))
}
