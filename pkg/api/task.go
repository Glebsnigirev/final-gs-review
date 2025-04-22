package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/glebsnigirev/final-GS/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			writeJson(w, map[string]string{"error": "Не указан идентификатор"})
			return
		}
		task, err := db.GetTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, task)

	case http.MethodPut:
		var task db.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}

		if task.Title == "" {
			writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
			return
		}

		now := time.Now()
		if err := checkDate(&task, now); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}

		err := db.UpdateTask(&task)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}

		writeJson(w, map[string]string{})

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			writeJson(w, map[string]string{"error": "Не указан идентификатор"})
			return
		}
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, map[string]string{})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
