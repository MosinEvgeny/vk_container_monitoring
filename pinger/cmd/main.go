package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"pinger/internal/config"
	"pinger/internal/usecase"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	pingUsecase := usecase.NewPingUsecase(cfg.App.BackendURL, time.Duration(cfg.App.PingInterval)*time.Second)

	// 3. Получение списка IP-адресов контейнеров (Заглушка - нужно реализовать динамическое обнаружение!)
	containerIPs := []string{"127.0.0.1", "8.8.8.8"} // Пример IP-адресов

	ticker := time.NewTicker(pingUsecase.GetPingInterval())
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigChan
		log.Println("Получен сигнал завершения. Завершаем работу...")
		cancel()
	}()

	for {
		select {
		case <-ticker.C:
			for _, ip := range containerIPs {
				pingTime := ping(ip)
				fmt.Printf("Ping %s, time: %dms\n", ip, pingTime)

				err := pingUsecase.SendPingData(ip, pingTime)
				if err != nil {
					log.Printf("Ошибка отправки данных для %s: %v", ip, err)
				}
			}
		case <-ctx.Done():
			log.Println("Завершение работы Pinger...")
			return
		}
	}
}

func ping(ip string) int {
	start := time.Now()
	conn, err := net.DialTimeout("ip4:icmp", ip, 2*time.Second)
	if err != nil {
		fmt.Printf("Ошибка пингования %s: %v\n", ip, err)
		return -1 // пинг не удался
	}
	defer conn.Close()

	duration := time.Since(start)
	return int(duration.Milliseconds())
}
