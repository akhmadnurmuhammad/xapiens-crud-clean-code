package database

import (
	"xapiens/database/migrations"

	"gorm.io/gorm"
)

func Migrations(db *gorm.DB) {
	db.AutoMigrate(&migrations.Credential{})
}
