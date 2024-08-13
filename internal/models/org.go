package models

import (
	"github.com/lemmego/lemmego/api/db"
)

type Org struct {
	db.Model
    OrgUsername string `json:"org_username" gorm:"unique"`
    OrgName string `json:"org_name"`
    OrgEmail string `json:"org_email" gorm:"unique"`
}
