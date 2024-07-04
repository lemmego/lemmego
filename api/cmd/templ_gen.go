package cmd

import (
	_ "embed"
	"fmt"
	"pressebo/api/cmder"
	"pressebo/api/fsys"
	"strings"
	"text/template"

	"github.com/charmbracelet/huh"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

//go:embed templ_form.txt
var templFormStub string

var templFieldTypes = []string{"text", "textarea", "integer", "decimal", "boolean", "radio", "checkbox", "dropdown", "date", "time", "image"}

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

		nameForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Enter the resource name in snake_case").
					Value(&templName).
					Validate(cmder.SnakeCase),
				huh.NewInput().
					Title("Enter the route where the form should be submitted (e.g. /login)").
					Value(&route),
			),
		)

		err := nameForm.Run()
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			var fieldName, fieldType string
			var choices []string

			fieldNameForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Enter the field name in snake_case").
						Validate(cmder.SnakeCaseEmptyAllowed).
						Value(&fieldName),
				),
			)

			err = fieldNameForm.Run()
			if err != nil {
				fmt.Println(err)
				return
			}

			if fieldName == "" {
				break
			}

			fieldTypeForm := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Select the field type").
						Value(&fieldType).
						Options(huh.NewOptions(templFieldTypes...)...),
				),
			)

			err = fieldTypeForm.Run()
			if err != nil {
				fmt.Println(err)
				return
			}

			if fieldType == "radio" || fieldType == "checkbox" || fieldType == "dropdown" {

				for {
					var choice string
					choicesForm := huh.NewForm(
						huh.NewGroup(
							huh.NewInput().
								Title(fmt.Sprintf("Add new choice for %s %s (Press enter to finish)", fieldName, fieldType)).
								Value(&choice),
						),
					)

					err = choicesForm.Run()
					if err != nil {
						fmt.Println(err)
						return
					}

					if choice == "" {
						break
					}
					choices = append(choices, choice)
				}
			}
			fields = append(fields, &TemplField{Name: fieldName, Type: fieldType, Choices: choices})
		}

		mg := NewTemplGenerator(&TemplConfig{Name: templName, Fields: fields, Route: route})
		err = mg.Generate()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Template generated successfully.")
	},
}
