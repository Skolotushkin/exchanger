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

type exchangeHandler struct {
	svc    *service.Service
	logger *zap.Logger
}

func NewExchangeHandler(svc *service.Service, logger *zap.Logger) *exchangeHandler {
	return &exchangeHandler{svc: svc, logger: logger}
}

// GetExchangeRates godoc
// @Summary Get exchange rates
// @Description Get list of available exchange rates
// @Tags Exchange
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /exchange/rates [get]
func (h *exchangeHandler) GetExchangeRates(c *gin.Context) {
	rates, err := h.svc.ExchangeService.GetRates(c)
	if err != nil {
		h.logger.Error("failed to get exchange rates", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exchange rates"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rates": rates})
}

// ExchangeCurrency godoc
// @Summary Exchange currency
// @Description Exchange one currency to another
// @Tags Exchange
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body models.ExchangeRequest true "Exchange request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /api/v1/exchange [post]
func (h *exchangeHandler) ExchangeCurrency(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var req models.ExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid exchange request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		h.logger.Warn("amount is zero", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	// получаем курс и считаем exchangedAmount
	rateStr, err := h.svc.ExchangeService.GetRate(c, req.FromCurrency, req.ToCurrency)
	if err != nil {
		h.logger.Warn("invalid exchange request getrate", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid currencies"})
		return
	}
	rate, err := decimal.NewFromString(rateStr)
	if err != nil {
		h.logger.Warn("invalid exchange decimal.newformstring", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse rate"})
		return
	}
	exchanged := amount.Mul(rate)

	// совершаем обмен (транзакция в storage.Exchange)
	resp, err := h.svc.ExchangeService.ExchangeCurrency(c, userID, req.FromCurrency, req.ToCurrency, amount, exchanged)
	if err != nil {
		h.logger.Warn("exchange failed svc.ExchangeService.ExchangeCurrency", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds or invalid currencies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Exchange successful",
		"exchanged_amount": exchanged.String(),
		"new_balance":      resp.Balances,
	})
}
