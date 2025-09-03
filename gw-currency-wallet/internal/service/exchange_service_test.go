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

func TestExchangeService_ExchangeCurrency_Success(t *testing.T) {
	expected := models.WalletResponse{Balances: map[string]decimal.Decimal{"EUR": decimal.NewFromInt(90)}}

	mStor := &mockExchangeStorage{
		exchangeFn: func(ctx context.Context, userID uuid.UUID, from, to string, amount, exchanged decimal.Decimal) (models.WalletResponse, uuid.UUID, error) {
			return expected, uuid.New(), nil
		},
	}

	svc := NewExchangeService(nil, nil, zapLogger(), mStor, nil)

	resp, err := svc.ExchangeCurrency(context.Background(), uuid.New(), "USD", "EUR", decimal.NewFromInt(100), decimal.NewFromInt(90))
	require.NoError(t, err)
	require.Equal(t, expected, resp)
}

func TestExchangeService_ExchangeCurrency_Failure(t *testing.T) {
	expectedErr := errors.New("exchange failed")

	mStor := &mockExchangeStorage{
		exchangeFn: func(ctx context.Context, userID uuid.UUID, from, to string, amount, exchanged decimal.Decimal) (models.WalletResponse, uuid.UUID, error) {
			return models.WalletResponse{}, uuid.Nil, expectedErr
		},
	}

	svc := NewExchangeService(nil, nil, zapLogger(), mStor, nil)

	_, err := svc.ExchangeCurrency(context.Background(), uuid.New(), "USD", "EUR", decimal.NewFromInt(100), decimal.NewFromInt(90))
	require.Error(t, err)
	require.Equal(t, expectedErr, err)
}
