package models

import (
	"github.com/lemmego/lemmego/api/db"
)

type Org struct {
	db.Model
    OrgUsername string `json:"org_username" gorm:"unique"`
    OrgName string `json:"org_name" gorm:"not null"`
    OrgEmail string `json:"org_email" gorm:"unique not null"`
    OrgLogo string `json:"org_logo" gorm:"not null"`
}
