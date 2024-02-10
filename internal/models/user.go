package models

type User struct {
	ID        int64  `json:"id" db:"id,omitempty"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name,omitempty"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
	OrgID     int64  `json:"org_id" db:"org_id,omitempty"`
}
