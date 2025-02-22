package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"wallet-service/internal/model"

	"log/slog"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateWallet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	repo := NewRepository(db, logger)

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO wallets").
			WithArgs("test-uuid").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateWallet(context.Background(), "test-uuid")
		assert.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO wallets").
			WithArgs("test-uuid").
			WillReturnError(sql.ErrConnDone)

		err := repo.CreateWallet(context.Background(), "test-uuid")
		assert.Error(t, err)
	})
}
func TestGetBalanceByUuid(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	repo := NewRepository(db, logger)

	t.Run("success", func(t *testing.T) {
		expectedBalance := decimal.NewFromFloat(100.0)
		mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1").
			WithArgs("test-uuid").
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow("100.00"))

		balance, err := repo.GetBalanceByUuid(context.Background(), "test-uuid")
		assert.NoError(t, err)
		assert.True(t, expectedBalance.Equal(balance), "expected balance %v, got %v", expectedBalance, balance)
	})

	t.Run("wallet not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1").
			WithArgs("test-uuid").
			WillReturnError(sql.ErrNoRows)

		balance, err := repo.GetBalanceByUuid(context.Background(), "test-uuid")
		assert.Error(t, err)
		assert.Equal(t, decimal.Zero, balance)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1").
			WithArgs("test-uuid").
			WillReturnError(sql.ErrConnDone)

		balance, err := repo.GetBalanceByUuid(context.Background(), "test-uuid")
		assert.Error(t, err)
		assert.Equal(t, decimal.Zero, balance)
	})
}
func TestTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	repo := NewRepository(db, logger)

	t.Run("deposit success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1 FOR UPDATE").
			WithArgs("test-uuid").
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow("100.00"))
		mock.ExpectExec("UPDATE wallets SET balance = \\$1 WHERE id = \\$2").
			WithArgs("200.00", "test-uuid").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Transaction(context.Background(), "test-uuid", decimal.NewFromFloat(100.0), model.TransactionDeposit)
		assert.NoError(t, err)
	})

	t.Run("withdraw success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1 FOR UPDATE").
			WithArgs("test-uuid").
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow("200.00"))
		mock.ExpectExec("UPDATE wallets SET balance = \\$1 WHERE id = \\$2").
			WithArgs("100.00", "test-uuid").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Transaction(context.Background(), "test-uuid", decimal.NewFromFloat(100.0), model.TransactionWithdraw)
		assert.NoError(t, err)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1 FOR UPDATE").
			WithArgs("test-uuid").
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow("50.00"))
		mock.ExpectRollback()

		err := repo.Transaction(context.Background(), "test-uuid", decimal.NewFromFloat(100.0), model.TransactionWithdraw)
		assert.Error(t, err)
		assert.Equal(t, "balance is not enough", err.Error())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1 FOR UPDATE").
			WithArgs("test-uuid").
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.Transaction(context.Background(), "test-uuid", decimal.NewFromFloat(100.0), model.TransactionDeposit)
		assert.Error(t, err)
	})

	t.Run("update error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT balance FROM wallets WHERE id = \\$1 FOR UPDATE").
			WithArgs("test-uuid").
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow("100.00"))
		mock.ExpectExec("UPDATE wallets SET balance = \\$1 WHERE id = \\$2").
			WithArgs("200.00", "test-uuid").
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := repo.Transaction(context.Background(), "test-uuid", decimal.NewFromFloat(100.0), model.TransactionDeposit)
		assert.Error(t, err)
	})
}
