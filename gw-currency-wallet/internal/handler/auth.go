package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gw-currency-wallet/internal/service"
	"gw-currency-wallet/internal/storage/models"
)

type authHandler struct {
	svc    *service.Service
	logger *zap.Logger
}

func NewAuthHandler(svc *service.Service, logger *zap.Logger) *authHandler {
	return &authHandler{svc: svc, logger: logger}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.UserRegister true "User register request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *authHandler) Register(c *gin.Context) {
	var req models.UserRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.svc.AuthService.Register(c.Request.Context(), req); err != nil {
		h.logger.Warn("register failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.UserLogin true "User login request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	var req models.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.svc.AuthService.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Warn("login failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
