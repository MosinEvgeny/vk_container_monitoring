package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Port         string
	PingInterval int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл: ", err)
	}

	config := &Config{}

	config.App.Port = getEnv("APP_PORT", "8080")
	config.App.PingInterval, err = strconv.Atoi(getEnv("PING_INTERVAL", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid PING_INTERVAL: %w", err)
	}
	config.App.ReadTimeout = time.Duration(getEnvInt("READ_TIMEOUT", 5))
	config.App.WriteTimeout = time.Duration(getEnvInt("WRITE_TIMEOUT", 5))
	config.App.IdleTimeout = time.Duration(getEnvInt("IDLE_TIMEOUT", 60))

	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnv("DB_PORT", "5432")
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "postgres")
	config.Database.Name = getEnv("DB_NAME", "container_monitoring")

	return config, nil
}

func ConnectDB(cfg DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть подключение к базе данных: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("не удалось проверить связь с базой данных: %w", err)
	}

	log.Println("Успешное подключение к базе данных!")
	return db, nil
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
