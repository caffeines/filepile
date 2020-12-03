package app_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/caffeines/filepile/app"
	"github.com/caffeines/filepile/config"
)

func TestMain(m *testing.M) {
	config.LoadConfig()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCheckConnection(t *testing.T) {
	t.Run("it should connect mongoDB without error", func(t *testing.T) {
		connErr := app.ConnectMongo()
		if nil != connErr {
			t.Errorf("expecting error 'nil', but got %v", connErr)
		}
	})

	t.Run("it should return mongoDB client", func(t *testing.T) {
		client, err := app.GetMongoClient()
		if err != nil {
			t.Errorf("expecting error 'nil', but got %v", err)
		}
		clientType := fmt.Sprintf("%T", client)
		if client == nil {
			t.Errorf("expecting '*mongo.Client', but got %v", client)
		}
		if clientType != "*mongo.Client" {
			t.Errorf("expecting '*mongo.Client', but got %T", client)
		}
	})

	t.Run("it should return a mongoDB", func(t *testing.T) {
		db := app.GetDB()
		if db == nil {
			t.Errorf("expecting db, but got %v", db)
		}
	})

	t.Run("it should disconnect from mongoDB", func(t *testing.T) {
		err := app.DisconnectMongo()
		if err != nil {
			t.Errorf("expecting 'nil', but got %v", err)
		}
	})
}
