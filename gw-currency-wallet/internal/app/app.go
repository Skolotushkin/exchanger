package app

import (
	"fmt"
	// "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"go.uber.org/zap"
	config "gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/handler"
	grpc "gw-currency-wallet/internal/infrastructure"
	"gw-currency-wallet/internal/server"
	"gw-currency-wallet/internal/service"
	"gw-currency-wallet/internal/storage"
	"gw-currency-wallet/internal/storage/cache"
	"gw-currency-wallet/internal/storage/postgres"
	"gw-currency-wallet/pkg/kafka"
	"gw-currency-wallet/pkg/utils"
	"strings"
)

func StartApplication(cfg *config.Config) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Postgres
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSL,
	)
	pgStore, err := postgres.NewPostgresConnector(dsn, logger)
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}

	// Cache: Redis или InMemory
	var c cache.Cache
	if cfg.Redis.RedisEnabled {
		c, err = cache.NewRedisCache(cfg.Redis.Addr /* cfg.Redis.Password*/, cfg.Redis.DB, logger)
		if err != nil {
			logger.Fatal("failed to init redis", zap.Error(err))
		}
	} else {
		c = cache.NewInMemoryCache()
	}
	defer c.Close()

	// gRPC клиент к exchanger
	exClient, err := grpc.NewExchangeClient(cfg.ExchangeService.Addr, logger)
	if err != nil {
		logger.Fatal("failed to init exchange client", zap.Error(err))
	}

	// JWT
	jwtManager := utils.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.JWTTTL)

	// Kafka
	kafkaProducer := kafka.NewProducer(strings.Split(cfg.Kafka.Brokers, ","), cfg.Kafka.Topic, logger)
	defer kafkaProducer.Close()

	// DI
	repo := storage.NewStorage(pgStore)
	services := service.NewService(repo, logger, jwtManager, exClient, c, kafkaProducer)
	handlers := handler.NewHandler(services, logger)

	// HTTP server
	server.SetupAndRunServer(cfg, handlers.InitRoutes(logger, jwtManager), logger)
	return nil
}
