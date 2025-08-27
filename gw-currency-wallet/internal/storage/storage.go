package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"gw-currency-wallet/internal/storage/models"
	"gw-currency-wallet/internal/storage/postgres"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user models.UserRegister) error
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type WalletStorage interface {
	GetBalance(ctx context.Context, userID uuid.UUID) (models.WalletResponse, error)
	Deposit(ctx context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, uuid.UUID, error)
	Withdraw(ctx context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, uuid.UUID, error)
	Exchange(
		ctx context.Context,
		userID uuid.UUID,
		fromCurrency string,
		toCurrency string,
		amount decimal.Decimal,
		exchangedAmount decimal.Decimal,
	) (models.WalletResponse, uuid.UUID, error)
}

type Storage struct {
	UserStorage
	WalletStorage
}

func NewStorage(db *postgres.PostgresStorage) *Storage {
	return &Storage{
		UserStorage:   db,
		WalletStorage: db,
	}
}
