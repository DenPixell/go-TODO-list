package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go-TODO-list/pkg/db"
)

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func writeError(w http.ResponseWriter, errText string) {
	writeJson(w, map[string]string{"error": errText})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		if afterNow(now, t) {
			task.Date = next
		}
	} else {
		if afterNow(now, t) {
			task.Date = now.Format(dateFormat)
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	if task.Title == "" {
		writeError(w, "не указан заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	writeJson(w, map[string]any{"id": strconv.FormatInt(id, 10)})
}
