package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"gw-currency-wallet/internal/storage/models"
)

func zapLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func TestAuthService_Register_Success(t *testing.T) {
	mStor := &mockUserStorage{
		createUserFn: func(ctx context.Context, user models.UserRegister) error {
			return nil
		},
	}

	mJWT := &mockJWT{}

	svc := NewAuthService(mStor, zapLogger(), mJWT)
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	err := svc.Register(context.Background(), models.UserRegister{Email: "test@test.com", Password: string(hashed)})
	require.NoError(t, err)
}

func TestAuthService_Register_Failure(t *testing.T) {
	expectedErr := errors.New("db error")

	mStor := &mockUserStorage{
		createUserFn: func(ctx context.Context, user models.UserRegister) error {
			return expectedErr
		},
	}

	mJWT := &mockJWT{}

	svc := NewAuthService(mStor, zapLogger(), mJWT)
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	err := svc.Register(context.Background(), models.UserRegister{Email: "test@test.com", Password: string(hashed)})
	require.Error(t, err)
	require.Equal(t, expectedErr, err)
}

func TestAuthService_Login_Success(t *testing.T) {
	// хэшируем пароль так же, как в Register
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := models.User{
		ID:       uuid.New(),
		Email:    "test@test.com",
		Password: string(hashedPassword),
	}

	mStor := &mockUserStorage{
		getUserByEmailFn: func(ctx context.Context, email string) (models.User, error) {
			return user, nil
		},
	}

	mJWT := &mockJWT{
		generateTokenFn: func(userID uuid.UUID, email string) (string, error) {
			return "test-token", nil
		},
	}

	svc := NewAuthService(mStor, zapLogger(), mJWT)

	token, err := svc.Login(context.Background(), &models.UserLogin{Email: user.Email, Password: "password"})
	require.NoError(t, err)
	require.Equal(t, "test-token", token)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mStor := &mockUserStorage{
		getUserByEmailFn: func(ctx context.Context, email string) (models.User, error) {
			return models.User{}, errors.New("not found")
		},
	}

	mJWT := &mockJWT{}

	svc := NewAuthService(mStor, zapLogger(), mJWT)

	_, err := svc.Login(context.Background(), &models.UserLogin{Email: "unknown@test.com", Password: "any"})
	require.Error(t, err)
}
