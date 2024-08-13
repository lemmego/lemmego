package models

import (
	"github.com/lemmego/lemmego/api/db"
)

type User struct {
	db.Model
    OrgId uint `json:"org_id"`
    Org Org `json:"org"`
    FirstName string `json:"first_name" gorm:"not null"`
    LastName string `json:"last_name" gorm:"not null"`
    Logo string `json:"logo"`
    Email string `json:"email" gorm:"unique not null"`
    Password string `json:"password" gorm:"not null"`
}
