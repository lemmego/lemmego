package cmder

import (
	"errors"

	"github.com/manifoldco/promptui"
)

var SnakeCaseValidator = func(allowEmpty bool) promptui.ValidateFunc {
	return func(input string) error {
		if !allowEmpty && len(input) == 0 {
			return errors.New("Input cannot be empty")
		}

		if allowEmpty && len(input) == 0 {
			return nil
		}

		// input must start with a letter and contain only letters, numbers, and underscores
		if len(input) > 0 && input[0] < 'a' || input[0] > 'z' {
			return errors.New("Field Name must start with a lowercase letter")
		}

		for _, c := range input {
			if c < 'a' || c > 'z' {
				if c < '0' || c > '9' {
					if c != '_' {
						return errors.New("Field Name must contain only lowercase letters, numbers, and underscores")
					}
				}
			}
		}
		return nil
	}
}
