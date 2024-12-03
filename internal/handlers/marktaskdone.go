package handlers

import (
	"net/http"
	"scheduler/internal/scheduler"
	"scheduler/internal/storage"
	"time"
)

func MarkTaskDoneHandler(w http.ResponseWriter, r *http.Request, storage *storage.Storage) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	var repeat string
	var date string
	err := storage.QueryRow("SELECT repeat, date FROM scheduler WHERE id = ?", id).Scan(&repeat, &date)
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	if repeat == "" {
		_, err := storage.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{}`))
		return
	}

	nextDate, err := scheduler.NextDate(time.Now(), date, repeat)
	if err != nil {
		http.Error(w, `{"error":"Ошибка при вычислении следующей даты"}`, http.StatusInternalServerError)
		return
	}

	_, err = storage.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
	if err != nil {
		http.Error(w, `{"error":"Ошибка при обновлении задачи"}`, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{}`))
}
