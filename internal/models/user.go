package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	OrgId       uint   `json:"org_id" `
	FirstName   string `json:"first_name" `
	LastName    string `json:"last_name" `
	Email       string `json:"email" gorm:"not null,unique"`
	Password    string `json:"password" gorm:"not null"`
	OrgUsername string `json:"org_username" gorm:"not null,unique"`
	OrgName     string `json:"org_name" gorm:"not null"`
}
