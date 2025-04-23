package server

import (
	"net/http"
	"os"

	"github.com/glebsnigirev/final-GS/pkg/api"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	api.Init() // ← регистрируем API-обработчики

	http.Handle("/", http.FileServer(http.Dir("web")))
	return http.ListenAndServe(":"+port, nil)
}
