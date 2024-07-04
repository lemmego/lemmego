package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"pressebo/api"
	"pressebo/api/cmd"
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
			username, password := "email", "password"
			fields := []*Field{}
			hasOrg := false

			cmder.Confirm("Should your users belong to an org? (useful for multitenant apps)", 'n').Fill(&hasOrg).
				AskRepeat(
					"Enter the field name in snake_case",
					cmder.NotIn(
						[]string{"id", "email", "password", "org_username", "created_at", "updated_at", "deleted_at"},
						"No need to add this field, it will be provided.",
						cmder.SnakeCaseEmptyAllowed,
					),
					func(result any) cmder.Prompter {
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
								return prompt.AskRepeat("Enter choices", cmder.SnakeCaseEmptyAllowed).Fill(&choices)
							}).
							Confirm("Is this a required field?", 'n').Fill(&required).
							Confirm("Is this a unique field?", 'n').Fill(&unique)

						fields = append(fields, &Field{FieldName: result.(string), FieldType: selectedType, Choices: choices, IsRequired: required, IsUnique: unique})

						return prompt
					})

			fields = append(fields, &Field{FieldName: username, FieldType: "text", IsUsername: true, IsRequired: true, IsUnique: true})
			fields = append(fields, &Field{FieldName: password, FieldType: "text", IsPassword: true, IsRequired: true})
			if hasOrg {
				fields = append(fields, &Field{FieldName: "org_username", FieldType: "text", IsRequired: true, IsUnique: true})
			}

			createInputFiles(fields)
			createMigrationFiles(fields)
			createModelFiles(fields)
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

func createModelFiles(fields []*Field) {
	// hasOrg := slices.ContainsFunc(fields, func(f *Field) bool { return f.FieldName == "org_username" })
	userFields := []*cmd.ModelField{}
	// for _, f := range fields {
	// }
	m := cmd.NewModelGenerator(&cmd.ModelConfig{
		Name:   "org",
		Fields: userFields,
	})
	m.Generate()
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

func createInputDir() {
	fs := fsys.NewLocalStorage("")
	err := fs.CreateDirectory("./internal/inputs")
	if err != nil {
		fmt.Println("Error creating inputs directory:", err.Error())
		return
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
