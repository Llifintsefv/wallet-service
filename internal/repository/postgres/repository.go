package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"wallet-service/internal/model"

	"github.com/shopspring/decimal"
)

type Repository interface {
	CreateWallet(ctx context.Context, uuid string) error
	GetBalanceByUuid(ctx context.Context, uuid string) (decimal.Decimal, error)
	Transaction(ctx context.Context, uuid string, amount decimal.Decimal, op string) error
}
type repository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewRepository(db *sql.DB, logger *slog.Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

var (
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

func (r *repository) CreateWallet(ctx context.Context, uuid string) error {
	const query = `INSERT INTO wallets (id) VALUES ($1)`

	_, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	return nil

}

func (r *repository) GetBalanceByUuid(ctx context.Context, uuid string) (decimal.Decimal, error) {
	const query = `SELECT balance FROM wallets WHERE id = $1`

	var balance decimal.Decimal
	err := r.db.QueryRowContext(ctx, query, uuid).Scan(&balance)
	if err != nil {
		return decimal.Zero, err
	}

	return balance, nil

}

func (r *repository) Transaction(ctx context.Context, uuid string, amount decimal.Decimal, op string) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	var balance decimal.Decimal
	err = tx.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", uuid).Scan(&balance)
	if err != nil {
		return fmt.Errorf("get balance: %w", err)
	}

	if op == model.TransactionWithdraw && balance.LessThan(amount) {
		err = fmt.Errorf("balance is not enough")
		return err
	}

	if op == model.TransactionDeposit {
		balance = balance.Add(amount)
	} else {
		balance = balance.Sub(amount)
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = $1 WHERE id = $2", balance.StringFixed(2), uuid)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	return nil
}
