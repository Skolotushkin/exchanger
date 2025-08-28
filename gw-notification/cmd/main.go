package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"gw-notification/internal/config"
	"gw-notification/internal/kafka"
	"gw-notification/internal/mongodb"
)

func main() {
	// конфигурация
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		os.Exit(1)
	}

	// init logger
	logger, _ := zap.NewProduction()
	logger.Named("gw-notification")
	defer logger.Sync()

	// init Mongo
	client, err := mongodb.NewClient(cfg.MongoURI)
	if err != nil {
		logger.Fatal("mongo connection failed", zap.Error(err))
	}
	repo := mongodb.NewRepository(client, cfg.MongoDBName, cfg.MongoCollection)
	logger.Info("connected to MongoDB")

	// start kafka consumer
	go kafka.ConsumeKafka(cfg.KafkaBrokers, cfg.KafkaTopic, repo, logger)

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down gw-notification gracefully")
}
