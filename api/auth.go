package api

import (
	"log"
	"net/http"
	"time"

	"github.com/caffeines/filepile/app"
	"github.com/caffeines/filepile/constants"
	"github.com/caffeines/filepile/constants/errors"
	"github.com/caffeines/filepile/data"
	"github.com/caffeines/filepile/lib"
	"github.com/caffeines/filepile/validators"
	"github.com/labstack/echo/v4"
)

// RegisterAuthRoutes registers authintacation routes
func RegisterAuthRoutes(endpoint *echo.Group) {
	endpoint.POST("/register/", register)
	endpoint.POST("/login/", login)
}

func login(ctx echo.Context) error {
	resp := lib.Response{}
	body, err := validators.ValidateLogin(ctx)
	if err != nil {
		log.Println(err)
		resp.Title = "Invalid login request data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.InvalidRegisterData
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	db := app.GetDB()
	userRepo := data.NewUserRepo()
	user, err := userRepo.FindUserByEmail(db, body.Email)
	if err != nil {
		log.Println(err)
		if lib.IsDocumentNotFoundError(err) {
			resp.Title = "User not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.UserNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}
		resp.Title = "User login failed"
		resp.Status = http.StatusUnauthorized
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	if ok := lib.CheckPasswordHash(body.Password, user.Password); !ok {
		resp.Title = "Email or password incorrect"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.InvalidLoginCredential
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	signedToken, err := lib.BuildJWTToken(user.Username, constants.USER_SCOPE, user.ID.Hex())
	if err != nil {
		log.Println(err)

		resp.Title = "Failed to sign auth token"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.UserLoginFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	result := map[string]interface{}{
		"access_token":  signedToken,
		"refresh_token": lib.NewRefresToken(),
		"expire_on":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"permission":    constants.USER_SCOPE,
	}
	resp.Status = http.StatusOK
	resp.Data = result
	return resp.ServerJSON(ctx)
}

func register(ctx echo.Context) error {
	resp := lib.Response{}
	user, err := validators.ValidateRegister(ctx)
	if err != nil {
		log.Println(err)
		resp.Title = "Invalid request data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.InvalidRegisterData
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
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
