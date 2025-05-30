package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sprint_14/pkg/db"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// Обработчик получения задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на получение задач")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		writeError(w, "Метод не поддерживается")
		return
	}

	tasks, err := db.Tasks(50)
	if err != nil {
		log.Printf("Ошибка получения задач: %v", err)
		writeError(w, err.Error())
		return
	}

	today := time.Now().Format("20060102")
	for i := range tasks {
		if tasks[i].Date == "" {
			tasks[i].Date = today
		}
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}

// Вывод ошибки
func writeError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// Отправить ответ
func writeJSON(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}
