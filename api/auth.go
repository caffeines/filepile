package api

import (
	"log"
	"net/http"

	"github.com/caffeines/filepile/app"
	"github.com/caffeines/filepile/constants/errors"
	"github.com/caffeines/filepile/data"
	"github.com/caffeines/filepile/lib"
	"github.com/caffeines/filepile/validators"
	"github.com/labstack/echo/v4"
)

// RegisterAuthRoutes registers authintacation routes
func RegisterAuthRoutes(endpoint *echo.Group) {
	endpoint.POST("/", register)
}
func register(ctx echo.Context) error {
	resp := lib.Response{}
	user, err := validators.ValidateRegister(ctx)
	hash, err := lib.HashPassword(user.Password)
	if err != nil {
		log.Println(err)
		resp.Title = "User registration failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	user.Password = hash
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
