package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var router = echo.New()

// GetRouter returns the api router
func GetRouter() http.Handler {
	router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 6,
		Skipper: func(ctx echo.Context) bool {
			return strings.Contains(ctx.Path(), "/fs/") || strings.Contains(ctx.Path(), "/download/")
		},
	}))

	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time: ${time_rfc3339}, method: ${method}, uri: ${uri}, status: ${status}\n",
	}))

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	router.Pre(middleware.AddTrailingSlash())
	router.Use(middleware.Recover())

	router.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"health": "OK"})
	})
	registerV1Routes()
	return router
}

func registerV1Routes() {
	v1 := router.Group("/v1")
	v1.GET("/", hello)
	v1.GET("/hello/", hello2)
	v1.GET("/hello/:id/", hello3)
}

func hello(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{"messsage": "hello"})
}
func hello2(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{"messsage": "hello2"})
}
func hello3(ctx echo.Context) error {
	id := ctx.Param("id")
	return ctx.JSON(http.StatusOK, map[string]string{"messsage": id})
}
