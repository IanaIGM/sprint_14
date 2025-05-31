package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// Ошибки для работы с БД
var (
	ErrNotFound = errors.New("запись не найдена")
	ErrConflict = errors.New("конфликт версий")
)

// Структура
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Функция добавления задачи
func AddTask(task *Task) (int64, error) {
	var id int64
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

// Функция отображения задач
func Tasks(limit int) ([]*Task, error) {
	rows, err := db.Query(`SELECT id, date, title, comment, repeat
	FROM scheduler
	ORDER BY date ASC
	LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка в SELECT %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка поиска ближайщих задач: %w", err)
		}
		tasks = append(tasks, &t)
	}
	//Обработка ошибки
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}

	if tasks == nil {
		tasks = make([]*Task, 0)
	}
	return tasks, nil
}

// Функция получения задач
func GetTask(id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("ошибка считывания задачи по ID")
	}

	t := &Task{}
	query := "SELECT * FROM scheduler WHERE id=?"
	err := db.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("ошибка при возврате строк")
	}
	return t, nil
}

// Функция редактирования задач
func UpdateTask(t *Task) error {
	query := `
UPDATE scheduler
SET date=?,
    title=?,
    comment=?,
    repeat=?
    WHERE id=?`
	res, err := db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("неправильный ID")
	}
	return nil

}

// Функция удаления задач
func DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("не указан идентификатор")
	}
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении задачи: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка подсчёта удалённых строк: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// Функция переноса даты
func UpdateDate(nextDate, id string) error {
	if id == "" {
		return fmt.Errorf("не указан идентификатор")
	}
	res, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, nextDate, id)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении даты: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка подсчёта обновлённых строк: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
