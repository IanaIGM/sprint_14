package main

import (
	"log"
	"sprint_14/pkg/db"
	"sprint_14/pkg/server"
)

func main() {
	// Инициализация БД
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	//Закрываем подключение
	defer db.GetDB().Close()
	// Запуск сервера
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
