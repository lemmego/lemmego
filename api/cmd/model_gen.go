package cmd

import (
	_ "embed"
	"fmt"
	"pressebo/api/cmder"
	"pressebo/api/fsys"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

//go:embed model.txt
var modelStub string

type ModelField struct {
	Name        string
	Type        string
	Unique      bool
	Required    bool
	SkipIfEmpty bool
}

type ModelConfig struct {
	Name   string
	Fields []*ModelField
}

type ModelGenerator struct {
	name   string
	fields []*ModelField
}

func NewModelGenerator(mc *ModelConfig) *ModelGenerator {
	return &ModelGenerator{mc.Name, mc.Fields}
}

func (mg *ModelGenerator) GetReplacables() []*Replacable {
	var fieldLines string
	for index, f := range mg.fields {
		omitEmpty := ""
		if f.SkipIfEmpty {
			omitEmpty = ",omitempty"
		}
		fieldLines += fmt.Sprintf("\t%s %s `json:\"%s\" db:\"%s%s\"`", strcase.ToCamel(f.Name), f.Type, f.Name, f.Name, omitEmpty)
		if index < len(mg.fields)-1 {
			fieldLines += "\n"
		}
	}
	return []*Replacable{
		{Placeholder: "ModelName", Value: strcase.ToCamel(mg.name)},
		{Placeholder: "Fields", Value: fieldLines},
	}
}

func (mg *ModelGenerator) GetPackagePath() string {
	return "internal/models"
}

func (mg *ModelGenerator) GetStub() string {
	return modelStub
}

func (mg *ModelGenerator) Generate() error {
	fs := fsys.NewLocalStorage("")
	parts := strings.Split(mg.GetPackagePath(), "/")
	packageName := mg.GetPackagePath()

	if len(parts) > 0 {
		packageName = parts[len(parts)-1]
	}

	tmplData := map[string]interface{}{
		"PackageName": packageName,
	}

	for _, v := range mg.GetReplacables() {
		tmplData[v.Placeholder] = v.Value
	}

	output, err := ParseTemplate(tmplData, mg.GetStub(), nil)

	if err != nil {
		return err
	}

	err = fs.Write(mg.GetPackagePath()+"/"+mg.name+".go", []byte(output))

	if err != nil {
		return err
	}

	return nil
}

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generate a db model",
	Long:  `Generate a db model`,
	Run: func(cmd *cobra.Command, args []string) {
		var modelName string
		var fields []*ModelField

		cmder.Ask("Enter the model name in snake_case", cmder.SnakeCase).Fill(&modelName).
			AskRecurring("Enter the field name in snake_case", cmder.SnakeCaseEmptyAllowed, func(result any) cmder.Prompter {
				const required = "Required"
				const unique = "Unique"
				const omitEmpty = "Omit Empty (Skip from db when the value is empty)"
				selectedAttrs := []string{}
				selectedType := ""
				prompt := cmder.Ask("What should the data type be? (https://go.dev/ref/spec#Types)", cmder.SnakeCase).
					Fill(&selectedType).
					MultiSelect("Select the attributes for this field:", []*cmder.Item{
						{Label: required}, {Label: unique}, {Label: omitEmpty},
					}, 0).Fill(&selectedAttrs)

				fields = append(
					fields,
					&ModelField{
						Name:        result.(string),
						Type:        selectedType,
						Required:    slices.Contains(selectedAttrs, required),
						Unique:      slices.Contains(selectedAttrs, unique),
						SkipIfEmpty: slices.Contains(selectedAttrs, omitEmpty),
					},
				)

				return prompt
			})

		mg := NewModelGenerator(&ModelConfig{Name: modelName, Fields: fields})
		err := mg.Generate()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Model generated successfully.")
	},
}
