package api

import (
	"fmt"
	"net/http"
	"sprint_14/pkg/utils"
	"time"
)

// Обработчик для /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Если now не указан, используем текущую дату
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "Некорректный формат now", http.StatusBadRequest)
			return
		}
	}

	// Вычисляем следующую дату
	nextDate, err := utils.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Формируем ответ
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, nextDate)
}
