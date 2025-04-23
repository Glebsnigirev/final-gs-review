package main

import (
	"log"
	"os"

	"github.com/glebsnigirev/final-GS/pkg/db"
	"github.com/glebsnigirev/final-GS/pkg/server"
)

func main() {
	// Получаем путь к файлу базы данных из переменной окружения или по умолчанию
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Инициализация БД
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("не удалось инициализировать БД: %v", err)
	}

	// Запуск сервера
	if err := server.Run(); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
}
