package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	resp, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := w.Write(resp); err != nil {
		log.Println(err)
	}
}

func writeError(w http.ResponseWriter, status int, errText string) {
	writeJson(w, status, map[string]string{"error": errText})
}
