package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error load .env")
	}

	DSN := os.Getenv("DSN")
	return gorm.Open(mysql.Open(DSN), &gorm.Config{})
}
