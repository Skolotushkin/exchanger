package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"gw-currency-wallet/internal/handler"
	"gw-currency-wallet/internal/service"
	"gw-currency-wallet/internal/storage/models"
)

// --- mock WalletService ---
type mockWalletService struct{ mock.Mock }

func (m *mockWalletService) GetBalance(ctx context.Context, userID uuid.UUID) (models.WalletResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(models.WalletResponse), args.Error(1)
}

func (m *mockWalletService) Deposit(ctx context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, error) {
	args := m.Called(ctx, userID, currency, amount)
	return args.Get(0).(models.WalletResponse), args.Error(1)
}

func (m *mockWalletService) Withdraw(ctx context.Context, userID uuid.UUID, currency string, amount decimal.Decimal) (models.WalletResponse, error) {
	args := m.Called(ctx, userID, currency, amount)
	return args.Get(0).(models.WalletResponse), args.Error(1)
}

func setupWalletHandler(mockSvc *mockWalletService) (*gin.Engine, uuid.UUID) {
	h := handler.NewWalletHandler(&service.Service{WalletService: mockSvc}, zap.NewNop())
	r := gin.New()
	uid := uuid.New()
	r.POST("/deposit", func(c *gin.Context) { c.Set("userID", uid); h.Deposit(c) })
	r.POST("/withdraw", func(c *gin.Context) { c.Set("userID", uid); h.Withdraw(c) })
	r.GET("/balance", func(c *gin.Context) { c.Set("userID", uid); h.GetBalance(c) })
	return r, uid
}

func TestWalletHandler_Deposit_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockWalletService)
	r, uid := setupWalletHandler(mockSvc)

	reqBody := models.WalletTransaction{Currency: "USD", Amount: "100"}
	body, _ := json.Marshal(reqBody)

	mockSvc.On("Deposit", mock.Anything, uid, "USD", mock.Anything).
		Return(models.WalletResponse{Balances: map[string]decimal.Decimal{"USD": decimal.NewFromInt(100)}}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/deposit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Deposit successful")
}
