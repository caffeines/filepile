package api

import (
	"bytes"
	"fmt"
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
	endpoint.GET("/my-bucket/", myBuckets, middlewares.JWTAuth())
	endpoint.POST("/upload-file/:bucket/", uploadFile, middlewares.JWTAuth())
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

func myBuckets(ctx echo.Context) error {
	resp := lib.Response{}
	db := app.GetDB()
	storageRepo := data.NewStorageRepo()
	userID := ctx.Get(constants.USER_ID).(primitive.ObjectID)
	lastID := ctx.QueryParam("lastId")
	log.Println(lastID)
	buckets, err := storageRepo.FindBucketsByCreatorID(db, userID, lastID)
	if err != nil {
		log.Println(err)
		resp.Title = "Can not fetch folders"
		resp.Status = http.StatusInternalServerError
		resp.Errors = lib.NewError(err.Error())
		resp.Code = errors.DatabaseQueryFailed
		return resp.ServerJSON(ctx)
	}
	resp.Data = buckets
	resp.Status = http.StatusOK
	return resp.ServerJSON(ctx)
}

func uploadFile(ctx echo.Context) error {
	resp := lib.Response{}
	bucket := ctx.Param("bucket")
	if err := ctx.Request().ParseMultipartForm(32 << 22); err != nil {
		resp.Title = "Couldn't parse multipart form"
		resp.Status = http.StatusBadRequest
		resp.Code = errors.InvalidMultiPartBody
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	r := ctx.Request()
	r.Body = http.MaxBytesReader(ctx.Response(), r.Body, 32<<21) // 32 Mb
	f, h, e := r.FormFile("file")

	if e != nil {
		resp.Title = "No multipart file"
		resp.Status = http.StatusBadRequest
		resp.Code = errors.InvalidMultiPartBody
		resp.Errors = e
		return resp.ServerJSON(ctx)
	}

	body := make([]byte, h.Size)
	_, errR := f.Read(body)
	if errR != nil {
		resp.Title = "Unable to read multipart data"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.SomethingWentWrong
		resp.Errors = errR
		return resp.ServerJSON(ctx)
	}
	fileRepo := data.NewFileRepo()
	contentType := h.Header.Get("Content-Type")
	minioClient := app.GetMinioClient()
	errU := fileRepo.UploadToMinio(bucket, h.Filename, contentType, bytes.NewReader(body), h.Size, minioClient)
	if errU != nil {
		resp.Title = "Minio service failed"
		resp.Status = http.StatusInternalServerError
		// resp.Code = errors.
		resp.Errors = errU
		return resp.ServerJSON(ctx)
	}
	newFileNameWithBucket := fmt.Sprintf("%s/%s", bucket, h.Filename)
	resp.Status = http.StatusCreated
	resp.Data = map[string]interface{}{
		"path": newFileNameWithBucket,
	}
	return resp.ServerJSON(ctx)
}
