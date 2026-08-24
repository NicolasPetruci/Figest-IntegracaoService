package database

import (
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/NicolasPetruci/Figest-IntegracaoService/internal/models"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database")
		return
	}
	DB = db
	
	err = db.AutoMigrate(&models.ImportLog{})
	if err != nil {
		log.Println("Migration failed")
	}
}
