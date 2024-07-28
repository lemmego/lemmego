package handlers

import (
	"fmt"
	"lemmego/api"
	"lemmego/internal/inputs"
)

func RegistrationStoreHandler(ctx *api.Context) error {
	fmt.Println("RegistrationStoreHandler")
	body := &inputs.RegistrationInput{}
	if err := ctx.Validate(body); err != nil {
		return err
	}
	return nil
}
