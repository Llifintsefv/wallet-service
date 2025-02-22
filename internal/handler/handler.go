package handler

import (
	"log/slog"
	"wallet-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service service.Service
	logger  slog.Logger
}

func NewHandler(service service.Service, logger slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) CreateWallet(c *fiber.Ctx) error {
	ctx := c.Context()

	uuid, err := h.service.CreateWallet(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"uuid": uuid})
}
