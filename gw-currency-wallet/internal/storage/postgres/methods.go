package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"gw-currency-wallet/internal/storage/models"
)

// --- USERS ---
func (s *PostgresStorage) CreateUser(ctx context.Context, user models.UserRegister) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO users (email, password) VALUES ($1, $2)`,
		user.Email, user.Password,
	)
	if err != nil {
		s.logger.Error("failed to insert user",
			zap.Error(err),
			zap.String("email", user.Email),
		)
	}
	return err
}

func (s *PostgresStorage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, email, password FROM users WHERE email=$1`, email,
	)
	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return u, models.ErrUserNotFound
		}
		s.logger.Error("failed to get user by email",
			zap.Error(err),
			zap.String("email", email),
		)
		return u, err
	}
	return u, nil
}

// --- WALLETS ---
func (s *PostgresStorage) GetBalance(ctx context.Context, userID uuid.UUID) (models.WalletResponse, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT currency, balance FROM wallets WHERE user_id=$1`, userID,
	)
	if err != nil {
		s.logger.Error("failed to get balance",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return models.WalletResponse{}, err
	}
	defer rows.Close()

	resp := models.WalletResponse{Balances: make(map[string]decimal.Decimal)}
	for rows.Next() {
		var currency string
		var balance decimal.Decimal
		if err := rows.Scan(&currency, &balance); err != nil {
			return resp, err
		}
		resp.Balances[currency] = balance
	}

	// Проверяем ошибки итерации
	if err := rows.Err(); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *PostgresStorage) Deposit(
	ctx context.Context,
	userID uuid.UUID,
	currency string,
	amount decimal.Decimal,
) (models.WalletResponse, uuid.UUID, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error("failed to start tx for deposit", zap.Error(err))
		return models.WalletResponse{}, uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO wallets (user_id, currency, balance)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, currency)
		 DO UPDATE SET balance = wallets.balance + $3`,
		userID, currency, amount,
	)
	if err != nil {
		s.logger.Error("failed to deposit",
			zap.Error(err),
			zap.String("user_id", userID.String()),
			zap.String("currency", currency),
			zap.String("amount", amount.String()),
		)
		return models.WalletResponse{}, uuid.Nil, err
	}

	transactionID := uuid.New()
	_, err = tx.Exec(ctx,
		`INSERT INTO transactions (id, user_id, type, currency, amount)
		 VALUES ($1, $2, $3, $4, $5)`,
		transactionID, userID, models.TransactionDeposit, currency, amount,
	)
	if err != nil {
		s.logger.Error("failed to insert deposit transaction",
			zap.Error(err),
			zap.String("transaction_id", transactionID.String()),
			zap.String("user_id", userID.String()),
		)
		return models.WalletResponse{}, uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("failed to commit deposit tx", zap.Error(err))
		return models.WalletResponse{}, uuid.Nil, err
	}

	balance, err := s.GetBalance(ctx, userID)
	return balance, transactionID, err
}

func (s *PostgresStorage) Withdraw(
	ctx context.Context,
	userID uuid.UUID,
	currency string,
	amount decimal.Decimal,
) (models.WalletResponse, uuid.UUID, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error("failed to start tx for withdraw", zap.Error(err))
		return models.WalletResponse{}, uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	cmdTag, err := tx.Exec(ctx,
		`UPDATE wallets SET balance = balance - $3
		 WHERE user_id=$1 AND currency=$2 AND balance >= $3`,
		userID, currency, amount,
	)
	if err != nil {
		s.logger.Error("failed to withdraw",
			zap.Error(err),
			zap.String("user_id", userID.String()),
			zap.String("currency", currency),
			zap.String("amount", amount.String()),
		)
		return models.WalletResponse{}, uuid.Nil, err
	}

	if cmdTag.RowsAffected() == 0 {
		return models.WalletResponse{}, uuid.Nil, errors.New("insufficient funds or invalid amount")
	}

	transactionID := uuid.New()
	_, err = tx.Exec(ctx,
		`INSERT INTO transactions (id, user_id, type, currency, amount)
		 VALUES ($1, $2, $3, $4, $5)`,
		transactionID, userID, models.TransactionWithdraw, currency, amount,
	)
	if err != nil {
		s.logger.Error("failed to insert withdraw transaction",
			zap.Error(err),
			zap.String("transaction_id", transactionID.String()),
			zap.String("user_id", userID.String()),
		)
		return models.WalletResponse{}, uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("failed to commit withdraw tx", zap.Error(err))
		return models.WalletResponse{}, uuid.Nil, err
	}

	balance, err := s.GetBalance(ctx, userID)
	return balance, transactionID, err
}

func (s *PostgresStorage) Exchange(
	ctx context.Context,
	userID uuid.UUID,
	fromCurrency string,
	toCurrency string,
	amount decimal.Decimal,
	exchangedAmount decimal.Decimal,
) (models.WalletResponse, uuid.UUID, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error("failed to start tx for exchange", zap.Error(err))
		return models.WalletResponse{}, uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	cmdTag, err := tx.Exec(ctx,
		`UPDATE wallets SET balance = balance - $3
		 WHERE user_id=$1 AND currency=$2 AND balance >= $3`,
		userID, fromCurrency, amount,
	)
	if err != nil {
		s.logger.Error("failed to debit in exchange",
			zap.Error(err),
			zap.String("user_id", userID.String()),
			zap.String("from_currency", fromCurrency),
			zap.String("amount", amount.String()),
		)
		return models.WalletResponse{}, uuid.Nil, err
	}
	if cmdTag.RowsAffected() == 0 {
		return models.WalletResponse{}, uuid.Nil, errors.New("insufficient funds or invalid amount")
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO wallets (user_id, currency, balance)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, currency)
		 DO UPDATE SET balance = wallets.balance + $3`,
		userID, toCurrency, exchangedAmount,
	)
	if err != nil {
		s.logger.Error("failed to credit in exchange",
			zap.Error(err),
			zap.String("user_id", userID.String()),
			zap.String("to_currency", toCurrency),
			zap.String("exchanged_amount", exchangedAmount.String()),
		)
		return models.WalletResponse{}, uuid.Nil, err
	}

	transactionID := uuid.New()
	_, err = tx.Exec(ctx,
		`INSERT INTO transactions (id, user_id, type, currency, amount, to_currency, exchanged_amount)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		transactionID, userID, models.TransactionExchange, fromCurrency, amount, toCurrency, exchangedAmount,
	)
	if err != nil {
		s.logger.Error("failed to insert exchange transaction",
			zap.Error(err),
			zap.String("transaction_id", transactionID.String()),
			zap.String("user_id", userID.String()),
		)
		return models.WalletResponse{}, uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("failed to commit exchange tx", zap.Error(err))
		return models.WalletResponse{}, uuid.Nil, err
	}

	balance, err := s.GetBalance(ctx, userID)
	return balance, transactionID, err
}
