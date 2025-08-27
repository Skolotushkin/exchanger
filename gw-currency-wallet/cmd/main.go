package main

import (
	"fmt"
	"gw-currency-wallet/internal/app"
	"gw-currency-wallet/internal/config"
	"log"
	"os"
)

func main() {
	// конфигурация
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		os.Exit(1)
	}

	// запуск приложения
	if err := app.StartApplication(cfg); err != nil {
		log.Fatalf("failed to start application: %v", err)
	}
}
