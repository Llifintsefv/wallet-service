package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"wallet-service/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type MockService struct {
	CreateWalletFn     func(ctx context.Context) (string, error)
	TransactionFn      func(ctx context.Context, transaction model.Transaction) error
	GetBalanceByUuidFn func(ctx context.Context, uuid string) (decimal.Decimal, error)
}

func (m *MockService) CreateWallet(ctx context.Context) (string, error) {
	return m.CreateWalletFn(ctx)
}

func (m *MockService) Transaction(ctx context.Context, transaction model.Transaction) error {
	return m.TransactionFn(ctx, transaction)
}

func (m *MockService) GetBalanceByUuid(ctx context.Context, uuid string) (decimal.Decimal, error) {
	return m.GetBalanceByUuidFn(ctx, uuid)
}

func TestCreateWallet(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	t.Run("Success", func(t *testing.T) {
		mockService := &MockService{
			CreateWalletFn: func(ctx context.Context) (string, error) {
				return "test-uuid", nil
			},
		}
		h := NewHandler(mockService, logger)

		app := fiber.New()
		app.Post("/wallets", h.CreateWallet)

		req := httptest.NewRequest(http.MethodPost, "/wallets", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "test-uuid", body["uuid"])
	})

	t.Run("Failed to create wallet", func(t *testing.T) {
		mockService := &MockService{
			CreateWalletFn: func(ctx context.Context) (string, error) {
				return "", errors.New("service error")
			},
		}
		h := NewHandler(mockService, logger)

		app := fiber.New()
		app.Post("/wallets", h.CreateWallet)

		req := httptest.NewRequest(http.MethodPost, "/wallets", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "failed to create wallet", body["error"])
	})
}

func TestTransaction(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	t.Run("Success Deposit", func(t *testing.T) {
		mockService := &MockService{
			TransactionFn: func(ctx context.Context, transaction model.Transaction) error {
				return nil
			},
		}
		h := NewHandler(mockService, logger)

		app := fiber.New()
		app.Post("/transaction", h.Transaction)

		reqBody := `{"valletId": "test-uuid", "operationType": "DEPOSIT", "amount": "100"}`
		req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "OK", body["message"])
	})

	t.Run("Success Withdraw", func(t *testing.T) {
		mockService := &MockService{
			TransactionFn: func(ctx context.Context, transaction model.Transaction) error {
				return nil
			},
		}
		h := NewHandler(mockService, logger)

		app := fiber.New()
		app.Post("/transaction", h.Transaction)

		reqBody := `{"valletId": "test-uuid", "operationType": "WITHDRAW", "amount": "50"}`
		req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "OK", body["message"])
	})

	t.Run(" invalid operation type", func(t *testing.T) {
		mockService := &MockService{}
		h := NewHandler(mockService, logger)

		app := fiber.New()
		app.Post("/transaction", h.Transaction)

		reqBody := `{"valletId": "test-uuid", "operationType": "INVALID", "amount": "100"}`
		req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Contains(t, body["error"], "invalid operation type")
	})

	t.Run("negative amount", func(t *testing.T) {
		mockService := &MockService{}
		h := NewHandler(mockService, logger)

		app := fiber.New()
		app.Post("/transaction", h.Transaction)

		reqBody := `{"valletId": "test-uuid", "operationType": "DEPOSIT", "amount": "-10"}`
		req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Contains(t, body["error"], "amount must be positive")
	})

	t.Run("Service transaction error", func(t *testing.T) {
		mockService := &MockService{
			TransactionFn: func(ctx context.Context, transaction model.Transaction) error {
				return errors.New("service transaction error")
			},
		}
		h := NewHandler(mockService, logger)

		app := fiber.New()
		app.Post("/transaction", h.Transaction)

		reqBody := `{"valletId": "test-uuid", "operationType": "DEPOSIT", "amount": "100"}`
		req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Contains(t, body["error"], "service transaction error")
	})
}

func TestGetWallet(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	t.Run("Success", func(t *testing.T) {
		mockService := &MockService{
			GetBalanceByUuidFn: func(ctx context.Context, uuid string) (decimal.Decimal, error) {
				return decimal.NewFromInt(100), nil
			},
		}
		h := NewHandler(mockService, logger)

		app := fiber.New()
		app.Get("/wallet/:uuid", h.GetWallet)

		req := httptest.NewRequest(http.MethodGet, "/wallet/test-uuid", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "100", body["balance"])
	})

	t.Run("Service error", func(t *testing.T) {
		mockService := &MockService{
			GetBalanceByUuidFn: func(ctx context.Context, uuid string) (decimal.Decimal, error) {
				return decimal.Zero, errors.New("service error")
			},
		}
		h := NewHandler(mockService, logger)

		app := fiber.New()
		app.Get("/wallet/:uuid", h.GetWallet)

		req := httptest.NewRequest(http.MethodGet, "/wallet/test-uuid", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Contains(t, body["error"], "service error")
	})
}
