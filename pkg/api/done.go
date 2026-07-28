package api

import (
	"net/http"
	"time"

	"go-TODO-list/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err.Error())
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, err.Error())
			return
		}
		if err := db.UpdateDate(next, id); err != nil {
			writeError(w, err.Error())
			return
		}
	}

	writeJson(w, map[string]any{})
}
