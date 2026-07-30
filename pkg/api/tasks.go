package api

import (
	"net/http"
	"time"

	"go-TODO-list/pkg/db"
)

const searchDateFormat = "02.01.2006"

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.FormValue("search")

	var tasks []*db.Task
	var err error

	switch {
	case search == "":
		tasks, err = db.Tasks(50)
	default:
		if t, dateErr := time.Parse(searchDateFormat, search); dateErr == nil {
			tasks, err = db.TasksByDate(t.Format(dateFormat), 50)
		} else {
			tasks, err = db.TasksSearch(search, 50)
		}
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJson(w, http.StatusOK, TasksResp{Tasks: tasks})
}
