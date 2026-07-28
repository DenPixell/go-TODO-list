package main

import (
	"log"
	"os"

	"go-TODO-list/pkg/db"
	"go-TODO-list/pkg/server"
)

const webDir = "./web"
const defaultPort = "7540"
const defaultDBFile = "scheduler.db"

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	if err := server.Run(webDir, port); err != nil {
		log.Fatal(err)
	}
}
