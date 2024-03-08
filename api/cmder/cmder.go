package cmder

import (
	"reflect"

	"github.com/manifoldco/promptui"
)

type PromptResultType int

const (
	PromptResultTypeNormal PromptResultType = iota
	PromptResultTypeBoolean
	PromptResultTypeRecurring
	PromptResultTypeSelect
)

type PromptResult struct {
	Type          PromptResultType
	ShouldAskNext bool
	Result        interface{}
	Error         error
}

type Prompter interface {
	Ask(question string, validator promptui.ValidateFunc) Prompter
	AskBoolean(question string) Prompter
	AskRecurring(question string, validator promptui.ValidateFunc, prompts ...func(result any) Prompter) Prompter
	Select(label string, items []string) Prompter
	If(cb func(result interface{}) bool) Prompter
	Then() Prompter
	Fill(ptr any) Prompter
}

func (pr *PromptResult) Then() Prompter {
	pr.ShouldAskNext = true
	return pr
}

func (pr *PromptResult) Fill(ptr any) Prompter {
	if pr.ShouldAskNext {
		if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
			panic("Fill() must be called with a pointer")
		}

		// if reflect.TypeOf(ptr).Elem().Kind() != reflect.TypeOf(pr.Result).Kind() {
		// 	panic("Fill() must be called with a pointer of the same type as the result")
		// }

		// Set the result to the pointer
		reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(pr.Result))
	}
	return pr
}

func (pr *PromptResult) Ask(question string, validator promptui.ValidateFunc) Prompter {
	if pr.ShouldAskNext {
		return Ask(question, validator)
	}
	return pr
}

func (pr *PromptResult) AskBoolean(question string) Prompter {
	if pr.ShouldAskNext {
		return AskBoolean(question)
	}
	return pr
}

func (pr *PromptResult) AskRecurring(question string, validator promptui.ValidateFunc, prompts ...func(result any) Prompter) Prompter {
	if pr.ShouldAskNext {
		return AskRecurring(question, validator, prompts...)
	}
	return pr
}

func (pr *PromptResult) Select(label string, items []string) Prompter {
	if pr.ShouldAskNext {
		return Select(label, items)
	}
	return pr
}

func (pr *PromptResult) If(cb func(result interface{}) bool) Prompter {
	if pr.ShouldAskNext {
		if cb(pr.Result) == true {
			pr.ShouldAskNext = true
		} else {
			pr.ShouldAskNext = false
		}
	}
	return pr
}

func Ask(question string, validator promptui.ValidateFunc) Prompter {
	if validator == nil {
		validator = func(input string) error {
			return nil
		}
	}

	prompt := promptui.Prompt{
		Label:    question,
		Validate: validator,
	}

	res, err := prompt.Run()
	if err != nil {
		return &PromptResult{Type: PromptResultTypeNormal, ShouldAskNext: false, Result: nil, Error: err}
	}
	return &PromptResult{Type: PromptResultTypeNormal, ShouldAskNext: true, Result: res, Error: nil}
}

func AskBoolean(question string) Prompter {
	q := promptui.Prompt{
		Label: question + " (y/n)",
	}

	res, err := q.Run()
	if err != nil {
		return &PromptResult{Type: PromptResultTypeBoolean, ShouldAskNext: false, Result: false, Error: err}
	}

	return &PromptResult{Type: PromptResultTypeBoolean, ShouldAskNext: true, Result: res == "y" || res == "Y", Error: nil}
}

func Select(label string, items []string) Prompter {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return &PromptResult{Type: PromptResultTypeSelect, ShouldAskNext: false, Result: nil, Error: err}
	}

	return &PromptResult{Type: PromptResultTypeSelect, ShouldAskNext: true, Result: result, Error: nil}
}

func AskRecurring(question string, validator promptui.ValidateFunc, prompts ...func(result any) Prompter) Prompter {
	if validator == nil {
		validator = func(input string) error {
			return nil
		}
	}

	inputsFinished := false
	inputs := []string{}

	for !inputsFinished {
		prompt := promptui.Prompt{
			Label:    question + " (press enter when finished)",
			Validate: validator,
		}

		input, err := prompt.Run()

		if err != nil {
			return &PromptResult{Type: PromptResultTypeRecurring, ShouldAskNext: false, Result: nil, Error: err}
		}

		if input == "" {
			inputsFinished = true
			break
		}

		for _, p := range prompts {
			p(input)
		}

		inputs = append(inputs, input)

	}

	return &PromptResult{Type: PromptResultTypeRecurring, ShouldAskNext: true, Result: inputs, Error: nil}
}

// API I want to achieve:

// cmder.AskBoolean("Should your users belong to an organization?").
//		IfTrue().
//      AskIndefinitely("Enter the field name in snake_case", orgFields).
//      Each().Select("Select the field type", fieldTypes).
//      Ask("What should your username field be?", &username)
//		Ask("What should your password field be?", &password)
//      AskIndefinitely("Enter the field name in snake_case", userFields)
