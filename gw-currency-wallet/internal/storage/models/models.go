package models

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	TransactionDeposit  = "deposit"
	TransactionWithdraw = "withdraw"
	TransactionExchange = "exchange"
)

var ErrUserNotFound = fmt.Errorf("user not found")

// --- User models ---

type User struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"` // хранится хэш пароля
}

type UserRegister struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// --- Wallet models ---

// WalletResponse — ответ с балансами
type WalletResponse struct {
	Balances map[string]decimal.Decimal `json:"balances"`
}

// WalletTransaction — запрос на пополнение или снятие средств
type WalletTransaction struct {
	Currency string `json:"currency" binding:"required,oneof=USD EUR RUB"`
	Amount   string `json:"amount" binding:"required"` // decimal string
}

// --- Exchange models ---

type ExchangeRequest struct {
	FromCurrency string `json:"from_currency" binding:"required,oneof=USD EUR RUB"`
	ToCurrency   string `json:"to_currency" binding:"required,oneof=USD EUR RUB"`
	Amount       string `json:"amount" binding:"required"` // decimal string
}

// --- JWT Claims ---

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
