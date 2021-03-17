package main

import (
	"context"
	"errors"
	"github.com/koind/cacher/internal/config"
	"github.com/koind/cacher/internal/domain/service"
	"github.com/koind/cacher/internal/handler"
	"github.com/koind/cacher/internal/storage/memory"
	"github.com/spf13/pflag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configPath string

const DefaultConfigPath = "config/config.toml"

func init() {
	pflag.StringVarP(&configPath, "config", "c", DefaultConfigPath, "Путь до конфигурационного файла")
}

// @title Cacher
// @version 1.0
// @description Микросервис для управления структурами данных типа «ключ — значение».

// @host localhost:8080
// @BasePath /

func main() {
	pflag.Parse()

	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatal(err)
	}

	cacheRepository := memory.NewCacheRepository()
	cacheService := service.NewCacheService(cacheRepository)
	srv := handler.NewHTTPServer(cacheService, cfg.HTTPServer.GetDomain())

	go func() {
		log.Printf("Запуск сервера, %s", cfg.HTTPServer.GetDomain())

		if err := srv.Start(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Ошибка при старте сервера: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Выключение сервера ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Сервер принудительно остановлен:", err)
	}

	log.Println("Сервер остановлен")
}
