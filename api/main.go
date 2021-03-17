package main

import (
	"context"
	"github.com/koind/cacher/api/internal/config"
	"github.com/koind/cacher/api/internal/domain/service"
	"github.com/koind/cacher/api/internal/handler"
	"github.com/koind/cacher/api/internal/storage/memory"
	"github.com/spf13/pflag"
	"log"
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

func main() {
	pflag.Parse()

	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatal(err)
	}

	cacheRepository := memory.NewCacheRepository()
	cacheService := service.NewCacheService(cacheRepository)
	srv := handler.NewHTTPServer(cacheService, cfg.HTTPServer.GetDomain())

	// handle shutdown gracefully
	quit := make(chan os.Signal, 1)
	done := make(chan error, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		err := srv.Shutdown(ctx)

		done <- err
	}()

	log.Printf("Запуск сервера, %s", cfg.HTTPServer.GetDomain())
	log.Printf("Результат запуска сервера, %v", srv.Start())

	err = <-done
	log.Printf("Остановка сервера, %v", err)
}
