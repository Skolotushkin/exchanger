package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"gw-currency-wallet/pkg/kafka"
)

const (
	EventLargeDeposit  = "LARGE_DEPOSIT"
	EventLargeWithdraw = "LARGE_WITHDRAW"
	EventLargeExchange = "LARGE_EXCHANGE"
	Threshold          = 30000
)

type KafkaEvent struct {
	TransactionID string `json:"transaction_id"`
	UserID        string `json:"user_id"`
	Operation     string `json:"operation"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
}

func sendLargeOperationEvent(
	ctx context.Context,
	userID uuid.UUID,
	operation string,
	amount decimal.Decimal,
	currency string,
	kafkaProducer *kafka.Producer,
	logger *zap.Logger,
	transactionID uuid.UUID,
) {
	if amount.GreaterThan(decimal.NewFromInt(Threshold)) {
		msg := KafkaEvent{
			TransactionID: transactionID.String(),
			UserID:        userID.String(),
			Operation:     operation,
			Amount:        amount.String(),
			Currency:      currency,
		}
		payload, _ := json.Marshal(msg)
		if err := kafkaProducer.SendMessage(ctx, []byte(userID.String()), payload); err != nil {
			logger.Error("failed to send kafka event", zap.Error(err))
		} else {
			logger.Info("large operation event sent",
				zap.String("operation", operation),
				zap.String("amount", amount.String()),
				zap.String("currency", currency),
				zap.String("transaction_id", transactionID.String()),
			)
		}
	}
}
