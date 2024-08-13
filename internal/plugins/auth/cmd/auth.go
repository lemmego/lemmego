package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"slices"
	"text/template"

	"github.com/lemmego/lemmego/api"
	"github.com/lemmego/lemmego/api/cli"
	"github.com/lemmego/lemmego/api/fsys"

	_ "embed"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var formFieldTypes = []string{"text", "textarea", "integer", "decimal", "boolean", "radio", "checkbox", "dropdown", "date", "time", "datetime", "file"}

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
			const required = "Required"
			const unique = "Unique"
			selectedAttrs := []string{}
			choices := []string{}
			fieldNameForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Enter the field name in snake_case.\nThe following fields will be provided:\nid, email, password, org_name, org_email, org_username, created_at, updated_at, deleted_at").
						Value(&fieldName).
						Validate(cli.NotIn(
							[]string{"id", "email", "password", "org_name", "org_email", "org_username", "created_at", "updated_at", "deleted_at"},
							"No need to add this field, it will be provided.",
							cli.SnakeCaseEmptyAllowed,
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

			fields = append(fields, &Field{
				FieldName:  fieldName,
				FieldType:  fieldType,
				Choices:    choices,
				IsRequired: slices.Contains(selectedAttrs, required),
				IsUnique:   slices.Contains(selectedAttrs, unique),
			})
		}

		fields = append(fields, &Field{FieldName: username, FieldType: "text", IsUsername: true, IsRequired: true, IsUnique: true})
		fields = append(fields, &Field{FieldName: password, FieldType: "text", IsPassword: true, IsRequired: true})

		//if hasOrg {
		//	fields = append(fields, &Field{FieldName: "org_username", FieldType: "text", IsRequired: true, IsUnique: true})
		//	fields = append(fields, &Field{FieldName: "org_name", FieldType: "text", IsRequired: true, IsUnique: false})
		//}

		createInputFiles(fields, hasOrg)
		createMigrationFiles(fields, hasOrg)
		createModelFiles(fields, hasOrg)
		createFormFiles(fields, selectedFrontend, hasOrg)
	},
}

func GetInstallCommand(p api.Plugin) *cobra.Command {
	return AuthCmd
}

func generateOrgMigration() {
	orgFields := []*cli.MigrationField{
		{Name: "id", Type: "bigIncrements", Primary: true},
		{Name: "org_username", Type: "string", Unique: true},
		{Name: "org_name", Type: "string"},
		{Name: "org_email", Type: "string", Unique: true},
	}
	om := cli.NewMigrationGenerator(&cli.MigrationConfig{
		TableName:  "orgs",
		Fields:     orgFields,
		Timestamps: true,
	})
	om.Generate()
}

func generateUserMigration(userFields []*cli.MigrationField) {
	um := cli.NewMigrationGenerator(&cli.MigrationConfig{
		TableName:      "users",
		Fields:         userFields,
		Timestamps:     true,
		PrimaryColumns: []string{"id", "org_id"},
	})
	um.BumpVersion().Generate()
}

func generateOrgModel() {
	orgFields := []*cli.ModelField{
		{Name: "org_username", Type: "string", Unique: true},
		{Name: "org_name", Type: "string"},
		{Name: "org_email", Type: "string", Unique: true},
	}

	om := cli.NewModelGenerator(&cli.ModelConfig{
		Name:   "org",
		Fields: orgFields,
	})
	om.Generate()
}

func generateUserModel(userFields []*cli.ModelField) {
	um := cli.NewModelGenerator(&cli.ModelConfig{
		Name:   "user",
		Fields: userFields,
	})
	um.Generate()
}

func createModelFiles(fields []*Field, hasOrg bool) {
	createModelDir()
	userFields := []*cli.ModelField{}
	if hasOrg {
		generateOrgModel()
		userFields = append(userFields, []*cli.ModelField{
			{
				Name: "org_id",
				Type: cli.UiDataTypeMap["integer"],
			},
			{
				Name: "org",
				Type: "Org",
			},
		}...)
	}
	for _, f := range fields {
		userFields = append(userFields, &cli.ModelField{
			Name:     f.FieldName,
			Type:     cli.UiDataTypeMap[f.FieldType],
			Required: f.IsRequired,
			Unique:   f.IsUnique,
		})
	}
	generateUserModel(userFields)
}

