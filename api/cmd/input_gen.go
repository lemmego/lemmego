package cmd

import (
	_ "embed"
	"fmt"
	"pressebo/api/cmder"
	"pressebo/api/fsys"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

//go:embed input.txt
var inputStub string

type InputField struct {
	Name string
	Type string
}

type InputConfig struct {
	Name   string
	Fields []*InputField
}

type InputGenerator struct {
	name   string
	fields []*InputField
}

func NewInputGenerator(mc *InputConfig) *InputGenerator {
	return &InputGenerator{mc.Name, mc.Fields}
}

func (ig *InputGenerator) GetReplacables() []*Replacable {
	var fieldLines string
	for index, f := range ig.fields {
		fieldLines += fmt.Sprintf("\t%s %s `json:\"%s\" in:\"form=%s\"`", strcase.ToCamel(f.Name), f.Type, f.Name, f.Name)
		if index < len(ig.fields)-1 {
			fieldLines += "\n"
		}
	}
	return []*Replacable{
		{Placeholder: "InputName", Value: strcase.ToCamel(ig.name)},
		{Placeholder: "Fields", Value: fieldLines},
	}
}

func (ig *InputGenerator) GetPackagePath() string {
	return "internal/inputs"
}

func (ig *InputGenerator) GetStub() string {
	return inputStub
}

func (ig *InputGenerator) Generate() error {
	fs := fsys.NewLocalStorage("")
	parts := strings.Split(ig.GetPackagePath(), "/")
	packageName := ig.GetPackagePath()

	if len(parts) > 0 {
		packageName = parts[len(parts)-1]
	}

	tmplData := map[string]interface{}{
		"PackageName": packageName,
	}

	for _, v := range ig.GetReplacables() {
		tmplData[v.Placeholder] = v.Value
	}

	output, err := ParseTemplate(tmplData, ig.GetStub(), nil)

	if err != nil {
		return err
	}

	err = fs.Write(ig.GetPackagePath()+"/"+ig.name+"_input.go", []byte(output))

	if err != nil {
		return err
	}

	return nil
}

var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "Generate a request input",
	Long:  `Generate a request input`,
	Run: func(cmd *cobra.Command, args []string) {
		var inputName string
		var fields []*InputField

		cmder.Ask("Enter the input name in snake_case", cmder.SnakeCase).Fill(&inputName).
			AskRecurring("Enter the field name in snake_case", cmder.SnakeCaseEmptyAllowed, func(result any) cmder.Prompter {
				selectedType := ""
				prompt := cmder.Ask("What should the data type be? (https://go.dev/ref/spec#Types)", cmder.SnakeCase).Fill(&selectedType)

				fields = append(fields, &InputField{Name: result.(string), Type: selectedType})

				return prompt
			})

		mg := NewInputGenerator(&InputConfig{Name: inputName, Fields: fields})
		err := mg.Generate()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Input generated successfully.")
	},
}
