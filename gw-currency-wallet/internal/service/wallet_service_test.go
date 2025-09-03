package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gw-currency-wallet/internal/storage/models"
)

func TestWalletService_GetBalance(t *testing.T) {
	expected := models.WalletResponse{Balances: map[string]decimal.Decimal{"USD": decimal.NewFromInt(100)}}

	mStor := &mockWalletStorage{
		getBalanceFn: func(ctx context.Context, userID uuid.UUID) (models.WalletResponse, error) {
			return expected, nil
		},
	}

	svc := NewWalletService(mStor, zapLogger(), nil)

	resp, err := svc.GetBalance(context.Background(), uuid.New())
	require.NoError(t, err)
	require.Equal(t, expected, resp)
}

func TestWalletService_Deposit_Error(t *testing.T) {
	expectedErr := errors.New("deposit failed")

	mStor := &mockWalletStorage{
		depositFn: func(ctx context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, uuid.UUID, error) {
			return models.WalletResponse{}, uuid.Nil, expectedErr
		},
	}

	svc := NewWalletService(mStor, zapLogger(), nil)

	_, err := svc.Deposit(context.Background(), uuid.New(), "USD", decimal.NewFromInt(100))
	require.Error(t, err)
	require.Equal(t, expectedErr, err)
}
