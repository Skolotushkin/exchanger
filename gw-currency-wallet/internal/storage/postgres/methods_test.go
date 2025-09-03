package postgres_test

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"gw-currency-wallet/internal/storage/models"
	"gw-currency-wallet/internal/storage/postgres"
)

func newMockStorage(t *testing.T) (pgxmock.PgxPoolIface, *postgres.PostgresStorage) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock: %v", err)
	}
	logger := zap.NewNop() // глушилка логов
	st := postgres.NewPostgresStorageWithPool(mock, logger)
	return mock, st
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	mock, st := newMockStorage(t)
	defer mock.Close()

	user := models.UserRegister{Email: "test@example.com", Password: "hashedpwd"}

	// успешная вставка
	mock.ExpectExec(`INSERT INTO users`).
		WithArgs(user.Email, user.Password).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := st.CreateUser(ctx, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// дубликат пользователя
	mock.ExpectExec(`INSERT INTO users`).
		WithArgs(user.Email, user.Password).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))

	err = st.CreateUser(ctx, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate")
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()
	mock, st := newMockStorage(t)
	defer mock.Close()

	uid := uuid.New()
	expected := models.User{ID: uid, Email: "test@example.com", Password: "hash"}

	// пользователь найден
	rows := pgxmock.NewRows([]string{"id", "email", "password"}).
		AddRow(expected.ID, expected.Email, expected.Password)
	mock.ExpectQuery(`SELECT id, email, password FROM users WHERE email`).
		WithArgs(expected.Email).
		WillReturnRows(rows)

	user, err := st.GetUserByEmail(ctx, expected.Email)
	assert.NoError(t, err)
	assert.Equal(t, expected.Email, user.Email)

	// не найден
	mock.ExpectQuery(`SELECT id, email, password FROM users WHERE email`).
		WithArgs("missing@example.com").
		WillReturnError(models.ErrUserNotFound)

	_, err = st.GetUserByEmail(ctx, "missing@example.com")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrUserNotFound))
}

func TestDeposit(t *testing.T) {
	ctx := context.Background()
	mock, st := newMockStorage(t)
	defer mock.Close()

	uid := uuid.New()
	currency := "USD"
	amount := decimal.NewFromInt(100)

	// успешный депозит
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO wallets`).
		WithArgs(uid, currency, amount).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	// ожидание вставки транзакции (UUID транзакции - любой аргумент)
	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(pgxmock.AnyArg(), uid, models.TransactionDeposit, currency, amount).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()
	// ожидание GetBalance
	rows := pgxmock.NewRows([]string{"currency", "balance"}).AddRow(currency, amount)
	mock.ExpectQuery(`SELECT currency, balance FROM wallets WHERE user_id`).
		WithArgs(uid).
		WillReturnRows(rows)

	_, _, err := st.Deposit(ctx, uid, currency, amount)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// ошибка — невалидная валюта
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO wallets`).
		WithArgs(uid, "BTC", amount).
		WillReturnError(errors.New("invalid currency"))
	mock.ExpectRollback()

	_, _, err = st.Deposit(ctx, uid, "BTC", amount)
	assert.Error(t, err)
}

func TestWithdraw(t *testing.T) {
	ctx := context.Background()
	mock, st := newMockStorage(t)
	defer mock.Close()

	uid := uuid.New()
	currency := "USD"
	amount := decimal.NewFromInt(50)

	// успешный withdraw
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE wallets SET balance = balance -`).
		WithArgs(uid, currency, amount).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	// транзакция
	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(pgxmock.AnyArg(), uid, models.TransactionWithdraw, currency, amount).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()
	// GetBalance
	rows := pgxmock.NewRows([]string{"currency", "balance"}).AddRow(currency, decimal.NewFromInt(0))
	mock.ExpectQuery(`SELECT currency, balance FROM wallets WHERE user_id`).
		WithArgs(uid).
		WillReturnRows(rows)

	_, _, err := st.Withdraw(ctx, uid, currency, amount)
	assert.NoError(t, err)

	// недостаточно средств
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE wallets SET balance = balance -`).
		WithArgs(uid, currency, amount).
		WillReturnError(errors.New("insufficient funds"))
	mock.ExpectRollback()

	_, _, err = st.Withdraw(ctx, uid, currency, amount)
	assert.Error(t, err)
}

func TestExchange(t *testing.T) {
	ctx := context.Background()
	mock, st := newMockStorage(t)
	defer mock.Close()

	uid := uuid.New()
	from, to := "USD", "EUR"
	amount := decimal.NewFromInt(100)
	exchanged := decimal.NewFromInt(85)

	// успешный обмен
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE wallets SET balance = balance -`).
		WithArgs(uid, from, amount).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec(`INSERT INTO wallets`).
		WithArgs(uid, to, exchanged).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(pgxmock.AnyArg(), uid, models.TransactionExchange, from, amount, to, exchanged).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()
	// GetBalance
	rows := pgxmock.NewRows([]string{"currency", "balance"}).AddRow(to, exchanged)
	mock.ExpectQuery(`SELECT currency, balance FROM wallets WHERE user_id`).
		WithArgs(uid).
		WillReturnRows(rows)

	_, _, err := st.Exchange(ctx, uid, from, to, amount, exchanged)
	assert.NoError(t, err)

	// ошибка — недостаточно средств
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE wallets SET balance = balance -`).
		WithArgs(uid, from, amount).
		WillReturnError(errors.New("insufficient funds"))
	mock.ExpectRollback()

	_, _, err = st.Exchange(ctx, uid, from, to, amount, exchanged)
	assert.Error(t, err)
}
