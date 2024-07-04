package cmd

import (
	_ "embed"
	"fmt"
	"pressebo/api/cmder"
	"pressebo/api/fsys"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

//go:embed templ_form.txt
var templFormStub string

type TemplField struct {
	Name    string
	Type    string
	Choices []string
}

type TemplConfig struct {
	Name   string
	Fields []*TemplField
	Route  string
}

type TemplGenerator struct {
	name   string
	fields []*TemplField
	route  string
}

func NewTemplGenerator(mc *TemplConfig) *TemplGenerator {
	return &TemplGenerator{mc.Name, mc.Fields, mc.Route}
}

func (ig *TemplGenerator) GetReplacables() []*Replacable {
	return []*Replacable{
		{Placeholder: "TemplName", Value: strcase.ToCamel(ig.name)},
		{Placeholder: "Route", Value: ig.route},
		{Placeholder: "Fields", Value: ig.fields},
	}
}

func (ig *TemplGenerator) GetPackagePath() string {
	return "templates"
}

func (ig *TemplGenerator) GetStub() string {
	return templFormStub
}

func (ig *TemplGenerator) Generate() error {
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

	output, err := ParseTemplate(tmplData, ig.GetStub(), template.FuncMap{
		"toCamel": func(str string) string {
			return strcase.ToCamel(str)
		},
		"toSnake": func(str string) string {
			return strcase.ToSnake(str)
		},
		"toSpaceDelimited": func(str string) string {
			return strcase.ToDelimited(str, ' ')
		},
	})

	if err != nil {
		return err
	}

	err = fs.Write(ig.GetPackagePath()+"/"+ig.name+".templ", []byte(output))

	if err != nil {
		return err
	}

	return nil
}

var templCmd = &cobra.Command{
	Use:   "templ",
	Short: "Generate a templ file",
	Long:  `Generate a templ file`,
	Run: func(cmd *cobra.Command, args []string) {
		var templName, route string
		var fields []*TemplField

		cmder.Ask("Enter the template name in snake_case", cmder.SnakeCase).Fill(&templName).
			Ask("Enter the route where the form should be submitted (e.g. /login)", nil).Fill(&route).
			AskRepeat("Enter the field name in snake_case", cmder.SnakeCaseEmptyAllowed, func(result any) cmder.Prompter {
				// var required, unique bool
				choices := []string{}
				selectedType := ""
				prompt := cmder.Select(
					"Select the field type",
					[]string{"text", "textarea", "integer", "decimal", "boolean", "radio", "checkbox", "dropdown", "date", "time", "image"},
				).Fill(&selectedType).
					When(func(res any) bool {
						if val, ok := res.(string); ok {
							return val == "radio" || val == "checkbox" || val == "dropdown"
						}
						return false
					}, func(prompt cmder.Prompter) cmder.Prompter {
						return prompt.AskRepeat("Enter choices", cmder.SnakeCaseEmptyAllowed).Fill(&choices)
					})
					// Confirm("Is this a required field?", 'n').Fill(&required).
					// Confirm("Is this a unique field?", 'n').Fill(&unique)

				fields = append(fields, &TemplField{Name: result.(string), Type: selectedType, Choices: choices})

				return prompt
			})

		mg := NewTemplGenerator(&TemplConfig{Name: templName, Fields: fields, Route: route})
		err := mg.Generate()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Template generated successfully.")
	},
}
