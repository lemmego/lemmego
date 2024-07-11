package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"lemmego/api"
	"lemmego/api/cmd"
	"lemmego/api/cmder"
	"lemmego/api/fsys"
	"os"
	"slices"
	"text/template"

	_ "embed"

	"github.com/charmbracelet/huh"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var formFieldTypes = []string{"text", "textarea", "integer", "decimal", "boolean", "radio", "checkbox", "dropdown", "date", "time", "datetime", "image"}

var wd, _ = os.Getwd()

type Field struct {
	FieldName  string
	FieldType  string // text, textarea, number, boolean, radio, checkbox, dropdown, date, time, image
	IsUsername bool
	IsPassword bool
	IsUnique   bool
	IsRequired bool
	Choices    []string
}

func GetInstallCommand(p api.Plugin) *cobra.Command {
	var AuthCmd = &cobra.Command{
		Use:   "auth",
		Short: "Generate auth related files",
		Long:  `Generate auth related files`,

		Run: func(cmd *cobra.Command, args []string) {
			selectedFrontend := ""
			username, password := "email", "password"
			fields := []*Field{}
			hasOrg := false

			orgForm := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Which frontend scaffolding should be generated?").
						Options(huh.NewOptions("templ", "react")...).
						Value(&selectedFrontend),
					huh.NewConfirm().
						Title("Should your users belong to an org? (useful for multitenant apps)").
						Value(&hasOrg),
				),
			)

			err := orgForm.Run()
			if err != nil {
				fmt.Println("Error:", err.Error())
				return
			}

			for {
				var fieldName, fieldType string
				var required, unique bool
				choices := []string{}
				fieldNameForm := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Enter the field name in snake_case").
							Value(&fieldName).
							Validate(cmder.NotIn(
								[]string{"id", "email", "password", "org_username", "created_at", "updated_at", "deleted_at"},
								"No need to add this field, it will be provided.",
								cmder.SnakeCaseEmptyAllowed,
							)),
					),
				)

				err = fieldNameForm.Run()
				if err != nil {
					fmt.Println("Error:", err.Error())
					return
				}

				if fieldName == "" {
					break
				}

				fieldTypeForm := huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("Select the field type").
							Options(huh.NewOptions(formFieldTypes...)...).
							Value(&fieldType),
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

				fields = append(fields, &Field{FieldName: fieldName, FieldType: fieldType, Choices: choices, IsRequired: required, IsUnique: unique})
			}

			fields = append(fields, &Field{FieldName: username, FieldType: "text", IsUsername: true, IsRequired: true, IsUnique: true})
			fields = append(fields, &Field{FieldName: password, FieldType: "text", IsPassword: true, IsRequired: true})
			if hasOrg {
				fields = append(fields, &Field{FieldName: "org_username", FieldType: "text", IsRequired: true, IsUnique: true})
			}

			// cmder.Confirm("Should your users belong to an org? (useful for multitenant apps)", 'n').Fill(&hasOrg).
			// 	AskRepeat(
			// 		"Enter the field name in snake_case",
			// 		cmder.NotIn(
			// 			[]string{"id", "email", "password", "org_username", "created_at", "updated_at", "deleted_at"},
			// 			"No need to add this field, it will be provided.",
			// 			cmder.SnakeCaseEmptyAllowed,
			// 		),
			// 		func(result any) cmder.Prompter {
			// 			var required, unique bool
			// 			choices := []string{}
			// 			selectedType := ""
			// 			prompt := cmder.Select(
			// 				"Select the field type",
			// 				[]string{"text", "textarea", "integer", "decimal", "boolean", "radio", "checkbox", "dropdown", "date", "time", "image"},
			// 			).Fill(&selectedType).
			// 				When(func(res any) bool {
			// 					if val, ok := res.(string); ok {
			// 						return val == "radio" || val == "checkbox" || val == "dropdown"
			// 					}
			// 					return false
			// 				}, func(prompt cmder.Prompter) cmder.Prompter {
			// 					return prompt.AskRepeat("Enter choices", cmder.SnakeCaseEmptyAllowed).Fill(&choices)
			// 				}).
			// 				Confirm("Is this a required field?", 'n').Fill(&required).
			// 				Confirm("Is this a unique field?", 'n').Fill(&unique)

			// 			fields = append(fields, &Field{FieldName: result.(string), FieldType: selectedType, Choices: choices, IsRequired: required, IsUnique: unique})

			// 			return prompt
			// 		})

			// fields = append(fields, &Field{FieldName: username, FieldType: "text", IsUsername: true, IsRequired: true, IsUnique: true})
			// fields = append(fields, &Field{FieldName: password, FieldType: "text", IsPassword: true, IsRequired: true})
			// if hasOrg {
			// 	fields = append(fields, &Field{FieldName: "org_username", FieldType: "text", IsRequired: true, IsUnique: true})
			// }

			createInputFiles(fields)
			createMigrationFiles(fields)
			createModelFiles(fields)
			createFormFiles(fields, selectedFrontend)
		},
	}

	return AuthCmd
}

func generateOrgMigration() {
	om := cmd.NewMigrationGenerator(&cmd.MigrationConfig{
		TableName: "orgs",
		Fields: []*cmd.MigrationField{
			{Name: "id", Type: "bigIncrements", Primary: true},
			{Name: "org_username", Type: "string", Unique: true},
			{Name: "org_name", Type: "string"},
			{Name: "email", Type: "string", Unique: true},
		},
		Timestamps: true,
	})
	om.Generate()
}

func generateUserMigration(userFields []*cmd.MigrationField) {
	um := cmd.NewMigrationGenerator(&cmd.MigrationConfig{
		TableName:      "users",
		Fields:         userFields,
		Timestamps:     true,
		PrimaryColumns: []string{"id", "org_id"},
	})
	um.BumpVersion().Generate()
}

