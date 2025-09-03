package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// PoolIface описывает минимальный интерфейс пула/мока, используемый в хранилище.
// Это нужно, чтобы тесты могли прокинуть pgxmock.PgxPoolIface без изменения продакшен-кода.
type PoolIface interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

type PostgresStorage struct {
	pool   PoolIface
	logger *zap.Logger
}

func NewPostgresConnector(dsn string, logger *zap.Logger) (*PostgresStorage, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres dsn: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Простая проверка соединения через Exec (SELECT 1)
	if _, err := pool.Exec(ctx, "SELECT 1"); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	logger.Info("connected to postgres")

	return &PostgresStorage{
		pool:   pool,
		logger: logger,
	}, nil
}

// NewPostgresStorageWithPool создает PostgresStorage из уже существующего пула или мока.
func NewPostgresStorageWithPool(pool PoolIface, logger *zap.Logger) *PostgresStorage {
	return &PostgresStorage{
		pool:   pool,
		logger: logger,
	}
}

func (p *PostgresStorage) Close() {
	p.pool.Close()
}
