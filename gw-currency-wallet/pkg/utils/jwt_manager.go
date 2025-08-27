package utils

import (
	"errors"
	"gw-currency-wallet/internal/storage/models"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JWTManager struct {
	secret   string
	tokenTTL time.Duration
}

func NewJWTManager(secret string, tokenTTL int) *JWTManager {
	return &JWTManager{
		secret:   secret,
		tokenTTL: time.Duration(tokenTTL) * time.Second, // В секундах
	}
}

func (m *JWTManager) GenerateToken(userID uuid.UUID, email string) (string, error) {
	claims := models.Claims{
		UserID: userID.String(),
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenTTL)), // Срок действия
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

// ParseJWT парсит и проверяет токен
func (m *JWTManager) ParseToken(tokenString string) (*models.Claims, error) {
	secret := []byte(m.secret)
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что алгоритм подписи корректен
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	// Извлекаем claims
	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
