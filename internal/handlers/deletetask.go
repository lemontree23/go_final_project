package handlers

import (
	"database/sql"
	"net/http"
	"scheduler/internal/config"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	result, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	w.Write([]byte(`{}`))
}