func createMigrationFiles(fields []*Field, hasOrg bool) {
	userFields := []*cli.MigrationField{
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
		userFields = append(userFields, &cli.MigrationField{
			Name:               "org_id",
			Type:               "bigIncrements",
			ForeignConstrained: true,
		})
	}

	for _, v := range fields {
		userFields = append(userFields, &cli.MigrationField{
			Name:     v.FieldName,
			Type:     cli.UiDbTypeMap[v.FieldType],
			Nullable: !v.IsRequired,
			Unique:   v.IsUnique,
		})
	}
	generateUserMigration(userFields)
}

func createInputFiles(fields []*Field, hasOrg bool) {
	createInputDir()
	loginFields := []*cli.InputField{}
	registrationFields := []*cli.InputField{}
	for _, f := range fields {
		if f.IsUsername || f.IsPassword {
			loginFields = append(loginFields, &cli.InputField{
				Name:     f.FieldName,
				Type:     cli.UiDataTypeMap[f.FieldType],
				Required: f.IsRequired,
				Unique:   f.IsUnique,
			})
		} else {
			registrationFields = append(registrationFields, &cli.InputField{
				Name:     f.FieldName,
				Type:     cli.UiDataTypeMap[f.FieldType],
				Required: f.IsRequired,
				Unique:   f.IsUnique,
			})
		}
	}

	if hasOrg {
		loginFields = append(loginFields, &cli.InputField{Name: "org_username", Type: "string", Required: true})
		registrationFields = append(registrationFields, []*cli.InputField{
			{Name: "org_name", Type: "string", Required: true},
			{Name: "org_email", Type: "string", Required: true},
			{Name: "org_username", Type: "string", Required: true},
		}...)
	}

	loginGen := cli.NewInputGenerator(&cli.InputConfig{
		Name:   "login",
		Fields: loginFields,
	})
	loginGen.Generate()

	registrationFields = append(registrationFields, []*cli.InputField{
		{Name: "email", Type: "string", Required: true},
		{Name: "password", Type: "string", Required: true},
		{Name: "password_confirmation", Type: "string", Required: true},
	}...)

	registrationGen := cli.NewInputGenerator(&cli.InputConfig{
		Name:   "registration",
		Fields: registrationFields,
	})
	registrationGen.Generate()
}

func createFormFiles(fields []*Field, flavor string, hasOrg bool) {
	createFormDir(flavor)
	registrationFields := []*cli.FormField{}
	for _, f := range fields {
		registrationFields = append(registrationFields, &cli.FormField{
			Name:    f.FieldName,
			Type:    f.FieldType,
			Choices: f.Choices,
		})
		if f.IsPassword {
			registrationFields = append(registrationFields, &cli.FormField{Name: "password_confirmation", Type: "text"})
		}
	}

	loginFields := []*cli.FormField{
		{Name: "email", Type: "text"},
		{Name: "password", Type: "text"},
	}

	if hasOrg {
		loginFields = append([]*cli.FormField{{Name: "org_username", Type: "text"}}, loginFields...)
		registrationFields = append(registrationFields, []*cli.FormField{
			{Name: "org_name", Type: "text"},
			{Name: "org_email", Type: "text"},
			{Name: "org_username", Type: "text"},
		}...)
	}

	loginForm := cli.NewFormGenerator(&cli.FormConfig{
		Name:   "login",
		Flavor: flavor,
		Fields: loginFields,
		Route:  "/login",
	})

	loginForm.Generate()

	regForm := cli.NewFormGenerator(&cli.FormConfig{
		Name:   "register",
		Flavor: flavor,
		Fields: registrationFields,
		Route:  "/register",
	})

	regForm.Generate()
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
