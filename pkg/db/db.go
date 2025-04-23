package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(128) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// Инициализация базы данных
func Init(dbFile string) error {
	log.Println("Инициализация базы данных:", dbFile)

	// Проверяем  файла базы данных
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открытие подключения к базе данных
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}

	// Если база данных не существует, создаём схему
	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("ошибка создания схемы: %w", err)
		}
		log.Println("Схема базы данных успешно создана")
	}

	log.Println("База данных успешно инициализирована")
	return nil
}
