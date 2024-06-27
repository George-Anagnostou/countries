package main

import (
	"github.com/George-Anagnostou/countries/internal/routes"
    "github.com/George-Anagnostou/countries/internal/middleware"
    "github.com/George-Anagnostou/countries/internal/templates"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
    e.Static("/static", "static")
	e.Use(echoMiddleware.Logger())
    e.Use(echoMiddleware.Recover())
    e.Use(middleware.InitSessionStore())
    e.Use(middleware.AuthMiddleware)
	e.Renderer = templates.NewTemplate()

    routes.RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":8000"))
}
