package services

import (
	"regexp"
	"strconv"
	"strings"
)

type OFXTransaction struct {
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
}

func ParseOFX(content string) ([]OFXTransaction, error) {
	var transactions []OFXTransaction

	stmttrnRegex := regexp.MustCompile(`(?is)<STMTTRN>(.*?)</STMTTRN>`)
	trnamtRegex := regexp.MustCompile(`(?i)<TRNAMT>([^<\r\n]+)`)
	dtpostRegex := regexp.MustCompile(`(?i)<DTPOST>([^<\r\n]+)`)
	memoRegex := regexp.MustCompile(`(?i)<MEMO>([^<\r\n]+)`)
	nameRegex := regexp.MustCompile(`(?i)<NAME>([^<\r\n]+)`)

	matches := stmttrnRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		block := match[1]

		var tx OFXTransaction

		amtMatch := trnamtRegex.FindStringSubmatch(block)
		if len(amtMatch) >= 2 {
			amtStr := strings.TrimSpace(amtMatch[1])
			amt, _ := strconv.ParseFloat(amtStr, 64)
			tx.Amount = amt
		}

		dtMatch := dtpostRegex.FindStringSubmatch(block)
		if len(dtMatch) >= 2 {
			tx.Date = strings.TrimSpace(dtMatch[1])
		}

		memoMatch := memoRegex.FindStringSubmatch(block)
		nameMatch := nameRegex.FindStringSubmatch(block)
		if len(nameMatch) >= 2 {
			tx.Description = strings.TrimSpace(nameMatch[1])
		} else if len(memoMatch) >= 2 {
			tx.Description = strings.TrimSpace(memoMatch[1])
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}
