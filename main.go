package main

import (
	"github.com/caffeines/filepile/app"
	"github.com/caffeines/filepile/config"
	"github.com/caffeines/filepile/lib"
	"github.com/caffeines/filepile/models"
	"github.com/caffeines/filepile/server"
)

func init() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	err = app.ConnectMongo()
	if err != nil {
		panic(err)
	}
	err = app.SetMinioClient()
	if err != nil {
		panic(err)
	}
	lib.InitValidator()
	db := app.GetDB()
	err = models.InitializeIndex(db)
	if err != nil {
		panic(err)
	}
}

func main() {
	server.StartServer()
	defer func() {
		err := app.DisconnectMongo()
		if err != nil {
			panic(err)
		}
	}()

}
