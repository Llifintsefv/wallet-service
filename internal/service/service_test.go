package service

import (
	"context"
	"errors"
	"io"
	"testing"
	"wallet-service/internal/model"

	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type Repository interface {
	CreateWallet(ctx context.Context, uuid string) error
	Transaction(ctx context.Context, uuid string, amount decimal.Decimal, operationType string) error
	GetBalanceByUuid(ctx context.Context, uuid string) (decimal.Decimal, error)
}

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateWallet(ctx context.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *mockRepository) Transaction(ctx context.Context, uuid string, amount decimal.Decimal, operationType string) error {
	args := m.Called(ctx, uuid, amount, operationType)
	return args.Error(0)
}

func (m *mockRepository) GetBalanceByUuid(ctx context.Context, uuid string) (decimal.Decimal, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func TestCreateWallet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockRepository)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		service := NewService(mockRepo, logger)
		ctx := context.Background()

		var calledUUID string
		mockRepo.On("CreateWallet", ctx, mock.AnythingOfType("string")).Run(func(args mock.Arguments) {
			calledUUID = args.String(1)
		}).Return(nil)

		uuidStr, err := service.CreateWallet(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, uuidStr)
		require.Equal(t, calledUUID, uuidStr)

		_, err = uuid.Parse(uuidStr)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mockRepository)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		service := NewService(mockRepo, logger)
		ctx := context.Background()

		expectedErr := errors.New("database error")
		mockRepo.On("CreateWallet", ctx, mock.AnythingOfType("string")).Return(expectedErr)

		uuidStr, err := service.CreateWallet(ctx)
		require.Error(t, err)
		require.Empty(t, uuidStr)
		require.Contains(t, err.Error(), "create wallet: database error")
	})
}

func TestTransaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockRepository)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		service := NewService(mockRepo, logger)
		ctx := context.Background()

		transactionRequest := model.Transaction{
			Uuid:          "some-uuid",
			Amount:        decimal.NewFromFloat(100.0),
			OperationType: "credit",
		}

		mockRepo.On("Transaction", ctx, transactionRequest.Uuid, transactionRequest.Amount, transactionRequest.OperationType).Return(nil)

		err := service.Transaction(ctx, transactionRequest)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mockRepository)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		service := NewService(mockRepo, logger)
		ctx := context.Background()

		transactionRequest := model.Transaction{
			Uuid:          "some-uuid",
			Amount:        decimal.NewFromFloat(100.0),
			OperationType: "credit",
		}

		expectedErr := errors.New("insufficient funds")
		mockRepo.On("Transaction", ctx, transactionRequest.Uuid, transactionRequest.Amount, transactionRequest.OperationType).Return(expectedErr)

		err := service.Transaction(ctx, transactionRequest)
		require.Error(t, err)
		require.Equal(t, expectedErr, err)
	})
}

func TestGetBalanceByUuid(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockRepository)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		service := NewService(mockRepo, logger)
		ctx := context.Background()

		uuidStr := "some-uuid"
		expectedBalance := decimal.NewFromFloat(100.0)
		mockRepo.On("GetBalanceByUuid", ctx, uuidStr).Return(expectedBalance, nil)

		balance, err := service.GetBalanceByUuid(ctx, uuidStr)
		require.NoError(t, err)
		require.Equal(t, expectedBalance, balance)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mockRepository)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		service := NewService(mockRepo, logger)
		ctx := context.Background()

		uuidStr := "some-uuid"
		expectedErr := errors.New("wallet not found")
		mockRepo.On("GetBalanceByUuid", ctx, uuidStr).Return(decimal.Zero, expectedErr)

		balance, err := service.GetBalanceByUuid(ctx, uuidStr)
		require.Error(t, err)
		require.Equal(t, expectedErr, err)
		require.True(t, balance.Equal(decimal.Zero))
	})
}
