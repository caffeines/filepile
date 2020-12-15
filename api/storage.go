package api

import (
	"log"
	"net/http"

	"github.com/caffeines/filepile/app"
	"github.com/caffeines/filepile/constants/errors"
	"github.com/caffeines/filepile/lib"
	"github.com/caffeines/filepile/middlewares"
	"github.com/caffeines/filepile/service"
	"github.com/labstack/echo/v4"
)

// RegisterStorageRoutes registers Bucketintacation routes
func RegisterStorageRoutes(endpoint *echo.Group) {
	endpoint.POST("/bucket/:name/", createBucket, middlewares.JWTAuth())
}

func createBucket(ctx echo.Context) error {
	resp := lib.Response{}
	name := ctx.Param("name")
	minioClient := app.GetMinioClient()
	minioService := service.NewMinioService()
	isExists, err := minioService.MakeBucket(minioClient, name)
	if err != nil {
		log.Println(err)
		if isExists {
			resp.Title = "Folder already exist"
			resp.Status = http.StatusConflict
			resp.Errors = lib.NewError(err.Error())
			resp.Code = errors.FolderAlreadyExist
			return resp.ServerJSON(ctx)
		}
		resp.Title = "Folder creation failed"
		resp.Status = http.StatusInternalServerError
		resp.Errors = err
		resp.Code = errors.SomethingWentWrong
		return resp.ServerJSON(ctx)
	}
	resp.Title = "Folder created successfully"
	return resp.ServerJSON(ctx)
}

func addFile(ctx echo.Context) error {
	resp := lib.Response{}

	return resp.ServerJSON(ctx)
}
