package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"gw-currency-wallet/internal/handler"
	"gw-currency-wallet/internal/service"
	"gw-currency-wallet/internal/storage/models"
)

// --- mock AuthService ---
type mockAuthService struct{ mock.Mock }

func (m *mockAuthService) Register(ctx context.Context, req models.UserRegister) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockAuthService) Login(ctx context.Context, req *models.UserLogin) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func TestAuthHandler_Register_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(mockAuthService)
	h := handler.NewAuthHandler(&service.Service{AuthService: mockSvc}, zap.NewNop())

	router := gin.New()
	router.POST("/auth/register", h.Register)

	reqBody := models.UserRegister{Email: "test@example.com", Password: "123456"}
	body, _ := json.Marshal(reqBody)

	mockSvc.On("Register", mock.Anything, reqBody).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := handler.NewAuthHandler(&service.Service{}, zap.NewNop())

	router := gin.New()
	router.POST("/auth/register", h.Register)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte(`{bad json}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(mockAuthService)
	h := handler.NewAuthHandler(&service.Service{AuthService: mockSvc}, zap.NewNop())

	router := gin.New()
	router.POST("/auth/login", h.Login)

	reqBody := models.UserLogin{Email: "test@example.com", Password: "123456"}
	body, _ := json.Marshal(reqBody)

	mockSvc.On("Login", mock.Anything, &reqBody).Return("fake-token", nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "fake-token")
}
