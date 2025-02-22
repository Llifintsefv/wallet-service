package service

import (
	"context"
	"fmt"
	"log/slog"
	"wallet-service/internal/model"
	"wallet-service/internal/repository/postgres"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service interface {
	CreateWallet(ctx context.Context) (string, error)
	Transaction(ctx context.Context, transactionRequest model.Transaction) error
	GetBalanceByUuid(ctx context.Context, uuid string) (decimal.Decimal, error)
}

type service struct {
	repo   postgres.Repository
	logger *slog.Logger
}

func NewService(repo postgres.Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) CreateWallet(ctx context.Context) (string, error) {
	uuid := uuid.New().String()
	if err := s.repo.CreateWallet(ctx, uuid); err != nil {
		s.logger.Error("failed to create wallet", slog.Any("error", err))
		return "", fmt.Errorf("create wallet: %w", err)
	}
	return uuid, nil
}

func (s *service) Transaction(ctx context.Context, transactionRequest model.Transaction) error {

	err := s.repo.Transaction(ctx, transactionRequest.Uuid, transactionRequest.Amount, transactionRequest.OperationType)
	if err != nil {
		return err
	}

	return nil

}

func (s *service) GetBalanceByUuid(ctx context.Context, uuid string) (decimal.Decimal, error) {

	balance, err := s.repo.GetBalanceByUuid(ctx, uuid)
	if err != nil {
		return decimal.Zero, err
	}

	return balance, nil
}