func generateOrgModel() {
	orgFields := []*cmd.ModelField{
		{Name: "id", Type: "bigIncrements"},
		{Name: "org_username", Type: "string", Unique: true},
		{Name: "org_name", Type: "string"},
		{Name: "email", Type: "string", Unique: true},
	}

	om := cmd.NewModelGenerator(&cmd.ModelConfig{
		Name:   "org",
		Fields: orgFields,
	})
	om.Generate()
}

func generateUserModel(userFields []*cmd.ModelField) {
	um := cmd.NewModelGenerator(&cmd.ModelConfig{
		Name:   "user",
		Fields: userFields,
	})
	um.Generate()
}

func createModelFiles(fields []*Field) {
	createModelDir()
	hasOrg := slices.ContainsFunc(fields, func(f *Field) bool { return f.FieldName == "org_username" })
	userFields := []*cmd.ModelField{}
	if hasOrg {
		generateOrgModel()
		userFields = append(userFields, &cmd.ModelField{
			Name: "org_id",
			Type: "bigIncrements",
		})
	}
	for _, f := range fields {
		userFields = append(userFields, &cmd.ModelField{
			Name:     f.FieldName,
			Type:     cmd.UiDbTypeMap[f.FieldType],
			Required: f.IsRequired,
			Unique:   f.IsUnique,
		})
	}
	generateUserModel(userFields)
}

func createMigrationFiles(fields []*Field) {
	hasOrg := slices.ContainsFunc(fields, func(f *Field) bool { return f.FieldName == "org_username" })
	userFields := []*cmd.MigrationField{
		{Name: "id", Type: "bigIncrements"},
	}
	uniqueColumns := []string{}
	for _, v := range fields {
		if v.IsUnique {
			uniqueColumns = append(uniqueColumns, v.FieldName)
		}
	}
	if hasOrg {
		generateOrgMigration()
		userFields = append(userFields, &cmd.MigrationField{
			Name:               "org_id",
			Type:               "bigIncrements",
			ForeignConstrained: true,
		})
	}

	for _, v := range fields {
		userFields = append(userFields, &cmd.MigrationField{
			Name:     v.FieldName,
			Type:     cmd.UiDbTypeMap[v.FieldType],
			Nullable: !v.IsRequired,
			Unique:   v.IsUnique,
		})
	}
	generateUserMigration(userFields)
}

func createInputFiles(fields []*Field) {
	createInputDir()
	inputFields := []*cmd.InputField{}
	registrationFields := []*cmd.InputField{}
	for _, f := range fields {
		if f.IsUsername || f.IsPassword {
			inputFields = append(inputFields, &cmd.InputField{
				Name: f.FieldName,
				Type: cmd.UiDataTypeMap[f.FieldType],
			})
		} else {
			registrationFields = append(registrationFields, &cmd.InputField{
				Name: f.FieldName,
				Type: cmd.UiDataTypeMap[f.FieldType],
			})
		}
	}

	loginGen := cmd.NewInputGenerator(&cmd.InputConfig{
		Name:   "login",
		Fields: inputFields,
	})
	loginGen.Generate()

	registrationGen := cmd.NewInputGenerator(&cmd.InputConfig{
		Name:   "registration",
		Fields: registrationFields,
	})
	registrationGen.Generate()
}

func createFormFiles(fields []*Field, flavor string) {
	createFormDir(flavor)
	formFields := []*cmd.FormField{}
	for _, f := range fields {
		formFields = append(formFields, &cmd.FormField{
			Name:    f.FieldName,
			Type:    cmd.UiDataTypeMap[f.FieldType],
			Choices: f.Choices,
		})
	}

	formGen := cmd.NewFormGenerator(&cmd.FormConfig{
		Name:   "auth",
		Flavor: flavor,
		Fields: formFields,
		Route:  "/auth",
	})
	formGen.Generate()
}

func createInputDir() {
	fs := fsys.NewLocalStorage("")
	err := fs.CreateDirectory("./internal/inputs")
	if err != nil {
		fmt.Println("Error creating inputs directory:", err.Error())
		return
	}
}

func createFormDir(flavor string) {
	if flavor == "react" {
		fs := fsys.NewLocalStorage("")
		err := fs.CreateDirectory("./resources/js/Pages/Forms")
		if err != nil {
			fmt.Println("Error creating forms directory:", err.Error())
			return
		}
	}

	if flavor == "templ" {
		fs := fsys.NewLocalStorage("")
		err := fs.CreateDirectory("./templates")
		if err != nil {
			fmt.Println("Error creating forms directory:", err.Error())
			return
		}
	}
}

func createModelDir() {
	fs := fsys.NewLocalStorage("")
	err := fs.CreateDirectory("./internal/models")
	if err != nil {
		fmt.Println("Error creating models directory:", err.Error())
		return
	}
}

func createInputFieldsString(fields []Field) string {
	var fieldsString string
	for index, f := range fields {
		fieldsString += fmt.Sprintf("\t%s string `json:\"%s\" in:\"form=%s\"`", strcase.ToCamel(f.FieldName), f.FieldName, f.FieldName)
		if index < len(fields)-1 {
			fieldsString += "\n"
		}
	}
	return fieldsString
}

func parseTemplate(tmplData map[string]string, fileContents string) (string, error) {
	var out bytes.Buffer
	tx := template.New("template")
	t := template.Must(tx.Parse(fileContents))
	err := t.Execute(&out, tmplData)
	if err != nil {
		return "", errors.New("Unable to execute template:" + err.Error())
	}
	return out.String(), nil
}
