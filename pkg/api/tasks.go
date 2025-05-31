package api

import (
	"log"
	"net/http"
	"sprint_14/pkg/db"
	"time"
)

// Лимит задач
const TaskLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// Обработчик получения задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на получение задач")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		writeError(w, "Метод не поддерживается", http.StatusBadRequest)
		return
	}

	tasks, err := db.Tasks(TaskLimit)
	if err != nil {
		log.Printf("Ошибка получения задач: %v", err)
		writeError(w, err.Error(), http.StatusBadRequest)
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
	}, http.StatusOK)
}
