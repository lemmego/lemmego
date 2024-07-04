package cmd

import (
	_ "embed"
	"fmt"
	"pressebo/api/cmder"
	"pressebo/api/fsys"
	"slices"
	"strconv"
	"strings"
	"time"

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

		primaryColumns := []*cmder.Item{}
		uniqueColumns := []*cmder.Item{}
		foreignColumns := []*cmder.Item{}
		timestamps := false
		selectedPrimaryColumns := []string{}
		selectedUniqueColumns := []string{}
		selectedForeignColumns := []string{}

		cmder.Ask("Enter the table name in snake_case", cmder.SnakeCase).Fill(&tableName).
			AskRepeat("Enter the field name in snake_case", cmder.SnakeCaseEmptyAllowed, func(result any) cmder.Prompter {
				selectedAttrs := []string{}
				selectedType := ""
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
				foreignColumns = append(foreignColumns, &cmder.Item{Label: result.(string)})

				return prompt
			}).
			MultiSelect("Select the column(s) for the primary key", primaryColumns, 0).Fill(&selectedPrimaryColumns).
			MultiSelect("Select the column(s) for the unique key", uniqueColumns, 0).Fill(&selectedUniqueColumns).
			MultiSelect("Select the column(s) for the foreign key", foreignColumns, 0).Fill(&selectedForeignColumns).
			Confirm("Do you want timestamp fields (created_at, updated_at, deleted_at) ?", 'y').Fill(&timestamps)

		mg := NewMigrationGenerator(&MigrationConfig{
			TableName:      tableName,
			Fields:         fields,
			PrimaryColumns: selectedPrimaryColumns,
			UniqueColumns:  [][]string{selectedUniqueColumns},
			ForeignColumns: [][]string{selectedForeignColumns},
			Timestamps:     timestamps,
		})
		err := mg.Generate()
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
