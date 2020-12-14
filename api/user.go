package api

import (
	"log"
	"net/http"

	"github.com/caffeines/filepile/app"
	"github.com/caffeines/filepile/constants"
	"github.com/caffeines/filepile/constants/errors"
	"github.com/caffeines/filepile/data"
	"github.com/caffeines/filepile/lib"
	"github.com/caffeines/filepile/middlewares"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterUserRoutes registers authintacation routes
func RegisterUserRoutes(endpoint *echo.Group) {
	endpoint.GET("/profile/", profile, middlewares.JWTAuth())
}

func profile(ctx echo.Context) error {
	resp := lib.Response{}
	db := app.GetDB()
	userRepo := data.NewUserRepo()

	userID := ctx.Get(constants.USER_ID).(primitive.ObjectID)
	user, err := userRepo.FindUserByID(db, userID)
	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "User not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.UserNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	resp.Data = user
	resp.Status = http.StatusOK
	return resp.ServerJSON(ctx)
}
