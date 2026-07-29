package api

import (
	"net/http"
	"os"
)

var appPassword string

func Init() {
	appPassword = os.Getenv("TODO_PASSWORD")

	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", signInHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
}
