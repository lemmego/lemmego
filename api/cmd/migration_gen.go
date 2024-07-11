package cmd

import (
	_ "embed"
	"fmt"
	"lemmego/api/cmder"
	"lemmego/api/fsys"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

//go:embed migration.txt
var migrationStub string

var migrationFieldTypes = []string{
	"increments", "bigIncrements", "int", "bigInt", "string", "text", "boolean", "unsignedInt", "unsignedBigInt", "decimal", "dateTime", "time",
}

type MigrationField struct {
	Name               string
	Type               string
	Nullable           bool
	Unique             bool
	Primary            bool
	ForeignConstrained bool
}

type MigrationConfig struct {
	TableName      string
	Fields         []*MigrationField
	PrimaryColumns []string
	UniqueColumns  [][]string
	ForeignColumns [][]string
	Timestamps     bool
}

type MigrationGenerator struct {
	name      string
	tableName string
	fields    []*MigrationField
	version   string

	primaryColumns []string
	uniqueColumns  [][]string
	foreignColumns [][]string
}

func NewMigrationGenerator(mc *MigrationConfig) *MigrationGenerator {
	version := time.Now().Format("20060102150405")
	if mc.Timestamps {
		timeStampFields := []*MigrationField{
			{Name: "created_at", Type: "dateTime", Nullable: true},
			{Name: "updated_at", Type: "dateTime", Nullable: true},
			{Name: "deleted_at", Type: "dateTime", Nullable: true},
		}
		mc.Fields = append(mc.Fields, timeStampFields...)
	}
	return &MigrationGenerator{
		fmt.Sprintf("create_%s_table", mc.TableName),
		mc.TableName,
		mc.Fields,
		version,
		mc.PrimaryColumns,
		mc.UniqueColumns,
		mc.ForeignColumns,
	}
}

func (mg *MigrationGenerator) BumpVersion() *MigrationGenerator {
	intVersion, _ := strconv.Atoi(mg.version)
	mg.version = fmt.Sprintf("%d", intVersion+1)
	return mg
}

func (mg *MigrationGenerator) GetReplacables() []*Replacable {
	var fieldLines string
	for index, f := range mg.fields {
		typeString := "\tt.%s(\"%s\""
		switch f.Type {
		case "string":
			typeString += ", 255)"
		case "dateTime":
			typeString += ", 0)"
		case "decimal":
			typeString += ", 8, 2)"
		default:
			typeString += ")"
		}

		if f.ForeignConstrained {
			fieldLines += fmt.Sprintf("\tt.ForeignID(\"%s\").Constrained()", f.Name)
		} else if f.Primary {
			fieldLines += fmt.Sprintf("\tt.BigIncrements(\"%s\").Primary()", f.Name)
		} else {
			fieldLines += fmt.Sprintf(typeString, strcase.ToCamel(f.Type), f.Name)
		}

		if f.Unique {
			fieldLines += ".Unique()"
		}

		if f.Nullable {
			fieldLines += ".Nullable()"
		}

		if index < len(mg.fields)-1 {
			fieldLines += "\n"
		}
	}

	if len(mg.primaryColumns) > 0 {
		primaryKeyString := fmt.Sprintf("\tt.PrimaryKey(")
		for i, c := range mg.primaryColumns {
			prefix := "\"%s\""
			if i > 0 {
				prefix = ", \"%s\""
			}
			primaryKeyString += fmt.Sprintf(prefix, c)
		}
		fieldLines += "\n" + primaryKeyString + ")"
	}

	if len(mg.uniqueColumns) > 0 && len(mg.uniqueColumns[0]) > 0 {
		for _, c := range mg.uniqueColumns {
			uniqueKeyString := fmt.Sprintf("\tt.UniqueKey(")
			for j, c2 := range c {
				prefix := "\"%s\""
				if j > 0 {
					prefix = ", \"%s\""
				}
				uniqueKeyString += fmt.Sprintf(prefix, c2)
			}
			fieldLines += "\n" + uniqueKeyString + ")"
		}
	}

	if len(mg.foreignColumns) > 0 && len(mg.foreignColumns[0]) > 0 {
		for _, columns := range mg.foreignColumns {
			foreignKeyString := fmt.Sprintf("\tt.Foreign(")
			for j, column := range columns {
				prefix := "\"%s\""
				if j > 0 {
					prefix = ", \"%s\""
				}
				foreignKeyString += fmt.Sprintf(prefix, column)
				table := guessPluralizedTableNameFromColumnName(column)
				suffix := fmt.Sprintf(").References(\"id\").On(\"%s\");", table)
				foreignKeyString += suffix
			}
			fieldLines += "\n" + foreignKeyString
		}
	}

	return []*Replacable{
		{Placeholder: "Name", Value: mg.name},
		{Placeholder: "TableName", Value: mg.tableName},
		{Placeholder: "Fields", Value: fieldLines},
		{Placeholder: "Version", Value: mg.version},
	}
}

func (mg *MigrationGenerator) GetPackagePath() string {
	return "cmd/migrations"
}

func (mg *MigrationGenerator) GetStub() string {
	return migrationStub
}

func (mg *MigrationGenerator) Generate() error {
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

	err = fs.Write(mg.GetPackagePath()+"/"+mg.version+"_"+mg.name+".go", []byte(output))

	if err != nil {
		return err
	}

	return nil
}

var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Generate a simple migration file",
	Long:  `Generate a simple migration file`,
	Run: func(cmd *cobra.Command, args []string) {
		var tableName string
		var fields []*MigrationField

		primaryColumns := []string{}
		uniqueColumns := []string{}
		foreignColumns := []string{}
		timestamps := false
		selectedPrimaryColumns := []string{}
		selectedUniqueColumns := []string{}
		selectedForeignColumns := []string{}

		nameForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Enter the table name in snake_case and plular form").
					Value(&tableName).
					Validate(cmder.SnakeCase),
			),
		)

		err := nameForm.Run()
		if err != nil {
			return
		}

		for {
			var fieldName, fieldType string
			var selectedAttrs []string

			fieldNameForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Enter the field name in snake_case").
						Value(&fieldName).
						Validate(cmder.SnakeCaseEmptyAllowed),
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
						Title("Enter the data type").
						Options(huh.NewOptions(migrationFieldTypes...)...).
						Value(&fieldType),
					huh.NewMultiSelect[string]().
						Title("Press x to select the attributes that apply to this field").
						Options(huh.NewOptions("Nullable", "Unique")...).
						Value(&selectedAttrs),
				),
			)

			err = fieldTypeForm.Run()
			if err != nil {
				return
			}

			fields = append(fields, &MigrationField{
				Name:     fieldName,
				Type:     fieldType,
				Nullable: slices.Contains(selectedAttrs, "Nullable"),
				Unique:   slices.Contains(selectedAttrs, "Unique"),
			})

			primaryColumns = append(primaryColumns, fieldName)
			uniqueColumns = append(uniqueColumns, fieldName)
			foreignColumns = append(foreignColumns, fieldName)
		}

		constraintForm := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Press x to select the primary keys").
					Options(huh.NewOptions(primaryColumns...)...).
					Value(&selectedPrimaryColumns),
				huh.NewMultiSelect[string]().
					Title("Press x to select the unique keys").
					Options(huh.NewOptions(uniqueColumns...)...).
					Value(&selectedUniqueColumns),
				huh.NewMultiSelect[string]().
					Title("Press x to select the foreign keys").
					Options(huh.NewOptions(foreignColumns...)...).
					Value(&selectedForeignColumns),
				huh.NewConfirm().
					Title("Do you want timestamp fields (created_at, updated_at, deleted_at) ?").
					Value(&timestamps),
			),
		)

		err = constraintForm.Run()
		if err != nil {
			return
		}

		mg := NewMigrationGenerator(&MigrationConfig{
			TableName:      tableName,
			Fields:         fields,
			PrimaryColumns: selectedPrimaryColumns,
			UniqueColumns:  [][]string{selectedUniqueColumns},
			ForeignColumns: [][]string{selectedForeignColumns},
			Timestamps:     timestamps,
		})
		err = mg.Generate()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Migration generated successfully.")
	},
}

func guessPluralizedTableNameFromColumnName(columnName string) string {
	pluralize := pluralize.NewClient()
	if strings.HasSuffix(columnName, "id") {
		nameParts := strings.Split(columnName, "_")
		if len(nameParts) > 1 {
			return pluralize.Plural(nameParts[len(nameParts)-2])
		}
		return pluralize.Plural(nameParts[0])
	}
	return pluralize.Plural(columnName)
}
