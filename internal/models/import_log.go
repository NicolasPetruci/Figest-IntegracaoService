package models

import "time"

type ImportLog struct {
	ID                uint       `gorm:"primaryKey"`
	BatchID           string     `gorm:"index"`
	UserID            string
	FileName          string
	FileType          string
	AccountID         string
	Subtag            string
	Status            string
	TotalTransactions int
	ProcessedAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
