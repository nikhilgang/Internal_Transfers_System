package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/InternalTransfer/internal/database"
)

const DefaultMaxTransferAmount int64 = 200000

type App struct {
	Env               string
	ServerPort        int
	DB                database.Config
	MaxTransferAmount int64
}

func Load() (App, error) {
	env := getEnv("APP_ENV", "development")

	port, err := getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return App{}, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	dbPort, err := getEnvInt("DB_PORT", 5432)
	if err != nil {
		return App{}, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	return App{
		Env:               env,
		ServerPort:        port,
		MaxTransferAmount: DefaultMaxTransferAmount,
		DB: database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "transaction_manager"),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	return strconv.Atoi(v)
}
