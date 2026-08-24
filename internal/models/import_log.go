package models

import "time"

type ImportLog struct {
	ID                uint       `gorm:"primaryKey"`
	UserID            string
	FileName          string
	FileType          string
	Status            string
	TotalTransactions int
	ProcessedAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
