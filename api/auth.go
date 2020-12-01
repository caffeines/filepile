package api

import (
	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(endpoint *echo.Group) {
	endpoint.POST("/", register)
}
func register(ctx echo.Context) error {

	return nil
}
