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

// --- mock ExchangeService ---
type mockExchangeService struct{ mock.Mock }

func (m *mockExchangeService) GetRates(ctx context.Context) (map[string]string, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *mockExchangeService) GetRate(ctx context.Context, fromCurrency, toCurrency string) (string, error) {
	args := m.Called(ctx, fromCurrency, toCurrency)
	return args.String(0), args.Error(1)
}

func (m *mockExchangeService) ExchangeCurrency(ctx context.Context, userID uuid.UUID, fromCurrency, toCurrency string, amount, exchangedAmount decimal.Decimal) (models.WalletResponse, error) {
	args := m.Called(ctx, userID, fromCurrency, toCurrency, amount, exchangedAmount)
	return args.Get(0).(models.WalletResponse), args.Error(1)
}

func setupExchangeHandler(mockSvc *mockExchangeService) (*gin.Engine, uuid.UUID) {
	h := handler.NewExchangeHandler(&service.Service{ExchangeService: mockSvc}, zap.NewNop())
	r := gin.New()
	uid := uuid.New()
	r.GET("/rates", h.GetExchangeRates)
	r.POST("/exchange", func(c *gin.Context) { c.Set("userID", uid); h.ExchangeCurrency(c) })
	return r, uid
}

func TestExchangeHandler_GetRates_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockExchangeService)
	r, _ := setupExchangeHandler(mockSvc)

	mockSvc.On("GetRates", mock.Anything).Return(map[string]string{"USD": "1.0"}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rates", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "USD")
}

func TestExchangeHandler_ExchangeCurrency_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockExchangeService)
	r, uid := setupExchangeHandler(mockSvc)

	reqBody := models.ExchangeRequest{FromCurrency: "USD", ToCurrency: "EUR", Amount: "100"}
	body, _ := json.Marshal(reqBody)

	mockSvc.On("GetRate", mock.Anything, "USD", "EUR").Return("0.9", nil)
	mockSvc.On("ExchangeCurrency", mock.Anything, uid, "USD", "EUR", mock.Anything, mock.Anything).
		Return(models.WalletResponse{Balances: map[string]decimal.Decimal{"EUR": decimal.NewFromInt(90)}}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/exchange", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Exchange successful")
}
