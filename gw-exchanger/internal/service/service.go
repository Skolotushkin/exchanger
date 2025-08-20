package service

import (
	"context"

	pb "github.com/Skolotushkin/proto-exchange/gen"
	"go.uber.org/zap"
	"gw-exchanger/internal/storages"
)

type ExchangerService struct {
	storage storages.Storage
	logger  *zap.Logger
	pb.UnimplementedExchangeServiceServer
}

func NewExchangerService(storage storages.Storage, logger *zap.Logger) *ExchangerService {
	return &ExchangerService{
		storage: storage,
		logger:  logger,
	}
}

func (s *ExchangerService) GetExchangeRates(ctx context.Context, _ *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	rates, err := s.storage.GetAllRates(ctx)
	if err != nil {
		s.logger.Error("failed to get all rates", zap.Error(err))
		return nil, err
	}
	return &pb.ExchangeRatesResponse{Rates: rates}, nil
}

func (s *ExchangerService) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	rate, err := s.storage.GetRate(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		s.logger.Warn("rate not found", zap.String("from", req.FromCurrency), zap.String("to", req.ToCurrency))
		return nil, err
	}

	return &pb.ExchangeRateResponse{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         float32(rate),
	}, nil
}
