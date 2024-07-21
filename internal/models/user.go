package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	OrgId     uint   `json:"org_id"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Logo      string `json:"logo" gorm:"not null"`
	Email     string `json:"email" gorm:"unique not null"`
	Password  string `json:"password" gorm:"not null"`
}
