package server

import (
	"log"
	"net/http"

	"go-TODO-list/pkg/api"
)

func Run(webDir, port string) error {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init()

	log.Printf("Starting server on port %s...", port)
	return http.ListenAndServe(":"+port, nil)
}
