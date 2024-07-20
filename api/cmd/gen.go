package cmd

import (
	"bytes"
	"errors"
	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

var UiDataTypeMap = map[string]string{
	"text":     reflect.String.String(),
	"textarea": reflect.String.String(),
	"integer":  reflect.Uint.String(),
	"decimal":  reflect.Float64.String(),
	"boolean":  reflect.Bool.String(),
	"radio":    reflect.String.String(),
	"checkbox": reflect.String.String(),
	"dropdown": reflect.String.String(),
	"date":     "time.Time",
	"time":     "time.Time",
	"image":    reflect.String.String(),
}

var UiDbTypeMap = map[string]string{
	"text":     "string",
	"textarea": "text",
	"integer":  "unsignedBigInt",
	"decimal":  "decimal",
	"boolean":  "boolean",
	"radio":    "string",
	"checkbox": "string",
	"dropdown": "string",
	"date":     "dateTime",
	"time":     "time",
	"image":    "string",
}

var commonFuncs = template.FuncMap{
	"toTitle": func(str string) string {
		caser := cases.Title(language.English)
		return caser.String(str)
	},
	"toCamel": func(str string) string {
		return strcase.ToCamel(str)
	},
	"toLowerCamel": func(str string) string {
		return strcase.ToLowerCamel(str)
	},
	"toSnake": func(str string) string {
		return strcase.ToSnake(str)
	},
	"toSpaceDelimited": func(str string) string {
		return strcase.ToDelimited(str, ' ')
	},
	"contains": strings.Contains,
}

type Replacable struct {
	Placeholder string
	Value       interface{}
}

type Generator interface {
	Generate() error
}

func ParseTemplate(tmplData map[string]interface{}, fileContents string, funcMap template.FuncMap) (string, error) {
	var out bytes.Buffer
	tx := template.New("template")
	if funcMap != nil {
		tx.Funcs(funcMap)
	}
	t := template.Must(tx.Parse(fileContents))
	err := t.Execute(&out, tmplData)
	if err != nil {
		return "", errors.New("Unable to execute template:" + err.Error())
	}
	// Replace &#34; with "
	result := strings.ReplaceAll(out.String(), "&#34;", "\"")
	return result, nil
}

// genCmd represents the generator command
var genCmd = &cobra.Command{
	Use:   "g",
	Short: "Generate code",
	Long:  `Generate code`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("An argument must be provided to the gen command (e.g. model, input, migration, handlers, etc.)")
	},
}
