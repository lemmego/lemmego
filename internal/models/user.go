package models

import (
	"encoding/gob"
	"strconv"
	"time"
)

func init() {
	gob.Register(&User{})
}

type User struct {
	ID        uint64     `json:"id" db:"id,omitempty"`
	Email     string     `json:"email" db:"email"`
	Name      string     `json:"name" db:"name"`
	Password  string     `json:"-" db:"password"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) GetID() string {
	return strconv.Itoa(int(u.ID))
}

func (u *User) GetUsername() string {
	return u.Email
}

func (u *User) GetPassword() string {
	return u.Password
}
