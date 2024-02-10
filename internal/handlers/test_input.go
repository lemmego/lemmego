package handlers

import "github.com/invopop/validation"

type TestInput struct {
	Username string `json:"username" in:"form=username"`
	Password string `json:"password" in:"form=password"`
}

func (r *TestInput) Validate() error {
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
