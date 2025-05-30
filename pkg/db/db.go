package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// Глобальная переменная для хранения подключения к БД
var db *sql.DB

// SQL-схема для создания таблицы и индекса
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

// Функция инициализации базы данных
func Init(dbFile string) error {
	// Проверка существования файла БД
	_, err := os.Stat(dbFile)
	install := errors.Is(err, os.ErrNotExist)

	// Открываем или создаем БД
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка в открытии SQL %v", err)
	}

	// Если БД новая, создаем схему
	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("ошибка в выполнении SQL-Запроса %v", err)
	}

	//Сообщение о создании таблицы
	if install {
		fmt.Println("Создана новая база данных:", dbFile)
	}

	//Удаляем из таблицы задачи без даты
	_, err = db.Exec(`DELETE FROM scheduler WHERE date = ''`)
	if err != nil {
		return fmt.Errorf("ошибка очистки базы: %v", err)

	}
	return nil
}

// Получение доступа к БД для других пакетов
func GetDB() *sql.DB {
	return db
}
