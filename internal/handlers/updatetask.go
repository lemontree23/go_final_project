package handlers

import (
	"encoding/json"
	"net/http"
	"scheduler/internal/model"
	"scheduler/internal/storage"
	"time"
)

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request, storage *storage.Storage) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task model.Tasks

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"Неверный формат данных"}`, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Заголовок обязателен"}`, http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		http.Error(w, `{"error":"Дата обязательна"}`, http.StatusBadRequest)
		return
	}

	if _, err := time.Parse("20060102", task.Date); err != nil {
		http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
		return
	}

	query := "UPDATE scheduler SET date = ?, title = ?"

	if task.Comment != "" {
		query += ", comment = ?"
	}
	if task.Repeat != "" {
		query += ", repeat = ?"
	}
	query += " WHERE id = ?"

	result, err := storage.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	w.Write([]byte(`{}`))
}
