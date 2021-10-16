package database

import (
	"time"

	"gorm.io/gorm"
)

type Credential struct {
	CredentialId string `gorm:"primaryId"`
	ClientKey    string `gorm:"type:varchar(64)"`
	SecretKey    string `gorm:"type:varchar(64)"`
	Platform     string `gorm:"type:varchar(20)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}
