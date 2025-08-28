package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"

	"go.uber.org/zap"

	"gw-notification/internal/models"
	"gw-notification/internal/mongodb"
)

func ConsumeKafka(broker, topic string, repo *mongodb.Repository, logger *zap.Logger) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "notification-group",
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			logger.Error("failed to read kafka message", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

		var event models.TransactionHistory
		if err := json.Unmarshal(m.Value, &event); err != nil {
			logger.Error("failed to unmarshal kafka message", zap.Error(err))
			continue
		}

		event.CreatedAt = time.Now()

		if err := repo.SaveTransaction(event); err != nil {
			logger.Error("failed to save transaction",
				zap.String("transaction_id", event.TransactionID),
				zap.String("user_id", event.UserID),
				zap.Error(err),
			)
			continue
		}

		logger.Info("transaction saved successfully",
			zap.String("transaction_id", event.TransactionID),
			zap.String("user_id", event.UserID),
			zap.String("operation", event.Operation),
			zap.String("amount", event.Amount),
			zap.String("currency", event.Currency),
		)
	}
}
