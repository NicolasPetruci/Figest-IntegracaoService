package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/NicolasPetruci/Figest-IntegracaoService/internal/database"
	"github.com/NicolasPetruci/Figest-IntegracaoService/internal/models"
	"github.com/NicolasPetruci/Figest-IntegracaoService/internal/services"
	"github.com/gofiber/fiber/v2"
)

func ImportOFX(c *fiber.Ctx) error {
	userID := c.Get("x-user-id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing x-user-id header"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File is required"})
	}

	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot open file"})
	}
	defer fileContent.Close()

	contentBytes, err := io.ReadAll(fileContent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot read file"})
	}

	transactions, err := services.ParseOFX(string(contentBytes))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse OFX"})
	}

	for _, tx := range transactions {
		jsonData, err := json.Marshal(tx)
		if err != nil {
			continue
		}

		req, err := http.NewRequest("POST", "http://localhost:3000/api/finance/transactions", bytes.NewBuffer(jsonData))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("x-user-id", userID)
			client := &http.Client{}
			resp, _ := client.Do(req)
			if resp != nil {
				resp.Body.Close()
			}
		}
	}

	importLog := models.ImportLog{
		UserID:            userID,
		FileName:          file.Filename,
		FileType:          "OFX",
		Status:            "DONE",
		TotalTransactions: len(transactions),
	}

	if result := database.DB.Create(&importLog); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save import log"})
	}

	return c.JSON(fiber.Map{
		"message":      "Import successful",
		"transactions": len(transactions),
	})
}
func ImportCSV(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "ImportCSV"}) }
func GetImportHistory(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "GetImportHistory"}) }
func GetImportStatus(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "GetImportStatus"}) }
