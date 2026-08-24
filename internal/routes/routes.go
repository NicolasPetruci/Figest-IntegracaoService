package routes

import (
	"github.com/NicolasPetruci/Figest-IntegracaoService/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/health", handlers.HealthCheck)
	
	importGroup := app.Group("/import")
	importGroup.Post("/ofx", handlers.ImportOFX)
	importGroup.Post("/csv", handlers.ImportCSV)
	importGroup.Get("/history", handlers.GetImportHistory)
	importGroup.Get("/:id/status", handlers.GetImportStatus)

	pluggyGroup := app.Group("/pluggy")
	pluggyGroup.Get("/token", handlers.GetPluggyToken)
	pluggyGroup.Get("/accounts", handlers.GetPluggyAccounts)
	pluggyGroup.Post("/webhook", handlers.PluggyWebhook)
}
