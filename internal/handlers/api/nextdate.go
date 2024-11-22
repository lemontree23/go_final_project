package api

import (
	"fmt"
	"net/http"
	"scheduler/internal/scheduler"
	"time"
)

func ApiNextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	// Проверка параметров
	if nowStr == "" || dateStr == "" || repeatStr == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	// Парсим текущую дату
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "invalid 'now' parameter", http.StatusBadRequest)
		return
	}

	// Вызываем NextDate
	nextDate, err := scheduler.NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, nextDate)
}
