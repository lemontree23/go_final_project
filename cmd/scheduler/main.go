package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"net/http"
	"os"
	"scheduler/internal/config"
	"scheduler/internal/handlers"
	"scheduler/internal/storage"
)

func main() {
	//init config
	cfg := config.MustLoad()

	//logger
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	//database
	storage, err := storage.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to initialize storage:", err)
		os.Exit(1)
	}

	_ = storage

	//router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/*", http.FileServer(http.Dir(cfg.FileServer)))
	r.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	r.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.AddTaskHandler(w, r)
		case http.MethodGet:
			handlers.GetTaskHandler(w, r)
		case http.MethodPut:
			handlers.UpdateTaskHandler(w, r)
		case http.MethodDelete:
			handlers.DeleteTaskHandler(w, r)
		default:
			http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	})
	r.HandleFunc("/api/tasks", handlers.GetTasksHandler)
	r.HandleFunc("/api/task/done", handlers.MarkTaskDoneHandler)

	log.Info("scheduler started", slog.String("Port", cfg.Port), slog.String("Database", cfg.StoragePath))

	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Error("Failed to start server", slog.String("Error", err.Error()))
	}
}
