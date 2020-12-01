package api

import (
	"log"
	"net/http"

	"github.com/caffeines/sharehub/constants/errors"
	"github.com/caffeines/sharehub/lib"
	"github.com/caffeines/sharehub/validators"
	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(endpoint *echo.Group) {
	endpoint.POST("/", register)
}
func register(ctx echo.Context) error {
	_, err := validators.ValidateRegister(ctx)
	resp := lib.Response{}
	if err != nil {
		log.Println(err)
		resp.Title = "Invalid request data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.InvalidRegisterData
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	return nil
}
