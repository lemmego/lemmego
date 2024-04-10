package cmd

import (
	_ "embed"
	"fmt"
	"pressebo/api/cmder"
	"pressebo/api/fsys"
	"slices"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

//go:embed migration.txt
var migrationStub string

var migrationFieldTypes = []string{
	"string", "text", "boolean", "int", "bigInt", "unsignedInt", "unsignedBigInt", "decimal", "date", "time",
}

type MigrationField struct {
	Name     string
	Type     string
	Nullable bool
	Unique   bool
}

type MigrationConfig struct {
	TableName      string
	Fields         []*MigrationField
	PrimaryColumns []string
	UniqueColumns  []string
}

type MigrationGenerator struct {
	name      string
	tableName string
	fields    []*MigrationField
	version   string

	primaryColumns []string
	uniqueColumns  []string
}

func NewMigrationGenerator(mc *MigrationConfig) *MigrationGenerator {
	version := time.Now().Format("20060102150405")
	return &MigrationGenerator{fmt.Sprintf("create_%s_table", mc.TableName), mc.TableName, mc.Fields, version, mc.PrimaryColumns, mc.UniqueColumns}
}

func (mg *MigrationGenerator) GetReplacables() []*Replacable {
	var fieldLines string
	for index, f := range mg.fields {
		typeString := "\tt.%s(\"%s\""
		if f.Type == "string" {
			typeString += ", 0)"
		} else {
			typeString += ")"
		}
		fieldLines += fmt.Sprintf(typeString, strcase.ToCamel(f.Type), f.Name)
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

	if len(mg.uniqueColumns) > 0 {
		uniqueKeyString := fmt.Sprintf("\tt.UniqueKey(")
		for i, c := range mg.uniqueColumns {
			prefix := "\"%s\""
			if i > 0 {
				prefix = ", \"%s\""
			}
			uniqueKeyString += fmt.Sprintf(prefix, c)
		}
		fieldLines += "\n" + uniqueKeyString + ")"
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

	tmplData := map[string]string{
		"PackageName": packageName,
	}

	for _, v := range mg.GetReplacables() {
		tmplData[v.Placeholder] = v.Value
	}

	output, err := ParseTemplate(tmplData, mg.GetStub())

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

		primaryColumns := []*cmder.Item{}
		uniqueColumns := []*cmder.Item{}
		selectedPrimaryColumns := []string{}
		selectedUniqueColumns := []string{}

		cmder.Ask("Enter the table name in snake_case", cmder.SnakeCase).Fill(&tableName).
			AskRecurring("Enter the field name in snake_case", cmder.SnakeCaseEmptyAllowed, func(result any) cmder.Prompter {
				selectedType := ""
				selectedAttrs := []string{}
				prompt := cmder.Select("What should the data type be?", migrationFieldTypes).Fill(&selectedType).
					MultiSelect("Select all the attributes that apply to this column", []*cmder.Item{
						{Label: "Nullable"}, {Label: "Unique"},
					}, 0).Fill(&selectedAttrs)

				fields = append(fields, &MigrationField{
					Name:     result.(string),
					Type:     selectedType,
					Nullable: slices.Contains(selectedAttrs, "Nullable"),
					Unique:   slices.Contains(selectedAttrs, "Unique"),
				})

				primaryColumns = append(primaryColumns, &cmder.Item{Label: result.(string)})
				uniqueColumns = append(uniqueColumns, &cmder.Item{Label: result.(string)})

				return prompt
			}).
			MultiSelect("Select the column(s) for the primary key", primaryColumns, 0).Fill(&selectedPrimaryColumns).
			MultiSelect("Select the column(s) for the unique key", uniqueColumns, 0).Fill(&selectedUniqueColumns)

		mg := NewMigrationGenerator(&MigrationConfig{TableName: tableName, Fields: fields, PrimaryColumns: selectedPrimaryColumns, UniqueColumns: selectedUniqueColumns})
		err := mg.Generate()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Migration generated successfully.")
	},
}
