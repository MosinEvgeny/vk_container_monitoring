package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/repository"
	"backend/internal/usecase"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	db, err := config.ConnectDB(cfg.Database)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	defer db.Close()

	pingRepo := repository.NewPostgresPingRepository(db)

	pingUsecase := usecase.NewPingUsecase(pingRepo, time.Duration(cfg.App.PingInterval)*time.Second)

	pingHandler := handlers.NewPingHandler(pingUsecase)

	router := mux.NewRouter()
	pingHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  cfg.App.ReadTimeout * time.Second,
		WriteTimeout: cfg.App.WriteTimeout * time.Second,
		IdleTimeout:  cfg.App.IdleTimeout * time.Second,
	}

	go func() {
		log.Printf("Сервер запущен на порту %s", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Завершение сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении сервера: %v", err)
	}

	log.Println("Сервер успешно остановлен")
}
