package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"scheduler/internal/config"
	"scheduler/internal/model"
	"scheduler/internal/scheduler"
	"strings"
	"time"
)

func ApiNextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	if nowStr == "" || dateStr == "" || repeatStr == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "invalid 'now' parameter", http.StatusBadRequest)
		return
	}

	nextDate, err := scheduler.NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, nextDate)
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.MustLoad()

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	today := now.Format("20060102")

	if strings.TrimSpace(task.Date) == "" {
		task.Date = today
	}

	taskDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, `{"error":"Дата указана в неправильном формате"}`, http.StatusBadRequest)
		return
	}

	if taskDate.Format("20060102") != today {
		if strings.TrimSpace(task.Repeat) == "" {
			task.Date = today
		} else {
			next_date, err := scheduler.NextDate(now, task.Date, task.Repeat)
			fmt.Println(next_date)
			if err != nil {
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
				return
			}
			task.Date = next_date
		}
	}

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		http.Error(w, `{"error":"Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"Ошибка записи в базу данных"}`, http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, `{"error":"Ошибка получения идентификатора задачи"}`, http.StatusInternalServerError)
		return
	}

	response := model.Response{ID: &id}
	json.NewEncoder(w).Encode(response)
}
