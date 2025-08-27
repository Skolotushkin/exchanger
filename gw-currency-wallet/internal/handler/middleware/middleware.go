package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"gw-currency-wallet/pkg/utils"
)

// GinZapLogger логирует все запросы
func GinZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Next()
		logger.Info("request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

// RecoverMiddleware обрабатывает паники и возвращает 500
func RecoverMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", zap.Any("error", r))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}

// AuthMiddleware проверяет JWT и кладет userID и email в контекст
func AuthMiddleware(jwtManager *utils.JWTManager, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			logger.Warn("invalid JWT", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		uid, err := uuid.Parse(claims.UserID)
		if err != nil {
			logger.Warn("invalid UUID in token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
			return
		}

		// Кладём в контекст
		c.Set("userID", uid)
		c.Set("email", claims.Email)
		c.Next()
	}
}
