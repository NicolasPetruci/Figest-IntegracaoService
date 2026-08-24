package main

import (
	"log"
	"github.com/NicolasPetruci/Figest-IntegracaoService/internal/config"
	"github.com/NicolasPetruci/Figest-IntegracaoService/internal/database"
	"github.com/NicolasPetruci/Figest-IntegracaoService/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()

	app := fiber.New()
	routes.RegisterRoutes(app)

	log.Println("Starting server on port 3005...")
	app.Listen(":3005")
}

// Refatorado para melhor legibilidade
