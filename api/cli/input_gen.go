package cli

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"

	"github.com/lemmego/lemmego/api/fsys"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

//go:embed input.txt
var inputStub string

var inputFieldTypes = []string{
	"int", "uint", "int64", "uint64", "float64", "string", "bool", "time.Time", "custom",
}

type InputField struct {
	Name     string
	Type     string
	Required bool
	Unique   bool
	Table    string
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
		"InputName":   ig.name,
		"Fields":      ig.fields,
	}

	output, err := ParseTemplate(tmplData, ig.GetStub(), commonFuncs)

	if err != nil {
		return err
	}

	err = fs.Write(ig.GetPackagePath()+"/"+ig.name+"_input.go", []byte(output))

	if err != nil {
		return err
	}

	return nil
}

func (ig *InputGenerator) Command() *cobra.Command {
	return inputCmd
}

var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "Generate a request input",
	Long:  `Generate a request input`,
	Run: func(cmd *cobra.Command, args []string) {
		var inputName string
		var fields []*InputField

		nameForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Enter the input name in snake_case").
					Value(&inputName).
					Validate(SnakeCase),
			),
		)
		err := nameForm.Run()
		if err != nil {
			return
		}

		for {
			var fieldName, fieldType string
			const required = "Required"
			const unique = "Unique"
			selectedAttrs := []string{}
			fieldNameForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Enter the field name in snake_case").
						Validate(SnakeCaseEmptyAllowed).
						Value(&fieldName),
				),
			)
			err := fieldNameForm.Run()
			if err != nil {
				return
			}
			if fieldName == "" {
				break
			}
			fieldTypeForm := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("What should the data type be?").
						Options(huh.NewOptions(inputFieldTypes...)...).
						Value(&fieldType),
				),
			)
			err = fieldTypeForm.Run()
			if err != nil {
				return
			}

			if fieldType == "custom" {
				fieldTypeForm := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Enter the data type (You'll need to import it if necessary)").
							Value(&fieldType),
					),
				)
				err = fieldTypeForm.Run()
				if err != nil {
					return
				}
			}
			selectedAttrsForm := huh.NewForm(
				huh.NewGroup(
					huh.NewMultiSelect[string]().
						Title("Press x to select the attributes").
						Options(huh.NewOptions(required, unique)...).
						Value(&selectedAttrs),
				),
			)
			err = selectedAttrsForm.Run()
			if err != nil {
				return
			}

			fields = append(fields, &InputField{
				Name:     fieldName,
				Type:     fieldType,
				Required: slices.Contains(selectedAttrs, required),
				Unique:   slices.Contains(selectedAttrs, unique),
			})
		}

		mg := NewInputGenerator(&InputConfig{Name: inputName, Fields: fields})
		err = mg.Generate()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Input generated successfully.")
	},
}
