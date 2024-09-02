package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/lemmego/api/app"
	"github.com/lemmego/api/cli"
	"github.com/lemmego/api/fsys"

	_ "embed"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var formFieldTypes = []string{"text", "textarea", "integer", "decimal", "boolean", "radio", "checkbox", "dropdown", "date", "time", "datetime", "file"}

var userFields = []string{"first_name", "last_name", "username", "bio", "phone", "avatar"}
var requiredUserFields = []string{"email", "password"}
var orgFields = []string{"org_name", "org_email", "org_logo"}
var requiredOrgFields = []string{"org_username"}
var wd, _ = os.Getwd()

type Field struct {
	Name     string
	Type     string
	Required bool
	Unique   bool
}

var uf = []*Field{
	{Name: "email", Type: "string", Required: true, Unique: true},
	{Name: "password", Type: "string", Required: true, Unique: false},
	{Name: "username", Type: "string", Required: true, Unique: false},
	{Name: "first_name", Type: "string", Required: true, Unique: false},
	{Name: "last_name", Type: "string", Required: true, Unique: false},
	{Name: "bio", Type: "string", Required: true, Unique: false},
	{Name: "phone", Type: "string", Required: true, Unique: false},
	{Name: "avatar", Type: "string", Required: true, Unique: false},
}

var of = []*Field{
	{Name: "avatar", Type: "string", Required: true, Unique: false},
}

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Generate auth related files",
	Long:  `Generate auth related files`,

	Run: func(cmd *cobra.Command, args []string) {
		selectedFrontend := ""
		// username, password := "email", "password"
		hasOrg := false

		selectedUserFields := []string{}
		selectedOrgFields := []string{}

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

		userFieldSelectionForm := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Select the fields for the user entity").
					Options(huh.NewOptions(userFields...)...).
					Value(&selectedUserFields),
			),
		)

		err = userFieldSelectionForm.Run()
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		if hasOrg {
			orgFieldSelectionForm := huh.NewForm(
				huh.NewGroup(
					huh.NewMultiSelect[string]().
						Title("Select the fields for the org entity").
						Options(huh.NewOptions(orgFields...)...).
						Value(&selectedOrgFields),
				),
			)

			err = orgFieldSelectionForm.Run()
			if err != nil {
				fmt.Println("Error:", err.Error())
				return
			}
		}

		createMigrationFiles(selectedUserFields, selectedOrgFields)
		createModelFiles(selectedUserFields, selectedOrgFields)
		createInputFiles(selectedUserFields, selectedOrgFields)
		createFormFiles(selectedFrontend, selectedUserFields, selectedOrgFields)
	},
}

func GetInstallCommand(p app.Plugin) *cobra.Command {
	return AuthCmd
}

func generateOrgMigration(oFields []*cli.MigrationField) {
	om := cli.NewMigrationGenerator(&cli.MigrationConfig{
		TableName:  "orgs",
		Fields:     oFields,
		Timestamps: true,
	})
	om.Generate()
}

func generateUserMigration(userFields []*cli.MigrationField, hasOrg bool) {
	config := &cli.MigrationConfig{
		TableName:  "users",
		Fields:     userFields,
		Timestamps: true,
	}
	if hasOrg {
		config.PrimaryColumns = []string{"id", "org_id"}
	} else {
		config.PrimaryColumns = []string{"id"}
	}
	um := cli.NewMigrationGenerator(config)
	um.BumpVersion().Generate()
}

func createMigrationFiles(userFields []string, orgFields []string) {
	hasOrg := len(orgFields) > 0
	uFields := []*cli.MigrationField{
		{Name: "id", Type: "bigIncrements"},
		{Name: "email", Type: "string", Unique: true},
		{Name: "password", Type: "text"},
	}

	for _, f := range userFields {
		field := &cli.MigrationField{Name: f, Type: "string"}
		if f == "username" || f == "email" {
			field.Unique = true
		}
		if f == "bio" {
			field.Nullable = true
			field.Type = "text"
		}
		uFields = append(uFields, field)
	}

	if hasOrg {
		oFields := []*cli.MigrationField{
			{Name: "id", Type: "bigIncrements", Primary: true},
			{Name: "org_username", Type: "string", Unique: true},
		}

		for _, f := range orgFields {
			field := &cli.MigrationField{Name: f, Type: "string"}
			if f == "username" || f == "email" {
				field.Unique = true
			}
			if f == "bio" {
				field.Nullable = true
				field.Type = "text"
			}
			oFields = append(oFields, field)
		}
		generateOrgMigration(oFields)

		uFields = append(uFields, &cli.MigrationField{
			Name:               "org_id",
			Type:               "bigIncrements",
			ForeignConstrained: true,
		})
	}

	generateUserMigration(uFields, hasOrg)
}

