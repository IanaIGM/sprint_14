package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Константа для формата даты
const DateFormat = "20060102"

// Функция для проверки, что дата больше текущей
func AfterNow(date, now time.Time) bool {
	return date.After(now) || date.Equal(now)
}

// Функция для проверки корректности интервала дней
func validInterval(interval int) bool {
	return interval >= 1 && interval <= 400
}

// Основная функция вычисления следующей даты
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Проверка корректности входных данных
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	// Парсинг начальной даты
	start, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата начала: %w", err)
	}

	// Разделение строки на части
	parts := strings.Split(repeat, " ")

	// Обработка разных типов правил
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("некорректный формат правила d")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || !validInterval(interval) {
			return "", errors.New("некорректный интервал дней")
		}

		// Вычисление следующей даты
		for {
			start = start.AddDate(0, 0, interval)
			if AfterNow(start, now) {
				break
			}
		}

	case "y":
		// Ежегодное повторение
		for {
			start = start.AddDate(1, 0, 0)
			if AfterNow(start, now) {
				break
			}
		}

	default:
		return "", errors.New("неподдерживаемый формат правила")
	}

	return start.Format(DateFormat), nil
}
