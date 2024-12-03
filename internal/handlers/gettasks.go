package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"scheduler/internal/config"
	"scheduler/internal/model"
	"strings"
	"time"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.MustLoad()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		http.Error(w, `{"error":"Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	search := r.URL.Query().Get("search")

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE 1=1"
	var args []interface{}

	if strings.TrimSpace(search) != "" {
		if searchDate, err := time.Parse(config.TimeFormat, search); err == nil {
			query += " AND date = ?"
			args = append(args, searchDate.Format(config.TimeFormat))
		} else {
			query += " AND (title LIKE ? OR comment LIKE ?)"
			searchTerm := "%" + search + "%"
			args = append(args, searchTerm, searchTerm)
		}
	}

	query += fmt.Sprintf(" ORDER BY date ASC LIMIT %d", config.StorageLimit)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, `{"error":"Ошибка выполнения запроса к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []model.Tasks{}
	for rows.Next() {
		var task model.Tasks
		if err := rows.Scan(&task.ID, &task.Task.Date, &task.Task.Title, &task.Task.Comment, &task.Task.Repeat); err != nil {
			http.Error(w, `{"error":"Ошибка обработки данных из базы"}`, http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, `{"error":"Ошибка чтения данных"}`, http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []model.Tasks{}
	}
	response := model.ResponseTasks{Tasks: tasks}
	json.NewEncoder(w).Encode(response)
}
