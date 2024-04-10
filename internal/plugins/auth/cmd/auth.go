package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"pressebo/api/cmder"
	"pressebo/api/fsys"
	"slices"
	"text/template"

	_ "embed"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var wd, _ = os.Getwd()

//go:embed login_input_stub.txt
var loginStoreInput string

//go:embed registration_input_stub.txt
var registrationStoreInput string

//go:embed 20230222004736_create_orgs_table.txt
var orgsMigration string

//go:embed 20231128193645_create_users_table.txt
var usersMigration string

type Config struct {
	//
}

type Field struct {
	FieldName  string
	FieldType  string // text, textarea, number, boolean, radio, checkbox, dropdown, date, time, image
	IsUsername bool
	IsPassword bool
	IsUnique   bool
	IsRequired bool
	Choices    []string
}

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Generate auth related files",
	Long:  `Generate auth related files`,

	Run: func(cmd *cobra.Command, args []string) {
		username, password := "", ""
		fields := []Field{}
		hasOrg := false

		cmder.Confirm("Should your users belong to an org? (useful for multitenant apps)", 'n').Fill(&hasOrg).
			Ask("What should your username field be called? (in snake_case)", cmder.SnakeCase).Fill(&username).
			Ask("What should your password field be called? (in snake_case)", cmder.SnakeCase).Fill(&password).
			AskRecurring("Enter the field name (in snake_case)", cmder.SnakeCaseEmptyAllowed, func(result any) cmder.Prompter {
				var required, unique bool
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
						return prompt.AskRecurring("Enter choices", cmder.SnakeCaseEmptyAllowed).Fill(&choices)
					}).
					Confirm("Is this a required field?", 'n').Fill(&required).
					Confirm("Is this a unique field?", 'n').Fill(&unique)

				fields = append(fields, Field{FieldName: result.(string), FieldType: selectedType, Choices: choices, IsRequired: required, IsUnique: unique})

				return prompt
			})

		fields = append(fields, Field{FieldName: username, FieldType: "text", IsUsername: true, IsRequired: true, IsUnique: true})
		fields = append(fields, Field{FieldName: password, FieldType: "text", IsPassword: true, IsRequired: true})

		createInputFiles(fields)
	},
}

func createMigrationFieldsString(fields []Field) string {
	var fieldsString string
	for index, f := range fields {
		fieldsString += fmt.Sprintf("\t%s %s `json:\"%s\" db:\"%s\"`", strcase.ToCamel(f.FieldName), f.FieldType, f.FieldName, f.FieldName)
		if index < len(fields)-1 {
			fieldsString += "\n"
		}
	}
	return fieldsString
}

func createInputFiles(fields []Field) {
	createInputDir()
	createLoginInputFile(fields)
	createRegistrationInputFile(fields)
}

func createInputDir() {
	fs := fsys.NewLocalStorage("")
	err := fs.CreateDirectory("./internal/inputs")
	if err != nil {
		fmt.Println("Error creating inputs directory:", err.Error())
		return
	}
}

func createLoginInputFile(fields []Field) {
	fs := fsys.NewLocalStorage("")
	fmt.Println("Creating login_input.go file")
	loginInputFilePath := "./internal/inputs/login_input.go"
	username := slices.IndexFunc(fields, func(f Field) bool {
		return f.IsUsername
	})
	password := slices.IndexFunc(fields, func(f Field) bool {
		return f.IsPassword
	})

	tmplData := map[string]string{
		"UsernameTitle": strcase.ToCamel(fields[username].FieldName),
		"Username":      fields[username].FieldName,
		"PasswordTitle": strcase.ToCamel(fields[password].FieldName),
		"Password":      fields[password].FieldName,
	}

	loginTmpl, err := parseTemplate(tmplData, loginStoreInput)

	if err != nil {
		fmt.Println("Error parsing template")
		return
	}

	err = fs.Write(loginInputFilePath, []byte(loginTmpl))
	if err != nil {
		fmt.Println("Error creating file:", err.Error())
		return
	}

	fmt.Printf("Created login_input.go file at %s\n", loginInputFilePath)
}

func createRegistrationInputFile(fields []Field) {
	fs := fsys.NewLocalStorage("")
	fmt.Println("Creating registration_input.go file")
	registrationInputFilePath := "./internal/inputs/registration_input.go"
	// registrationInputFileName := "registration_input.go"

	tmplData := map[string]string{
		"Fields": createInputFieldsString(fields),
	}

	registrationTmpl, err := parseTemplate(tmplData, registrationStoreInput)

	if err != nil {
		fmt.Println("Error parsing template")
		return
	}

	err = fs.Write(registrationInputFilePath, []byte(registrationTmpl))
	if err != nil {
		fmt.Println("Error creating file:", err.Error())
		return
	}

	fmt.Printf("Created registration_input.go file at %s\n", registrationInputFilePath)
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
