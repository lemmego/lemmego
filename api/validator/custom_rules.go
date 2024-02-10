package validator

import (
	"errors"

	"github.com/invopop/validation"
)

func StringEquals(str string, msg string) validation.RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(string)
		if s != str {
			if msg != "" {
				return errors.New(msg)
			}
			return errors.New("does not match")
		}
		return nil
	}
}
