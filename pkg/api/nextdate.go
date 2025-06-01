package api

import (
	"fmt"
	"log"
	"net/http"
	"sprint_14/pkg/utils"
	"time"
)

// Обработчик для /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		http.Error(
			w,
			"Метод не поддерживается",
			http.StatusMethodNotAllowed,
		)
		return
	}
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
		now, err = time.Parse(utils.DateFormat, nowStr)
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

	// Записываем ответ с обработкой ошибки
	_, err = fmt.Fprint(w, nextDate)
	if err != nil {
		// Логируем ошибку, но не меняем статус ответа
		log.Printf("Ошибка записи ответа: %s", err)
	}
}
