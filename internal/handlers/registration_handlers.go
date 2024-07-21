package handlers

import (
	"fmt"
	"lemmego/api"
	"lemmego/internal/inputs"
)

func RegistrationStoreHandler(ctx *api.Context) error {
	if body, err := ctx.ParseAndValidate(&inputs.RegistrationInput{}); err != nil {
		fmt.Printf("%s", err.Error())
		return err
	} else {
		input := body.(*inputs.RegistrationInput)
		fmt.Println(input)
		file, err := ctx.Upload("logo", "logos")
		if err != nil {
			return err
		}

		fmt.Println(file.Name())
	}

	return nil
}
