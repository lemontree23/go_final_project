package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"net/http"
	"os"
	"scheduler/internal/config"
	"scheduler/internal/handlers/api"
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
		log.Error("Failed to initialize storage")
		os.Exit(1)
	}

	_ = storage

	//router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/*", http.FileServer(http.Dir(cfg.FileServer)))
	r.HandleFunc("/api/nextdate", api.ApiNextDateHandler)

	log.Info("scheduler started", slog.String("Port", cfg.Port), slog.String("Database", cfg.StoragePath))

	http.ListenAndServe(":"+cfg.Port, r)
}
