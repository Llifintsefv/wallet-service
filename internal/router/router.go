package router

import (
	"wallet-service/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(handler handler.Handler) *fiber.App {
	app := fiber.New()

	app.Post("api/v1/wallets", handler.CreateWallet)
	app.Post("api/v1/wallet", handler.Transaction)
	app.Get("api/v1/wallet/:uuid", handler.GetWallet)

	return app
}
