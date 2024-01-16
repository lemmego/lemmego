package auth

import (
	"pressebo/framework/validator"

	"github.com/invopop/validation"
)

type LoginStoreRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistrationStoreRequest struct {
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
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

func (r *RegistrationStoreRequest) Validate() error {
	return validation.Errors{
		"username": validation.Validate(
			r.Username,
			validation.Required,
		),
		"password": validation.Validate(
			r.Password,
			validation.Required,
		),
		"password_confirmation": validation.Validate(
			r.PasswordConfirmation,
			validation.Required,
			validation.By(validator.StringEquals(r.Password, "Passwords do not match")),
		),
	}.Filter()
}
