package main

import (
	"log"

	"github.com/astgot/forum/internal/server"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	config := server.NewConfig() // generating config for server
	server := server.New(config) // creating new instance based on the 'config'
	// Starting the Server
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
