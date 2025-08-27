package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	exchpb "github.com/Skolotushkin/proto-exchange/gen"
)

const (
	defaultTimeout       = 2 * time.Second
	dialTimeout          = 5 * time.Second
	reconnectTimeout     = 3 * time.Second
	maxReconnectAttempts = 3
)

// ExchangeClient — thin wrapper над gRPC-клиентом
type ExchangeClient struct {
	addr   string
	conn   *grpc.ClientConn
	client exchpb.ExchangeServiceClient
	log    *zap.Logger
}

// NewExchangeClient создаёт соединение и gRPC-клиент.
// addr ожидается вида "host:port" (например, "gw-exchanger:50051").
func NewExchangeClient(addr string, logger *zap.Logger) (*ExchangeClient, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithReturnConnectionError(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, dialOpts...)
	if err != nil {
		logger.Error("failed to dial exchange grpc",
			zap.String("addr", addr),
			zap.Error(err))
		return nil, fmt.Errorf("dial exchange grpc: %w", err)
	}

	logger.Info("connected to exchange grpc", zap.String("addr", addr))

	return &ExchangeClient{
		addr:   addr,
		conn:   conn,
		client: exchpb.NewExchangeServiceClient(conn),
		log:    logger,
	}, nil
}

// Close закрывает gRPC-соединение.
func (c *ExchangeClient) Close() error {
	if c.conn != nil {
		c.log.Info("closing exchange grpc connection", zap.String("addr", c.addr))
		return c.conn.Close()
	}
	return nil
}

// EnsureConnection проверяет и восстанавливает соединение при необходимости
func (c *ExchangeClient) EnsureConnection(ctx context.Context) error {
	state := c.conn.GetState()
	if state == connectivity.Ready {
		return nil
	}

	c.log.Warn("connection not ready, attempting to reconnect",
		zap.String("state", state.String()))

	// Попытаться переподключиться
	for attempt := 1; attempt <= maxReconnectAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if c.conn.WaitForStateChange(ctx, state) {
				newState := c.conn.GetState()
				if newState == connectivity.Ready {
					c.log.Info("connection restored")
					return nil
				}
				state = newState
			}

			time.Sleep(reconnectTimeout)
		}
	}

	return fmt.Errorf("failed to restore connection after %d attempts", maxReconnectAttempts)
}

// GetRates возвращает все курсы в виде map[string]string
func (c *ExchangeClient) GetRates(ctx context.Context) (map[string]string, error) {
	if err := c.EnsureConnection(ctx); err != nil {
		return nil, fmt.Errorf("connection check failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	resp, err := c.client.GetExchangeRates(ctx, &exchpb.Empty{})
	if err != nil {
		c.log.Error("GetExchangeRates failed", zap.Error(err))
		return nil, fmt.Errorf("get exchange rates: %w", err)
	}

	out := make(map[string]string, len(resp.GetRates()))
	for k, v := range resp.GetRates() {
		out[k] = decimal.NewFromFloat32(v).String()
	}
	return out, nil
}

// GetRate возвращает один курс как string
func (c *ExchangeClient) GetRate(ctx context.Context, fromCurrency, toCurrency string) (string, error) {
	if err := c.EnsureConnection(ctx); err != nil {
		return "", fmt.Errorf("connection check failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	req := &exchpb.CurrencyRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	}

	resp, err := c.client.GetExchangeRateForCurrency(ctx, req)
	if err != nil {
		c.log.Error("GetExchangeRateForCurrency failed",
			zap.String("from", fromCurrency),
			zap.String("to", toCurrency),
			zap.Error(err),
		)
		return "", fmt.Errorf("get exchange rate for %s/%s: %w", fromCurrency, toCurrency, err)
	}

	return decimal.NewFromFloat32(resp.GetRate()).String(), nil
}

// CheckHealth проверяет доступность сервиса
func (c *ExchangeClient) CheckHealth(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err := c.client.GetExchangeRates(ctx, &exchpb.Empty{})
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	return nil
}
