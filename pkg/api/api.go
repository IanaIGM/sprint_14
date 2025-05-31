package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sprint_14/pkg/db"
	"sprint_14/pkg/utils"
	"time"
)

// Инициализация API обработчиков
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
}

// Обработчик операций задач
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		id := r.URL.Query().Get("id")
		task, err := db.GetTask(id)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, db.ErrNotFound) {
				statusCode = http.StatusNotFound
			}
			writeError(w, err.Error(), statusCode)
			return
		}
		writeJSON(w, task, http.StatusOK)

	case http.MethodDelete:

		id := r.URL.Query().Get("id")
		if err := db.DeleteTask(id); err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, db.ErrNotFound) {
				statusCode = http.StatusNotFound
			}
			writeError(w, err.Error(), statusCode)
		} else {

			w.Write([]byte(`{}`))
		}

	case http.MethodPut:
		var t db.Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			writeError(w, "Неверный формат json", http.StatusBadRequest)
			return
		}

		//Проверяем ID и заголовок задачи
		if t.ID == "" {
			writeError(w, "Поле ID не может быть пустым", http.StatusBadRequest)
			return
		}
		if t.Title == "" {
			writeError(w, "Поле Title не может быть пустым", http.StatusBadRequest)
			return
		}
		//Проверка формата даты
		layout := "20060102"
		if _, err := time.Parse(layout, t.Date); err != nil {
			writeError(w, "Неверный формат даты", http.StatusBadRequest)
			return
		}
		//Проверка правила повторения задачи
		if t.Repeat != "" {
			if _, err := utils.NextDate(time.Now(), t.Date, t.Repeat); err != nil {
				writeError(w, "Неверное правило повторения", http.StatusBadRequest)
				return
			}
		}
		//Обновление даты задачи с учётом правила повторения
		parsed, _ := time.Parse(layout, t.Date)
		if t.Repeat != "" {
			nextDate, _ := utils.NextDate(time.Now(), t.Date, t.Repeat)
			if parsed.Before(time.Now()) {
				t.Date = nextDate
			}
		} else if parsed.Before(time.Now()) {
			t.Date = time.Now().Format(layout)
		}

		if err := db.UpdateTask(&t); err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, db.ErrConflict) {
				statusCode = http.StatusConflict
			}
			writeError(w, err.Error(), statusCode)
			return
		}
		w.Write([]byte(`{}`))

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Обработчик выполнения задачи
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	//Читаем id задачи
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	//Получение задачи из БД
	t, err := db.GetTask(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, db.ErrNotFound) {
			statusCode = http.StatusNotFound
		}
		writeError(w, err.Error(), statusCode)
		return
	}

	//Удаление если нет повторения
	if t.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, db.ErrNotFound) {
				statusCode = http.StatusNotFound
			}
			writeError(w, err.Error(), statusCode)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
		return
	}

	//Парсинг даты
	prev, err := time.Parse("20060102", t.Date)
	if err != nil {
		writeError(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}

	nextDate, err := utils.NextDate(prev, t.Date, t.Repeat)
	if err != nil {
		writeError(w, "Неверное правило повторения", http.StatusBadRequest)
		return
	}

	if err := db.UpdateDate(nextDate, id); err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, db.ErrConflict) {
			statusCode = http.StatusConflict
			writeError(w, err.Error(), statusCode)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  message,
		"status": http.StatusText(statusCode),
	})
}
