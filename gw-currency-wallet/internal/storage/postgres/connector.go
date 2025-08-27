package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgresConnector(dsn string, logger *zap.Logger) (*PostgresStorage, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres dsn: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	logger.Info("connected to postgres")

	return &PostgresStorage{
		pool:   pool,
		logger: logger,
	}, nil
}

func (p *PostgresStorage) Close() {
	p.pool.Close()
}
