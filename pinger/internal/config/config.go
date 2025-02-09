package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App AppConfig
}

type AppConfig struct {
	BackendURL   string
	PingInterval int
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл: ", err)
	}

	config := &Config{}

	config.App.BackendURL = getEnv("BACKEND_URL", "http://localhost:8080")
	config.App.PingInterval, err = strconv.Atoi(getEnv("PING_INTERVAL", "5"))
	if err != nil {
		return nil, fmt.Errorf("Недействительный PING_INTERVAL: %w", err)
	}

	return config, nil
}

// getEnv - читает переменные окружения(строковые значения) или возвращает значения по умолчанию
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt - читает переменные окружения(с целочисленными значениями) или возвращает значения по умолчанию
func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Недействительное значение для переменной окружения %s: %v. Использование значения по умолчанию %d", key, err, defaultValue)
		return defaultValue
	}
	return value
}
