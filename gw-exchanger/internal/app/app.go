package app

import (
	"context"
	"go.uber.org/zap"
	"gw-exchanger/internal/config"
	"gw-exchanger/internal/server"
	"gw-exchanger/internal/service"
	"gw-exchanger/internal/storages/postgres"
	"os"
	"os/signal"
	"syscall"
)

func StartApplication(cfg *config.Config, logger *zap.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// подключаемся к Postgres
	storage, err := postgres.NewPostgresStorage(ctx,
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSL)
	if err != nil {
		return err
	}
	defer storage.Close()

	// создаем слой сервисов
	exchangerService := service.NewExchangerService(storage, logger)

	// запускаем gRPC сервер
	go func() {
		if err := server.RunGRPCServer(ctx, exchangerService, cfg, logger); err != nil {
			logger.Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	cancel()
	logger.Info("✅ application shutdown complete")

	return nil
}
