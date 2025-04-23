package api

import (
	"net/http"
)

func Init() {
	// Регистрируем только нужные обработчики
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
}
