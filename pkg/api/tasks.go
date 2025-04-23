package api

import (
	"net/http"

	"github.com/glebsnigirev/final-GS/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}
	writeJson(w, TasksResp{Tasks: tasks})
}
