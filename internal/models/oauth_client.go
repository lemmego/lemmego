package models

import (
	"gorm.io/gorm"
	"time"
)

type OauthClient struct {
	ID          string         `json:"id" gorm:"unique not null"`
	Secret      string         `json:"secret" gorm:"not null"`
	RedirectUri string         `json:"redirect_uri" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}
