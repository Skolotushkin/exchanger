// Package handler Currency Wallet API
//
// @title Currency Wallet API
// @version 1.0
// @description REST API для кошелька с поддержкой мультивалюты и обмена.
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.email support@example.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"gw-currency-wallet/internal/handler/middleware"
	"gw-currency-wallet/internal/service"
	"gw-currency-wallet/pkg/utils"
)

type Handler struct {
	auth     *authHandler
	wallet   *walletHandler
	exchange *exchangeHandler
}

func NewHandler(
	svc *service.Service,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		auth:     NewAuthHandler(svc, logger),
		wallet:   NewWalletHandler(svc, logger),
		exchange: NewExchangeHandler(svc, logger),
	}
}

func (h *Handler) InitRoutes(logger *zap.Logger, jwtManager *utils.JWTManager) *gin.Engine {
	router := gin.New()
	router.Use(middleware.GinZapLogger(logger), middleware.RecoverMiddleware(logger))

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := router.Group("/api/v1")

	// --- Public routes ---
	auth := apiV1.Group("/auth")
	{
		auth.POST("/register", h.auth.Register)
		auth.POST("/login", h.auth.Login)
	}

	// --- Protected routes ---
	protected := apiV1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager, logger))
	{
		w := protected.Group("/wallet")
		{
			w.GET("/balance", h.wallet.GetBalance)
			w.POST("/deposit", h.wallet.Deposit)
			w.POST("/withdraw", h.wallet.Withdraw)
		}
		ex := protected.Group("/exchange")
		{
			ex.GET("/rates", h.exchange.GetExchangeRates)
			ex.POST("", h.exchange.ExchangeCurrency)
		}
	}

	// healthcheck
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}
