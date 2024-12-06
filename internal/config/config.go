package config

import (
	"cmp"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string `yaml:"port"`
	StoragePath string `yaml:"storagePath"`
	FileServer  string `yaml:"fileServer"`
}

func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Errorf("error loading .env file: %w", err)
	}

	var cfg Config

	cfg.Port = cmp.Or(os.Getenv("TODO_PORT"), "7540")
	cfg.StoragePath = cmp.Or(os.Getenv("TODO_DBFILE"), "./storage/scheduler.db")
	cfg.FileServer = cmp.Or(os.Getenv("TODO_FILESERVER"), "./web")

	return &cfg
}
