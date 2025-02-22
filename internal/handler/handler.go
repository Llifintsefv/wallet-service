package handler

import (
	"log/slog"
	"wallet-service/internal/model"
	"wallet-service/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

type Handler struct {
	service service.Service
	logger  *slog.Logger
}

func NewHandler(service service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) CreateWallet(c *fiber.Ctx) error {
	ctx := c.Context()

	uuid, err := h.service.CreateWallet(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create wallet"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"uuid": uuid})
}

func (h *Handler) Transaction(c *fiber.Ctx) error {
	ctx := c.Context()
	var req model.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("failed to parse request", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request format"})
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		h.logger.Warn("invalid amount format", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid amount format"})
	}

	transaction := model.Transaction{
		Uuid:          req.ValletId,
		OperationType: req.OperationType,
		Amount:        amount,
	}

	if err := model.ValidateTransaction(transaction); err != nil {
		h.logger.Warn("validation failed", slog.Any("error", err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.Transaction(ctx, transaction); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
}

func (h *Handler) GetWallet(c *fiber.Ctx) error {
	ctx := c.Context()
	uuid := c.Params("uuid")

	balance, err := h.service.GetBalanceByUuid(ctx, uuid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"balance": balance.String()})
}
