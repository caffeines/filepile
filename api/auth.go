package api

import (
	"log"
	"net/http"

	"github.com/caffeines/sharehub/app"
	"github.com/caffeines/sharehub/constants/errors"
	"github.com/caffeines/sharehub/data"
	"github.com/caffeines/sharehub/lib"
	"github.com/caffeines/sharehub/validators"
	"github.com/labstack/echo/v4"
)

// RegisterAuthRoutes registers authintacation routes
func RegisterAuthRoutes(endpoint *echo.Group) {
	endpoint.POST("/", register)
}
func register(ctx echo.Context) error {
	user, err := validators.ValidateRegister(ctx)
	// TODO: hash password
	resp := lib.Response{}
	if err != nil {
		log.Println(err)
		resp.Title = "Invalid request data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.InvalidRegisterData
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	userRepo := data.NewUserRepo()
	db := app.GetDB()
	_, err = userRepo.Register(db, user)

	if err != nil {
		log.Println(err)
		resp.Title = "User registration failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	resp.Title = "User registration successful"
	resp.Status = http.StatusAccepted
	resp.Data = user

	return resp.ServerJSON(ctx)
}
