package main

import (
	"fmt"

	"github.com/caffeines/sharehub/app"
	"github.com/caffeines/sharehub/config"
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
	fmt.Println("Hello Sharehub")
	defer func() {
		err := app.DisconnectMongo()
		if err != nil {
			panic(err)
		}
	}()

}
