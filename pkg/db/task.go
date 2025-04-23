package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Структура задачи
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int, search string) ([]*Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler"
	args := []interface{}{}
	where := ""

	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			where = " WHERE date = ?"
			args = append(args, t.Format("20060102"))
		} else {
			where = " WHERE title LIKE ? OR comment LIKE ?"
			s := "%" + search + "%"
			args = append(args, s, s)
		}
	}

	query = query + where + " ORDER BY date LIMIT ?"
	args = append(args, limit)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	var t Task
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Задача не найдена")
		}
		return nil, err
	}
	return &t, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	_, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", next, id)
	return err
}
