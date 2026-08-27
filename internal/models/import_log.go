package models

import "time"

type ImportLog struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	BatchID           string     `gorm:"index" json:"batchId"`
	UserID            string     `json:"userId"`
	FileName          string     `json:"fileName"`
	FileType          string     `json:"fileType"`
	AccountID         string     `json:"accountId"`
	Subtag            string     `json:"subtag"`
	Status            string     `json:"status"`
	TotalTransactions int        `json:"totalTransactions"`
	ProcessedAt       *time.Time `json:"processedAt"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}
