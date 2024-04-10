package cmder

import (
	"errors"
)

var SnakeCase = func(input string) error {
	if len(input) == 0 {
		return errors.New("input cannot be empty")
	}

	// input must start with a letter and contain only letters, numbers, and underscores
	if len(input) > 0 && input[0] < 'a' || input[0] > 'z' {
		return errors.New("field name must start with a lowercase letter")
	}

	for _, c := range input {
		if c < 'a' || c > 'z' {
			if c < '0' || c > '9' {
				if c != '_' {
					return errors.New("field name must contain only lowercase letters, numbers, and underscores")
				}
			}
		}
	}
	return nil
}

var SnakeCaseEmptyAllowed = func(input string) error {
	if len(input) == 0 {
		return nil
	}

	// input must start with a letter and contain only letters, numbers, and underscores
	if len(input) > 0 && input[0] < 'a' || input[0] > 'z' {
		return errors.New("field name must start with a lowercase letter")
	}

	for _, c := range input {
		if c < 'a' || c > 'z' {
			if c < '0' || c > '9' {
				if c != '_' {
					return errors.New("field name must contain only lowercase letters, numbers, and underscores")
				}
			}
		}
	}
	return nil
}
