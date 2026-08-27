package services

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type OFXTransaction struct {
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
}

func FormatOFXDate(rawDate string) string {
	rawDate = strings.TrimSpace(rawDate)
	if len(rawDate) >= 8 {
		year := rawDate[0:4]
		month := rawDate[4:6]
		day := rawDate[6:8]
		return year + "-" + month + "-" + day + "T12:00:00Z"
	}
	return time.Now().Format(time.RFC3339)
}

func ParseOFX(content string) ([]OFXTransaction, error) {
	var transactions []OFXTransaction

	// Split content by <STMTTRN> (case insensitive) to handle both SGML (OFX 1.x) and XML (OFX 2.x)
	blocks := regexp.MustCompile(`(?i)<STMTTRN>`).Split(content, -1)
	if len(blocks) <= 1 {
		return transactions, nil
	}

	trnamtRegex := regexp.MustCompile(`(?i)<TRNAMT>\s*([-\d.,]+)`)
	dtpostRegex := regexp.MustCompile(`(?i)<DTPOST(?:ED)?>\s*(\d{8})`)
	memoRegex := regexp.MustCompile(`(?i)<MEMO>\s*([^<\r\n]+)`)
	nameRegex := regexp.MustCompile(`(?i)<NAME>\s*([^<\r\n]+)`)

	for i := 1; i < len(blocks); i++ {
		block := blocks[i]
		if idx := strings.Index(strings.ToUpper(block), "</BANKTRANLIST>"); idx != -1 {
			block = block[:idx]
		}

		var tx OFXTransaction

		amtMatch := trnamtRegex.FindStringSubmatch(block)
		if len(amtMatch) >= 2 {
			amtStr := strings.TrimSpace(amtMatch[1])
			amtStr = strings.ReplaceAll(amtStr, ",", ".")
			amt, _ := strconv.ParseFloat(amtStr, 64)
			tx.Amount = amt
		}

		dtMatch := dtpostRegex.FindStringSubmatch(block)
		if len(dtMatch) >= 2 {
			tx.Date = strings.TrimSpace(dtMatch[1])
		}

		nameMatch := nameRegex.FindStringSubmatch(block)
		memoMatch := memoRegex.FindStringSubmatch(block)
		if len(nameMatch) >= 2 && strings.TrimSpace(nameMatch[1]) != "" {
			tx.Description = strings.TrimSpace(nameMatch[1])
		} else if len(memoMatch) >= 2 && strings.TrimSpace(memoMatch[1]) != "" {
			tx.Description = strings.TrimSpace(memoMatch[1])
		} else {
			tx.Description = "Lançamento OFX"
		}

		if tx.Amount != 0 || tx.Date != "" {
			transactions = append(transactions, tx)
		}
	}

	return transactions, nil
}
