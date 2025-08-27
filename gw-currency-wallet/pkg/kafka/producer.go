package kafka

import (
	"context"

	kgo "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kgo.Writer
	logger *zap.Logger
}

func NewProducer(brokers []string, topic string, logger *zap.Logger) *Producer {
	return &Producer{
		writer: &kgo.Writer{
			Addr:     kgo.TCP(brokers...),
			Topic:    topic,
			Balancer: &kgo.LeastBytes{},
		},
		logger: logger,
	}
}

func (p *Producer) SendMessage(ctx context.Context, key, value []byte) error {
	err := p.writer.WriteMessages(ctx, kgo.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		p.logger.Error("failed to send kafka message", zap.Error(err))
		return err
	}
	p.logger.Info("sent kafka message")
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
