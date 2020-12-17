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
	"github.com/caffeines/filepile/middlewares"
	"github.com/caffeines/filepile/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RegisterStorageRoutes registers Bucketintacation routes
func RegisterStorageRoutes(endpoint *echo.Group) {
	endpoint.POST("/bucket/:name/", createBucket, middlewares.JWTAuth())
}

func createBucket(ctx echo.Context) error {
	resp := lib.Response{}
	name := ctx.Param("name")
	db := app.GetDB()
	storageRepo := data.NewStorageRepo()
	userID := ctx.Get(constants.USER_ID).(primitive.ObjectID)
	bucket := &models.Bucket{
		ID:        primitive.NewObjectID(),
		Name:      name,
		CreatedBy: userID,
		CreatedAt: time.Now().UTC(),
		UpdateAt:  time.Now().UTC(),
	}

	isExists, err := storageRepo.CreateNewBucket(db, bucket)
	if isExists {
		resp.Title = "Folder already exist"
		resp.Status = http.StatusConflict
		resp.Errors = err
		resp.Code = errors.FolderAlreadyExist
		return resp.ServerJSON(ctx)
	}
	if err != nil {
		log.Println(err)
		resp.Title = "Folder creation failed"
		resp.Status = http.StatusInternalServerError
		resp.Errors = lib.NewError(err.Error())
		resp.Code = errors.SomethingWentWrong
		return resp.ServerJSON(ctx)
	}
	resp.Title = "Folder created successfully"
	resp.Data = bucket
	resp.Status = http.StatusOK
	return resp.ServerJSON(ctx)
}

func addFile(ctx echo.Context) error {
	resp := lib.Response{}
	return resp.ServerJSON(ctx)
}
