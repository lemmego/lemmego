package cmd

import (
	_ "embed"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"lemmego/api/cmder"
	"lemmego/api/fsys"
	"strings"
)

//go:embed templ_form.txt
var templFormStub string

//go:embed react_form.txt
var reactFormStub string

var formFieldTypes = []string{"text", "textarea", "integer", "decimal", "boolean", "radio", "checkbox", "dropdown", "date", "time", "datetime", "image"}

type FormField struct {
	Name    string
	Type    string
	Choices []string
}

type FormConfig struct {
	Name   string
	Flavor string // templ, react
	Fields []*FormField
	Route  string
}

type FormGenerator struct {
	name   string
	flavor string // templ, react
	fields []*FormField
	route  string
}

func NewFormGenerator(mc *FormConfig) *FormGenerator {
	return &FormGenerator{mc.Name, mc.Flavor, mc.Fields, mc.Route}
}

func (fg *FormGenerator) GetReplacables() []*Replacable {
	return []*Replacable{
		{Placeholder: "Name", Value: strcase.ToCamel(fg.name)},
		{Placeholder: "Route", Value: fg.route},
		{Placeholder: "Fields", Value: fg.fields},
	}
}

func (fg *FormGenerator) GetPackagePath() string {
	if fg.flavor == "react" {
		return "resources/js/Pages/Forms"
	}
	if fg.flavor == "templ" {
		return "templates"
	}
	return ""
}

func (fg *FormGenerator) GetStub() string {
	if fg.flavor == "react" {
		return reactFormStub
	}
	if fg.flavor == "templ" {
		return templFormStub
	}
	return ""
}

func (fg *FormGenerator) Generate() error {
	fs := fsys.NewLocalStorage("")
	parts := strings.Split(fg.GetPackagePath(), "/")
	packageName := fg.GetPackagePath()

	if len(parts) > 0 {
		packageName = parts[len(parts)-1]
	}

	tmplData := map[string]interface{}{
		"PackageName": packageName,
	}

	for _, v := range fg.GetReplacables() {
		tmplData[v.Placeholder] = v.Value
	}

	output, err := ParseTemplate(tmplData, fg.GetStub(), commonFuncs)

	if err != nil {
		return err
	}

	if fg.flavor == "templ" {
		err = fs.Write(fg.GetPackagePath()+"/"+fg.name+".templ", []byte(output))
	} else if fg.flavor == "react" {
		// Check if resources/js/Pages/Forms directory exists,
		// create it if it doesn't
		if exists, _ := fs.Exists(fg.GetPackagePath()); !exists {
			err = fs.CreateDirectory(fg.GetPackagePath())
			if err != nil {
				return err
			}
		}
		err = fs.Write(fg.GetPackagePath()+"/"+strcase.ToCamel(fg.name)+".tsx", []byte(output))
	}

	if err != nil {
		return err
	}

	return nil
}

var formCmd = &cobra.Command{
	Use:   "form",
	Short: "Generate a form template/view",
	Long:  `Generate a form template/view`,
	Run: func(cmd *cobra.Command, args []string) {
		var templName, route string
		var fields []*FormField
		var flavor string

		nameForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Which flavor do you want?").
					Options(huh.NewOptions("templ", "react")...).
					Value(&flavor),
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
						Options(huh.NewOptions(formFieldTypes...)...),
				),
			)

			err = fieldTypeForm.Run()
			if err != nil {
				fmt.Println("Error:", err.Error())
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
			fields = append(fields, &FormField{Name: fieldName, Type: fieldType, Choices: choices})
		}

		fg := NewFormGenerator(&FormConfig{Name: templName, Flavor: flavor, Fields: fields, Route: route})
		err = fg.Generate()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Template generated successfully.")
	},
}
