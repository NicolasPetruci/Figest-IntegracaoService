package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// GetPluggyToken returns a mock Pluggy Connect Token
func GetPluggyToken(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"accessToken": "mock-pluggy-token-123",
	})
}

// GetPluggyAccounts returns a mock list of connected bank accounts
func GetPluggyAccounts(c *fiber.Ctx) error {
	accounts := []fiber.Map{
		{
			"id": "acc-1",
			"name": "Nubank Conta Corrente",
			"type": "CHECKING",
			"balance": 1500.50,
		},
		{
			"id": "acc-2",
			"name": "Itaú Cartão de Crédito",
			"type": "CREDIT",
			"balance": -300.00,
		},
	}
	return c.JSON(fiber.Map{
		"results": accounts,
	})
}

// PluggyWebhook is a mock webhook endpoint that logs new transactions
func PluggyWebhook(c *fiber.Ctx) error {
	log.Println("Received Pluggy webhook: new transactions")
	return c.SendStatus(fiber.StatusOK)
}
