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

var modelFieldTypes = []string{
	"int", "uint", "int64", "uint64", "float64", "string", "bool", "time.Time", "custom",
}

const (
	TagColumn                 = "column"
	TagType                   = "type"
	TagSerializer             = "serializer"
	TagSize                   = "size"
	TagPrimaryKey             = "primaryKey"
	TagUnique                 = "unique"
	TagDefault                = "default"
	TagPrecision              = "precision"
	TagScale                  = "scale"
	TagNotNull                = "not null"
	TagAutoIncrement          = "autoIncrement"
	TagAutoIncrementIncrement = "autoIncrementIncrement"
	TagEmbedded               = "embedded"
	TagEmbeddedPrefix         = "embeddedPrefix"
	TagAutoCreateTime         = "autoCreateTime"
	TagAutoUpdateTime         = "autoUpdateTime"
	TagIndex                  = "index"
	TagUniqueIndex            = "uniqueIndex"
	TagCheck                  = "check"
	TagWritePerm              = "<-"
	TagReadPerm               = "->"
	TagIgnore                 = "-"
	TagComment                = "comment"
)

type ModelField struct {
	Name               string
	Type               string
	Required           bool
	Unique             bool
	Primary            bool
	ForeignConstrained bool
}

type ModelConfig struct {
	Name   string
	Fields []*ModelField
}

type ModelGenerator struct {
	name   string
	fields []*ModelField
}

type DBTag struct {
	Name     string
	Argument string
}

type DBTagBuilder struct {
	tags       []*DBTag
	driverName string
}

func NewDBTagBuilder(tags []*DBTag, driverName string) *DBTagBuilder {
	return &DBTagBuilder{tags, driverName}
}

func (mtb *DBTagBuilder) Add(name, argument string) *DBTagBuilder {
	mtb.tags = append(mtb.tags, &DBTag{name, argument})
	return mtb
}

func (mtb *DBTagBuilder) Build() string {
	// Build the tag string in this format: gorm:"tagName1:tagArgument1,tagName2:tagArgument2".
	// If the argument is empty, it's omitted: gorm:"tagName1,tagName2".
	var tagStrs []string
	for _, t := range mtb.tags {
		if t.Argument != "" {
			tagStrs = append(tagStrs, fmt.Sprintf(`"%s:%s"`, t.Name, t.Argument))
		} else {
			tagStrs = append(tagStrs, fmt.Sprintf(`"%s"`, t.Name))
		}
	}
	if len(tagStrs) == 0 {
		return ""
	}
	return fmt.Sprintf("%s:", mtb.driverName) + strings.Join(tagStrs, ",")
}

func NewModelGenerator(mc *ModelConfig) *ModelGenerator {
	return &ModelGenerator{mc.Name, mc.Fields}
}

func (mg *ModelGenerator) GetReplacables() []*Replacable {
	var fieldLines string
	for index, f := range mg.fields {
		tb := NewDBTagBuilder(nil, "gorm")
		if f.Required {
			tb.Add(TagNotNull, "")
		}
		if f.Unique {
			tb.Add(TagUnique, "")
		}
		fieldLines += fmt.Sprintf("\t%s %s `json:\"%s\" %s`", strcase.ToCamel(f.Name), f.Type, f.Name, tb.Build())
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

	if exists, _ := fs.Exists(mg.GetPackagePath()); exists {
		err = fs.Write(mg.GetPackagePath()+"/"+mg.name+".go", []byte(output))

		if err != nil {
			return err
		}
	} else {
		fs.CreateDirectory(mg.GetPackagePath())
		err = fs.Write(mg.GetPackagePath()+"/"+mg.name+".go", []byte(output))

		if err != nil {
			return err
		}
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

		cmder.Ask("Enter the model name in singular form and in snake_case", cmder.SnakeCase).Fill(&modelName).
			AskRepeat(
				"Enter the field name in snake_case",
				cmder.NotIn(
					[]string{"id", "created_at", "updated_at", "deleted_at"},
					"This field will be provided for you",
					cmder.SnakeCaseEmptyAllowed,
				),
				func(result any) cmder.Prompter {
					const required = "Required"
					const unique = "Unique"
					selectedAttrs := []string{}
					selectedType := ""
					prompt := cmder.Select("What should the data type be?", modelFieldTypes).Fill(&selectedType).
						When(func(result interface{}) bool {
							return result.(string) == "custom"
						}, func(prompt cmder.Prompter) cmder.Prompter {
							return cmder.Ask("Enter the data type (You'll need to import it if necessary)", nil).Fill(&selectedType)
						}).
						MultiSelect("Select attribute(s) for this field:", []*cmder.Item{
							{Label: required}, {Label: unique},
						}, 0).Fill(&selectedAttrs)

					fields = append(
						fields,
						&ModelField{
							Name:     result.(string),
							Type:     selectedType,
							Required: slices.Contains(selectedAttrs, required),
							Unique:   slices.Contains(selectedAttrs, unique),
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
