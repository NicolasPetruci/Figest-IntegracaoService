package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

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

	financeUrl := os.Getenv("FINANCEIRO_SERVICE_URL")
	if financeUrl == "" {
		financeUrl = "http://financeiro-service:3002"
	}

	accountIDParam := c.FormValue("accountId")
	subtagParam := c.FormValue("subtag")

	batchID := fmt.Sprintf("ofx-batch-%d", time.Now().UnixNano())

	insertedCount := 0
	for _, tx := range transactions {
		txType := "EXPENSE"
		amt := tx.Amount
		if amt < 0 {
			amt = -amt
		} else {
			txType = "INCOME"
		}

		payload := map[string]interface{}{
			"type":          txType,
			"amount":        amt,
			"description":   tx.Description,
			"date":          services.FormatOFXDate(tx.Date),
			"importBatchId": batchID,
		}

		if accountIDParam != "" {
			payload["accountId"] = accountIDParam
		}
		if subtagParam != "" {
			payload["subtag"] = subtagParam
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			continue
		}

		req, err := http.NewRequest("POST", financeUrl+"/transactions", bytes.NewBuffer(jsonData))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("x-user-id", userID)
			client := &http.Client{}
			resp, err := client.Do(req)
			if err == nil && resp != nil {
				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					insertedCount++
				}
				resp.Body.Close()
			}
		}
	}

	now := time.Now()
	importLog := models.ImportLog{
		BatchID:           batchID,
		UserID:            userID,
		FileName:          file.Filename,
		FileType:          "OFX",
		AccountID:         accountIDParam,
		Subtag:            subtagParam,
		Status:            "DONE",
		TotalTransactions: insertedCount,
		ProcessedAt:       &now,
	}

	if database.DB != nil {
		database.DB.Create(&importLog)
	}

	return c.JSON(fiber.Map{
		"message":      "Import successful",
		"batchId":      batchID,
		"transactions": len(transactions),
		"inserted":     insertedCount,
	})
}

func GetImportHistory(c *fiber.Ctx) error {
	userID := c.Get("x-user-id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing x-user-id header"})
	}

	var logs []models.ImportLog
	if database.DB != nil {
		database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&logs)

		// Fix legacy logs without batchId
		for i := range logs {
			if logs[i].BatchID == "" {
				logs[i].BatchID = fmt.Sprintf("ofx-batch-legacy-%d", logs[i].ID)
				database.DB.Model(&logs[i]).Update("batch_id", logs[i].BatchID)
			}
		}
	}

	return c.JSON(logs)
}

func UpdateImportBatch(c *fiber.Ctx) error {
	userID := c.Get("x-user-id")
	batchID := c.Params("batchId")
	if userID == "" || batchID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing param"})
	}

	var body struct {
		AccountID string `json:"accountId"`
		Subtag    string `json:"subtag"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	logID, _ := strconv.Atoi(batchID)

	if database.DB != nil {
		database.DB.Model(&models.ImportLog{}).
			Where("(batch_id = ? OR id = ?) AND user_id = ?", batchID, logID, userID).
			Updates(map[string]interface{}{
				"account_id": body.AccountID,
				"subtag":     body.Subtag,
			})
	}

	financeUrl := os.Getenv("FINANCEIRO_SERVICE_URL")
	if financeUrl == "" {
		financeUrl = "http://financeiro-service:3002"
	}

	payload, _ := json.Marshal(body)
	req, err := http.NewRequest("PATCH", financeUrl+"/transactions/import-batch/"+batchID, bytes.NewBuffer(payload))
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-user-id", userID)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err == nil && resp != nil {
			resp.Body.Close()
		}
	}

	return c.JSON(fiber.Map{
		"message": "Import batch updated successfully",
		"batchId": batchID,
	})
}

func ImportCSV(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "ImportCSV"}) }
func GetImportStatus(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "GetImportStatus"}) }
