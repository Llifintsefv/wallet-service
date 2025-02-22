package router

import (
	"wallet-service/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(handler handler.Handler) *fiber.App {
	app := fiber.New()

	app.Post("api/wallet", handler.CreateWallet)

	return app
}
