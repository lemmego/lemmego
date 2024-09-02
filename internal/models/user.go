package models

import (
	"github.com/lemmego/api/db"
)

type User struct {
	db.Model
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"-"`
	OrgId     uint   `json:"org_id" gorm:"not null"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username" gorm:"unique not null"`
	Bio       string `json:"bio"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
}
