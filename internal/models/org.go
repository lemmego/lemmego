package models

import (
	"gorm.io/gorm"
)

type Org struct {
	gorm.Model
	OrgUsername string `json:"org_username" gorm:"unique "`
	OrgName     string `json:"org_name"`
	OrgEmail    string `json:"org_email" gorm:"unique "`
}
