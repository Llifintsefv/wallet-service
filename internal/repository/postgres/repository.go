package postgres

import (
	"context"
	"database/sql"
	"log/slog"
)

type Repository interface {
	CreateWallet(ctx context.Context, uuid string) error
}

type repository struct {
	db     *sql.DB
	logger slog.Logger
}

func NewRepository(db *sql.DB, logger slog.Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (r *repository) CreateWallet(ctx context.Context, uuid string) error {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO wallets (id) VALUES ($1)")
	if err != nil {
		r.logger.Error("Error creating wallet", slog.Any("error", err))
		return err

	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, uuid)
	if err != nil {
		r.logger.Error("Error creating wallet", slog.Any("error", err))
		return err
	}

	return nil

}
