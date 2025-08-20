package main

import (
	"fmt"
	"go.uber.org/zap"
	"gw-exchanger/internal/app"
	"gw-exchanger/internal/config"
	"os"
)

func main() {
	// конфигурация
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		os.Exit(1)
	}

	// инициализация логгера
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// запуск приложения
	if err := app.StartApplication(cfg, logger); err != nil {
		logger.Fatal("error starting application", zap.Error(err))
	}
}
