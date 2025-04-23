package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/glebsnigirev/final-GS/pkg/db"
)

// Функция для проверки повторений
func validateRepeat(repeat string) error {
	// Поддерживаем только следующие значения повторений
	if repeat != "" && repeat != "y" && repeat != "m" && repeat != "d" && !isValidRepeatFormat(repeat) {
		return fmt.Errorf("Неверное правило повторения: %v", repeat)
	}
	return nil
}

// Проверка формата повторений вида "d <число>"
func isValidRepeatFormat(repeat string) bool {
	re := regexp.MustCompile(`^d\s*\d+$`)
	return re.MatchString(repeat)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Некорректный JSON"})
		return
	}

	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	now := time.Now()
	if err := checkDate(&task, now); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func checkDate(task *db.Task, now time.Time) error {
	// Проверка повторений
	if err := validateRepeat(task.Repeat); err != nil {
		return err
	}

	// Если дата указана как "today", устанавливаем текущую дату
	if task.Date == "today" {
		task.Date = now.Format("20060102")
	}

	// Если дата пустая, устанавливаем текущую дату
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	// Проверка корректности формата даты
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("Неверный формат даты")
	}

	// Если дата в прошлом, обновляем ее на текущую или на следующую по правилу повторения
	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			// если t == now, можно оставить дату как есть
			if !t.Before(now.Truncate(24 * time.Hour)) {
				task.Date = t.Format("20060102")
			} else {
				next, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return fmt.Errorf("Неверное правило повторения: %w", err)
				}
				task.Date = next
			}
		}
	}

	return nil
}
