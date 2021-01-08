package server

import (
	"github.com/astgot/forum/internal/database"
)

// Config ...
type Config struct {
	WebPort  string
	Database *database.Config
}

// NewConfig generates configurations for the Server
func NewConfig() *Config {
	return &Config{
		WebPort:  ":8080",
		Database: database.NewConfig(),
	}
}
