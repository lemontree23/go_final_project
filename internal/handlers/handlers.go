package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"scheduler/internal/config"
	"scheduler/internal/model"
	"scheduler/internal/scheduler"
	"scheduler/internal/storage"
	"strings"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	if nowStr == "" || dateStr == "" || repeatStr == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	now, err := time.Parse(config.TimeFormat, nowStr)
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

func AddTaskHandler(w http.ResponseWriter, r *http.Request, storage *storage.Storage) {
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
	today := now.Format(config.TimeFormat)

	if strings.TrimSpace(task.Date) == "" {
		task.Date = today
	}

	taskDate, err := time.Parse(config.TimeFormat, task.Date)
	if err != nil {
		http.Error(w, `{"error":"Дата указана в неправильном формате"}`, http.StatusBadRequest)
		return
	}

	if taskDate.Format(config.TimeFormat) != today {
		if strings.TrimSpace(task.Repeat) == "" {
			task.Date = today
		} else {
			next_date, err := scheduler.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
				return
			}
			task.Date = next_date
		}
	}

	res, err := storage.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
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

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, storage *storage.Storage) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	result, err := storage.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	w.Write([]byte(`{}`))
}

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

	if _, err := time.Parse(config.TimeFormat, task.Date); err != nil {
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
