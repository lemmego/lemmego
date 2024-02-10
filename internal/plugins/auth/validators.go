package auth

import (
	// "mime/multipart"
	// "pressebo/framework"
	"pressebo/api/validator"

	"github.com/invopop/validation"
)

// type LoginStoreRequest struct {
// 	Username string           `json:"username"`
// 	Password string           `json:"password"`
// 	Logos    []multipart.File `json:"logos"`
// }

type LoginStoreInput struct {
	Username string `json:"username" in:"form=username"`
	Password string `json:"password" in:"form=password"`
}

type RegistrationStoreInput struct {
	FirstName            string `json:"first_name" in:"form=first_name"`
	LastName             string `json:"last_name" in:"form=last_name"`
	Username             string `json:"username" in:"form=username"`
	Password             string `json:"password" in:"form=password"`
	PasswordConfirmation string `json:"password_confirmation" in:"form=password_confirmation"`
}

// type RegistrationStoreRequest struct {
// 	FirstName            string `json:"firstName"`
// 	LastName             string `json:"lastName"`
// 	Username             string `json:"username"`
// 	Password             string `json:"password"`
// 	PasswordConfirmation string `json:"password_confirmation"`
// }

func (r *LoginStoreInput) Validate() error {
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

func (r *RegistrationStoreInput) Validate() error {
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
