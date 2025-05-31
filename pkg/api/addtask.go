package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sprint_14/pkg/db"
	"sprint_14/pkg/utils"
	"time"
)

// Обработчик добавления задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Десериализация запроса
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Ошибка десериализации: %v", err)
		writeJSON(w, map[string]string{"error": "Ошибка десериализации JSON"}, http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}
	// Обработка "today"
	if task.Date == "today" {
		task.Date = time.Now().Format(utils.DateFormat)
	}

	// Проверка даты
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Добавление задачи
	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("Ошибка добавления задачи: %v", err)
		writeJSON(w, map[string]string{"error": "Ошибка добавления задачи"}, http.StatusInternalServerError)
		return
	}

	fmt.Printf("Полученный ID: %d\n", id)

	// Возврат ID
	response := map[string]string{
		"id": fmt.Sprintf("%d", id),
	}
	writeJSON(w, response, http.StatusCreated)
}

// Проверка и корректировка даты задачи
func checkDate(task *db.Task) error {
	now := time.Now()
	layout := "20060102"
	nowDate := now.Format(layout)

	// Если дата пустая, ставим текущую
	if task.Date == "" {
		task.Date = nowDate
	}

	// Проверяем формат даты
	t, err := time.Parse(layout, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты: %s", task.Date)
	}

	// Если дата в прошлом
	if t.Before(now) {
		task.Date = now.Format(layout)
	}

	// Проверяем правило повторения
	if task.Repeat != "" {
		_, err = utils.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("некорректное правило повторения: %s", task.Repeat)
		}
	}

	return nil
}
func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}
