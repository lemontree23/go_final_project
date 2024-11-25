package handlers

import (
	"database/sql"
	"net/http"
	"scheduler/internal/config"
	"scheduler/internal/scheduler"
	"time"
)

func MarkTaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.MustLoad()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		http.Error(w, `{"error":"Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var repeat string
	var date string
	err = db.QueryRow("SELECT repeat, date FROM scheduler WHERE id = ?", id).Scan(&repeat, &date)
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	if repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
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

	_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
	if err != nil {
		http.Error(w, `{"error":"Ошибка при обновлении задачи"}`, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{}`))
}
