package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"gw-currency-wallet/internal/service"
	"gw-currency-wallet/internal/storage/models"
)

type walletHandler struct {
	svc    *service.Service
	logger *zap.Logger
}

func NewWalletHandler(svc *service.Service, logger *zap.Logger) *walletHandler {
	return &walletHandler{svc: svc, logger: logger}
}

// GetBalance godoc
// @Summary Get wallet balance
// @Description Get current balance for authenticated user
// @Tags Wallet
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /wallet/balance [get]
func (h *walletHandler) GetBalance(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	balance, err := h.svc.WalletService.GetBalance(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get balance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// Deposit godoc
// @Summary Deposit funds
// @Description Deposit money into wallet
// @Tags Wallet
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body models.WalletTransaction true "Deposit request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /wallet/deposit [post]
func (h *walletHandler) Deposit(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req models.WalletTransaction
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid deposit request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	resp, err := h.svc.WalletService.Deposit(c, userID, req.Currency, amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount or currency"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "Deposit successful",
		"new_balance": resp.Balances,
	})
}

// Withdraw godoc
// @Summary Withdraw funds
// @Description Withdraw money from wallet
// @Tags Wallet
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body models.WalletTransaction true "Withdraw request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /wallet/withdraw [post]
func (h *walletHandler) Withdraw(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var req models.WalletTransaction
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid withdraw request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		h.logger.Warn("amount zero", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	resp, err := h.svc.WalletService.Withdraw(c, userID, req.Currency, amount)
	if err != nil {
		h.logger.Warn("withdraw failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds or invalid amount"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "Withdrawal successful",
		"new_balance": resp.Balances,
	})
}