func generateOrgModel(orgFields []*cli.ModelField) {
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

func createModelFiles(userFields []string, orgFields []string) {
	createModelDir()
	uFields := []*cli.ModelField{
		{Name: "email", Type: "string", Unique: true},
		{Name: "password", Type: "string"},
	}

	if len(orgFields) > 0 {
		oFields := []*cli.ModelField{
			{Name: "org_username", Type: "string", Unique: true},
		}
		for _, f := range orgFields {
			field := &cli.ModelField{Name: f, Type: "string"}
			if f == "org_email" {
				field.Unique = true
			}
			field.Required = true
			oFields = append(oFields, field)
		}

		uFields = append(uFields, &cli.ModelField{
			Name: "org_id", Type: "uint", Required: true,
		})
		generateOrgModel(oFields)
	}

	for _, f := range userFields {
		field := &cli.ModelField{Name: f, Type: "string"}
		if f == "username" || f == "email" {
			field.Required = true
			field.Unique = true
		}
		uFields = append(uFields, field)
	}
	generateUserModel(uFields)

}

func createInputFiles(userFields []string, orgFields []string) {
	createInputDir()
	loginFields := []*cli.InputField{
		{Name: "email", Type: "string", Required: true},
		{Name: "password", Type: "string", Required: true},
	}

	registrationFields := []*cli.InputField{
		{Name: "email", Type: "string", Required: true, Unique: true, Table: "users"},
		{Name: "password", Type: "string", Required: true},
		{Name: "password_confirmation", Type: "string", Required: true},
	}

	for _, f := range userFields {
		registrationFields = append(registrationFields, &cli.InputField{
			Name:     f,
			Type:     "string",
			Required: true,
		})
	}

	if len(orgFields) > 0 {
		for _, f := range orgFields {
			field := &cli.InputField{
				Name:     f,
				Type:     "string",
				Required: true,
			}

			if f == "org_email" {
				field.Unique = true
				field.Table = "orgs"
			}

			registrationFields = append(registrationFields, field)
		}
		orgUsernameField := &cli.InputField{
			Name:     "org_username",
			Type:     "string",
			Required: true,
			Unique:   true,
			Table:    "orgs",
		}
		loginFields = append(loginFields, orgUsernameField)
		registrationFields = append(registrationFields, orgUsernameField)
	}

	loginGen := cli.NewInputGenerator(&cli.InputConfig{
		Name:   "login",
		Fields: loginFields,
	})
	loginGen.Generate()

	registrationGen := cli.NewInputGenerator(&cli.InputConfig{
		Name:   "registration",
		Fields: registrationFields,
	})
	registrationGen.Generate()
}

func createFormFiles(flavor string, userFields []string, orgFields []string) {
	createFormDir(flavor)
	loginFields := []*cli.FormField{
		{Name: "email", Type: "text"},
		{Name: "password", Type: "text"},
	}
	registrationFields := []*cli.FormField{}

	for _, f := range userFields {
		field := &cli.FormField{Name: f, Type: "text"}
		if f == "avatar" {
			field.Type = "file"
		}
		if f == "bio" {
			field.Type = "textarea"
		}
		registrationFields = append(registrationFields, field)
	}

	registrationFields = append(registrationFields, []*cli.FormField{
		{Name: "email", Type: "text"},
		{Name: "password", Type: "text"},
		{Name: "password_confirmation", Type: "text"},
	}...)

	if len(orgFields) > 0 {
		loginFields = append([]*cli.FormField{{Name: "org_username", Type: "text"}}, loginFields...)
		registrationFields = append(registrationFields, &cli.FormField{Name: "org_username", Type: "text"})
		for _, f := range orgFields {
			field := &cli.FormField{Name: f, Type: "text"}
			if f == "org_logo" {
				field.Type = "file"
			}
			registrationFields = append(registrationFields, field)
		}
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
