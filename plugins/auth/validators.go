package auth

import (
	"github.com/invopop/validation"
)

type LoginStoreRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginStoreRequest) Validate() error {
	return validation.Errors{
		"username": validation.Validate(
			r.Username,
			validation.Required,
		),
		"password": validation.Validate(
			r.Password,
			validation.Required,
		),
	}.Filter()
}
