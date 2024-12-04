package config

import (
	"cmp"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port        string `yaml:"port"`
	StoragePath string `yaml:"storagePath"`
	FileServer  string `yaml:"fileServer"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var cfg Config

	cfg.Port = cmp.Or(os.Getenv("TODO_PORT"), "7540")
	cfg.StoragePath = cmp.Or(os.Getenv("TODO_DBFILE"), "./storage/scheduler.db")
	cfg.FileServer = cmp.Or(os.Getenv("TODO_FILESERVER"), "./web")

	return &cfg
}
