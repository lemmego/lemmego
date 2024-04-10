package cmder

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/manifoldco/promptui"
)

type PromptResultType int

type Item struct {
	Label      string
	IsSelected bool
}

const (
	PromptResultTypeNormal PromptResultType = iota
	PromptResultTypeBoolean
	PromptResultTypeRecurring
	PromptResultTypeSelect
	PromptResultTypeMultiSelect
)

type PromptResult struct {
	Type          PromptResultType
	ShouldAskNext bool
	Result        interface{}
	Error         error
}

type Prompter interface {
	Ask(question string, validator promptui.ValidateFunc) Prompter
	Confirm(question string, defaultValue rune) Prompter
	AskRecurring(question string, validator promptui.ValidateFunc, prompts ...func(result any) Prompter) Prompter
	Select(label string, items []string) Prompter
	MultiSelect(label string, items []*Item, selectedPos int) Prompter
	When(cb func(result interface{}) bool, thenPrompt func(prompt Prompter) Prompter) Prompter
	Fill(ptr any) Prompter
}

func (pr *PromptResult) Fill(ptr any) Prompter {
	if pr.ShouldAskNext {
		if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
			panic("Fill() must be called with a pointer")
		}
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

func (pr *PromptResult) Confirm(question string, defaultValue rune) Prompter {
	if pr.ShouldAskNext {
		return Confirm(question, defaultValue)
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

func (pr *PromptResult) MultiSelect(label string, allItems []*Item, selectedPos int) Prompter {
	if pr.ShouldAskNext {
		return MultiSelect(label, allItems, selectedPos)
	}
	return pr
}

func (pr *PromptResult) When(cb func(result interface{}) bool, thenPrompt func(prompt Prompter) Prompter) Prompter {
	if pr.ShouldAskNext {
		if cb(pr.Result) {
			return thenPrompt(pr)
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
		if err == promptui.ErrInterrupt {
			os.Exit(-1)
		}
		return &PromptResult{Type: PromptResultTypeNormal, ShouldAskNext: false, Result: nil, Error: err}
	}
	return &PromptResult{Type: PromptResultTypeNormal, ShouldAskNext: true, Result: res, Error: nil}
}

func Confirm(question string, defaultVal rune) Prompter {
	if defaultVal != 'y' && defaultVal != 'Y' && defaultVal != 'n' && defaultVal != 'N' {
		panic("defaultVal argument must be either of y, Y, n, N")
	}

	labelSuffix := " (%s/%s)"

	if defaultVal == 'y' || defaultVal == 'Y' {
		labelSuffix = fmt.Sprintf(labelSuffix, "Y", "n")
	}

	if defaultVal == 'n' || defaultVal == 'N' {
		labelSuffix = fmt.Sprintf(labelSuffix, "y", "N")
	}

	q := promptui.Prompt{
		Label: question + labelSuffix,
		Validate: func(s string) error {
			if s != "" && s != "y" && s != "Y" && s != "n" && s != "N" {
				return errors.New("Input must be either of y, Y, n, N")
			}
			return nil
		},
	}

	res, err := q.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			os.Exit(-1)
		}
		return &PromptResult{Type: PromptResultTypeBoolean, ShouldAskNext: false, Result: false, Error: err}
	}

	if res == "" {
		res = string(defaultVal)
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
		if err == promptui.ErrInterrupt {
			os.Exit(-1)
		}
		return &PromptResult{Type: PromptResultTypeSelect, ShouldAskNext: false, Result: nil, Error: err}
	}

	return &PromptResult{Type: PromptResultTypeSelect, ShouldAskNext: true, Result: result, Error: nil}
}

// MultiSelect() prompts user to select one or more items in the given slice
func MultiSelect(label string, allItems []*Item, selectedPos int) Prompter {
	// Always prepend a "Done" item to the slice if it doesn't
	// already exist.
	const doneID = "Done ✅"
	if len(allItems) > 0 {
		lastIndex := len(allItems) - 1
		if allItems[lastIndex].Label != doneID {
			var items = []*Item{
				{
					Label: doneID,
				},
			}
			allItems = append(allItems, items...)
		}
	}

	// Define promptui template
	templates := &promptui.SelectTemplates{
		Label: `{{if .IsSelected}}
                    ✔
                {{end}} {{ .Label }} - label`,
		Active:   "→ {{if .IsSelected}}✔ {{end}}{{ .Label | cyan }}",
		Inactive: "{{if .IsSelected}}✔ {{end}}{{ .Label | cyan }}",
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     allItems,
		Templates: templates,
		Size:      5,
		// Start the cursor at the currently selected index
		CursorPos:    selectedPos,
		HideSelected: true,
	}

	selectionIdx, _, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			os.Exit(-1)
		}
		return &PromptResult{Type: PromptResultTypeMultiSelect, ShouldAskNext: false, Result: nil, Error: err}
	}

	chosenItem := allItems[selectionIdx]

	if chosenItem.Label != doneID {
		// If the user selected something other than "Done",
		// toggle selection on this item and run the function again.
		chosenItem.IsSelected = !chosenItem.IsSelected
		return MultiSelect(label, allItems, selectionIdx)
	}

	var selectedLabels []string
	for _, i := range allItems {
		if i.IsSelected {
			selectedLabels = append(selectedLabels, i.Label)
		}
	}
	return &PromptResult{Type: PromptResultTypeMultiSelect, ShouldAskNext: true, Result: selectedLabels, Error: nil}
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
			if err == promptui.ErrInterrupt {
				os.Exit(-1)
			}
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
