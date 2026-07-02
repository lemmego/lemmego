package models

import (
	"context"
	"encoding/gob"
	"strconv"
	"time"

	"github.com/lemmego/api/utils"
)

func init() {
	gob.Register(&User{})
}

type User struct {
	ID        uint64 `json:"id" db:"id,omitempty"`
	Email     string `json:"email" db:"email"`
	Name      string `json:"name" db:"name"`
	Password  string `json:"-" db:"password"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
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

func (u *User) BeforeCreate(ctx context.Context) error {
	if hashed, err := utils.Bcrypt(u.Password); err != nil {
		return err
	} else {
		u.Password = hashed
	}

	if u.CreatedAt == "" {
		u.CreatedAt = time.Now().Format(time.RFC3339)
	}

	if u.UpdatedAt == "" {
		u.UpdatedAt = time.Now().Format(time.RFC3339)
	}

	return nil
}
