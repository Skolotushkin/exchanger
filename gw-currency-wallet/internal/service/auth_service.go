package service

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"gw-currency-wallet/internal/storage/models"
)

type authService struct {
	storage    UserStorage
	logger     *zap.Logger
	jwtManager JWTManager
}

func NewAuthService(storage UserStorage, logger *zap.Logger, jwtManager JWTManager) AuthService {
	return &authService{storage: storage, logger: logger, jwtManager: jwtManager}
}

func (a *authService) Register(ctx context.Context, input models.UserRegister) error {
	// хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Создаем UserRegister с хешированным паролем
	userReg := models.UserRegister{
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := a.storage.CreateUser(ctx, userReg); err != nil {
		a.logger.Error("failed to create user", zap.Error(err))
		return err
	}

	a.logger.Info("user registered", zap.String("email", userReg.Email))
	return nil
}

func (a *authService) Login(ctx context.Context, creds *models.UserLogin) (string, error) {
	user, err := a.storage.GetUserByEmail(ctx, creds.Email)
	if err != nil {
		a.logger.Warn("login failed - user not found", zap.String("email", creds.Email))
		return "", errors.New("invalid email or password")
	}

	// проверка пароля через bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		a.logger.Warn("login failed - wrong password", zap.String("email", creds.Email))
		return "", errors.New("invalid email or password")
	}

	// создаём JWT
	token, err := a.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		a.logger.Error("failed to generate jwt", zap.Error(err))
		return "", err
	}

	a.logger.Info("user logged in", zap.String("email", user.Email))
	return token, nil
}
