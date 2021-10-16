package database

import (
	database "xapiens/database/migrations"

	"gorm.io/gorm"
)

func Migrations(db *gorm.DB) {
	db.AutoMigrate(&database.Credentials{})
}
