package main

import (
	"github.com/caffeines/sharehub/app"
	"github.com/caffeines/sharehub/config"
	"github.com/caffeines/sharehub/server"
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
