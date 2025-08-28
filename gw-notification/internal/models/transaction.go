package models

import "time"

// TransactionHistory - запись о крупной транзакции
type TransactionHistory struct {
	TransactionID string    `bson:"transaction_id" json:"transaction_id"`
	UserID        string    `bson:"user_id" json:"user_id"`
	Operation     string    `bson:"operation" json:"operation"`
	Amount        string    `bson:"amount" json:"amount"`
	Currency      string    `bson:"currency" json:"currency"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
}
