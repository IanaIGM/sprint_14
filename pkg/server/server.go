package server

import (
	"fmt"
	"net/http"
	"os"
	"sprint_14/pkg/api"
)

// Запускаем веб-сервер
func Run() error {
	// Инициализируем API обработчики
	api.Init()

	// Определение порта для подключения
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Настраиваем файловый сервер
	http.Handle("/", http.FileServer(http.Dir("web")))

	// Запускаем сервер
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
