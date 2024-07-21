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
	Timestamps     bool
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
		mc.Timestamps,
	}
}

func (mg *MigrationGenerator) BumpVersion() *MigrationGenerator {
	intVersion, _ := strconv.Atoi(mg.version)
	mg.version = fmt.Sprintf("%d", intVersion+1)
	return mg
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
		"PackageName":    packageName,
		"Name":           mg.name,
		"TableName":      mg.tableName,
		"Version":        mg.version,
		"Fields":         mg.fields,
		"PrimaryColumns": mg.primaryColumns,
		"UniqueColumns":  mg.uniqueColumns,
		"ForeignColumns": mg.foreignColumns,
		"Timestamps":     mg.Timestamps,
	}

	output, err := ParseTemplate(tmplData, mg.GetStub(), commonFuncs)

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
			var attrs = []string{"Nullable", "Unique", "Primary", "ForeignConstrained"}
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
						Title("Press x to mark the attributes that apply to this field").
						Options(huh.NewOptions(attrs...)...).
						Value(&selectedAttrs),
				),
			)

			err = fieldTypeForm.Run()
			if err != nil {
				return
			}

			fields = append(fields, &MigrationField{
				Name:               fieldName,
				Type:               fieldType,
				Nullable:           slices.Contains(selectedAttrs, "Nullable"),
				Unique:             slices.Contains(selectedAttrs, "Unique"),
				Primary:            slices.Contains(selectedAttrs, "Primary"),
				ForeignConstrained: slices.Contains(selectedAttrs, "ForeignConstrained"),
			})

			primaryColumns = append(primaryColumns, fieldName)
			uniqueColumns = append(uniqueColumns, fieldName)
			foreignColumns = append(foreignColumns, fieldName)
		}

		constraintForm := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Press x to mark the primary key columns").
					Options(huh.NewOptions(primaryColumns...)...).
					Value(&selectedPrimaryColumns),
				huh.NewMultiSelect[string]().
					Title("Press x to mark the unique key columns").
					Options(huh.NewOptions(uniqueColumns...)...).
					Value(&selectedUniqueColumns),
				huh.NewMultiSelect[string]().
					Title("Press x to mark the foreign key columns").
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
