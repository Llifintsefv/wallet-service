package service

import (
	"context"
	"log/slog"
	"wallet-service/internal/repository/postgres"

	"github.com/google/uuid"
)

type Service interface {
	CreateWallet(ctx context.Context) (string, error)
}

type service struct {
	repo   postgres.Repository
	logger slog.Logger
}

func NewService(repo postgres.Repository, logger slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) CreateWallet(ctx context.Context) (string, error) {
	uuid := uuid.New().String()

	err := s.repo.CreateWallet(ctx, uuid)
	if err != nil {
		return "", err
	}

	return uuid, err
}
