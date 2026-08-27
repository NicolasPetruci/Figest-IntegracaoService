package database

import (
	"log"
	"os"

	"github.com/NicolasPetruci/Figest-IntegracaoService/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=db user=figest password=password dbname=figestdb port=5432 sslmode=disable TimeZone=UTC"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	DB = db

	err = db.AutoMigrate(&models.ImportLog{})
	if err != nil {
		log.Printf("Migration failed: %v", err)
	}
}
