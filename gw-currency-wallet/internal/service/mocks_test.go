// service/mocks_test.go
package service

import (
	"context"
	"gw-currency-wallet/internal/storage/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ===== UserStorage (для AuthService) =====
type mockUserStorage struct {
	createUserFn     func(ctx context.Context, user models.UserRegister) error
	getUserByEmailFn func(ctx context.Context, email string) (models.User, error)
}

func (m *mockUserStorage) CreateUser(ctx context.Context, user models.UserRegister) error {
	return m.createUserFn(ctx, user)
}

func (m *mockUserStorage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return m.getUserByEmailFn(ctx, email)
}

// ===== WalletStorage (для WalletService) =====
type mockWalletStorage struct {
	getBalanceFn func(ctx context.Context, userID uuid.UUID) (models.WalletResponse, error)
	depositFn    func(ctx context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, uuid.UUID, error)
	withdrawFn   func(ctx context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, uuid.UUID, error)
}

func (m *mockWalletStorage) GetBalance(ctx context.Context, userID uuid.UUID) (models.WalletResponse, error) {
	return m.getBalanceFn(ctx, userID)
}

func (m *mockWalletStorage) Deposit(ctx context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, uuid.UUID, error) {
	return m.depositFn(ctx, userID, currency, amount)
}

func (m *mockWalletStorage) Withdraw(ctx context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, uuid.UUID, error) {
	return m.withdrawFn(ctx, userID, currency, amount)
}

// ===== ExchangeStorage (для ExchangeService) =====
type mockExchangeStorage struct {
	exchangeFn func(ctx context.Context, userID uuid.UUID, fromCurrency, toCurrency string, amount, exchangedAmount decimal.Decimal) (models.WalletResponse, uuid.UUID, error)
}

func (m *mockExchangeStorage) Exchange(ctx context.Context, userID uuid.UUID, fromCurrency, toCurrency string, amount, exchangedAmount decimal.Decimal) (models.WalletResponse, uuid.UUID, error) {
	return m.exchangeFn(ctx, userID, fromCurrency, toCurrency, amount, exchangedAmount)
}

// ===== JWTManager (для AuthService) =====
type mockJWT struct {
	generateTokenFn func(userID uuid.UUID, email string) (string, error)
	parseTokenFn    func(token string) (*models.Claims, error)
}

func (m *mockJWT) GenerateToken(userID uuid.UUID, email string) (string, error) {
	return m.generateTokenFn(userID, email)
}

func (m *mockJWT) ParseToken(token string) (*models.Claims, error) {
	return m.parseTokenFn(token)
}
