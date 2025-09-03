package service

import (
	"context"
	"gw-currency-wallet/internal/infrastructure"
	"gw-currency-wallet/internal/storage"
	"gw-currency-wallet/internal/storage/cache"
	"gw-currency-wallet/internal/storage/models"
	"gw-currency-wallet/pkg/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gw-currency-wallet/pkg/kafka"
)

type AuthService interface {
	Register(c context.Context, input models.UserRegister) error
	Login(c context.Context, userInput *models.UserLogin) (string, error)
}

type ExchangeService interface {
	GetRates(c context.Context) (map[string]string, error)
	GetRate(c context.Context, fromCurrency, toCurrency string) (string, error)
	ExchangeCurrency(c context.Context, userID uuid.UUID, fromCurrency string, toCurrency string, amount decimal.Decimal, exchangedAmount decimal.Decimal) (models.WalletResponse, error)
}

type WalletService interface {
	GetBalance(c context.Context, userID uuid.UUID) (models.WalletResponse, error)
	Deposit(c context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, error)
	Withdraw(c context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, error)
}

// для моков
type UserStorage interface {
	CreateUser(ctx context.Context, user models.UserRegister) error
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type WalletStorage interface {
	GetBalance(ctx context.Context, userID uuid.UUID) (models.WalletResponse, error)
	Deposit(ctx context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, uuid.UUID, error)
	Withdraw(ctx context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, uuid.UUID, error)
}

type ExchangeStorage interface {
	Exchange(ctx context.Context, userID uuid.UUID, fromCurrency, toCurrency string, amount, exchangedAmount decimal.Decimal) (models.WalletResponse, uuid.UUID, error)
}

type JWTManager interface {
	GenerateToken(userID uuid.UUID, email string) (string, error)
	ParseToken(token string) (*models.Claims, error)
}

type Service struct {
	AuthService
	ExchangeService
	WalletService
}

func NewService(
	stor *storage.Storage,
	logger *zap.Logger,
	jwtManager *utils.JWTManager,
	exClient *grpc.ExchangeClient,
	cache cache.Cache,
	kafkaProducer *kafka.Producer,
) *Service {
	return &Service{
		AuthService:     NewAuthService(stor, logger, jwtManager),
		ExchangeService: NewExchangeService(exClient, cache, logger, stor, kafkaProducer),
		WalletService:   NewWalletService(stor, logger, kafkaProducer),
	}
}
