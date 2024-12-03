package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"scheduler/internal/model"
	"scheduler/internal/storage"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request, storage *storage.Storage) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	var task model.Tasks
	err := storage.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
		Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, `{"error":"Ошибка выполнения запроса"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, `{"error":"Ошибка формирования ответа"}`, http.StatusInternalServerError)
		return
	}
}
