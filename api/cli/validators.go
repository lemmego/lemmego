package cli

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/lemmego/lemmego/api/cmder"
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

var NotIn = func(ignoreList []string, message string, validators ...cmder.ValidateFunc) cmder.ValidateFunc {
	return func(input string) error {
		if slices.Contains(ignoreList, strings.ToLower(input)) {
			if message == "" {
				return errors.New(fmt.Sprintf("input must not contain %s", strings.Join(ignoreList, ",")))
			}
			return errors.New(message)
		}
		for _, v := range validators {
			if err := v(input); err != nil {
				return err
			}
		}
		return nil
	}
}

var In = func(allowList []string, message string, validators ...cmder.ValidateFunc) cmder.ValidateFunc {
	return func(input string) error {
		if slices.Contains(allowList, strings.ToLower(input)) {
			if message == "" {
				return errors.New(fmt.Sprintf("input must contain %s", strings.Join(allowList, ",")))
			}
			return errors.New(message)
		}
		for _, v := range validators {
			if err := v(input); err != nil {
				return err
			}
		}
		return nil
	}
}
