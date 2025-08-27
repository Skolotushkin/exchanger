package service

import (
	"context"
	"gw-currency-wallet/internal/storage"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"gw-currency-wallet/internal/storage/models"
	"gw-currency-wallet/pkg/kafka"
)

type walletService struct {
	storage *storage.Storage
	logger  *zap.Logger
	kafka   *kafka.Producer
}

func NewWalletService(storage *storage.Storage, logger *zap.Logger, kafka *kafka.Producer) WalletService {
	return &walletService{storage: storage, logger: logger, kafka: kafka}
}

func (w *walletService) GetBalance(c context.Context, userID uuid.UUID) (models.WalletResponse, error) {
	return w.storage.GetBalance(c, userID)
}

func (w *walletService) Deposit(c context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, error) {
	resp, transactionID, err := w.storage.Deposit(c, userID, currency, amount)
	if err != nil {
		return resp, err
	}
	sendLargeOperationEvent(c, userID, EventLargeDeposit, amount, currency, w.kafka, w.logger, transactionID)
	return resp, nil
}

func (w *walletService) Withdraw(c context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, error) {
	resp, transactionID, err := w.storage.Withdraw(c, userID, currency, amount)
	if err != nil {
		return resp, err
	}
	sendLargeOperationEvent(c, userID, EventLargeWithdraw, amount, currency, w.kafka, w.logger, transactionID)
	return resp, nil
}
