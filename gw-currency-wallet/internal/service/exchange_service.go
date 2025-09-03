package service

import (
	"context"
	"fmt"
	cache "gw-currency-wallet/internal/storage/cache"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"gw-currency-wallet/internal/infrastructure"
	"gw-currency-wallet/internal/storage/models"
	"gw-currency-wallet/pkg/kafka"
)

const (
	cacheKeyRatePair = "exchange:rate:%s:%s" // rate:USD:EUR
	cacheTTL         = 2 * time.Minute
)

type exchangeService struct {
	exClient *grpc.ExchangeClient
	cache    cache.Cache
	logger   *zap.Logger
	storage  ExchangeStorage
	kafka    *kafka.Producer
}

func NewExchangeService(
	exClient *grpc.ExchangeClient,
	cache cache.Cache,
	logger *zap.Logger,
	storage ExchangeStorage,
	kafka *kafka.Producer,
) ExchangeService {
	return &exchangeService{
		exClient: exClient,
		cache:    cache,
		logger:   logger,
		storage:  storage,
		kafka:    kafka,
	}
}

func (e *exchangeService) GetRates(ctx context.Context) (map[string]string, error) {
	// 1) пробуем из кеша
	if e.cache != nil {
		if rates, err := e.cache.GetRates(ctx); err == nil {
			return rates, nil
		}
	}

	// 2) идём в gRPC
	rates, err := e.exClient.GetRates(ctx)
	if err != nil {
		return nil, err
	}

	// 3) положим в кеш (TTL 5 минут)
	if e.cache != nil {
		_ = e.cache.SetRates(ctx, rates, 5*time.Minute)
	}

	return rates, nil
}

func (e *exchangeService) GetRate(ctx context.Context, fromCurrency, toCurrency string) (string, error) {
	if fromCurrency == toCurrency {
		return "1", nil
	}

	key := fmt.Sprintf(cacheKeyRatePair, fromCurrency, toCurrency)

	// 1) пробуем кеш пары
	if e.cache != nil {
		if v, err := e.cache.Get(ctx, key); err == nil && v != "" {
			return v, nil
		}
	}

	// 2)  идём в gRPC
	rate, err := e.exClient.GetRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		return "", err
	}

	// 3) кладём в кеш
	if e.cache != nil {
		if err := e.cache.Set(ctx, key, rate, cacheTTL); err != nil {
			e.logger.Warn("failed to set rate cache", zap.Error(err))
		}
	}

	return rate, nil
}

func (e *exchangeService) ExchangeCurrency(
	ctx context.Context,
	userID uuid.UUID,
	fromCurrency string,
	toCurrency string,
	amount decimal.Decimal,
	exchangedAmount decimal.Decimal,
) (models.WalletResponse, error) {
	resp, transactionID, err := e.storage.Exchange(ctx, userID, fromCurrency, toCurrency, amount, exchangedAmount)
	if err != nil {

		return resp, err
	}

	sendLargeOperationEvent(ctx, userID, EventLargeExchange, amount, fromCurrency, e.kafka, e.logger, transactionID)
	return resp, nil
}
