package config

import (
	"github.com/spf13/viper"
)

// Database holds the database configuration
type Database struct {
	Name string
	URL  string
}

var db Database

// DB returns the default database configuration
func DB() Database {
	return db
}

// LoadDB loads database configuration
func LoadDB() {
	mu.Lock()
	defer mu.Unlock()

	db = Database{
		Name: viper.GetString("database.name"),
		URL:  viper.GetString("database.URL"),
	}
}
