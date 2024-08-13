package cli

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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
	"file":     reflect.String.String(),
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
	"file":     "string",
}

var commonFuncs = template.FuncMap{
	"contains":  strings.Contains,
	"hasSuffix": strings.HasSuffix,
	"join":      strings.Join,

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
	"concat": func(str string, strs ...string) string {
		return str + strings.Join(strs, "")
	},
}

type Replacable struct {
	Placeholder string
	Value       interface{}
}

type Generator interface {
	Generate() error
}

type PackagePathGetter interface {
	GetPackagePath() string
}

type StubGetter interface {
	GetStub() string
}

type Commander interface {
	Command() *cobra.Command
}

type CommandGenerator interface {
	Generator
	PackagePathGetter
	StubGetter
	Commander
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
	Use:     "gen",
	Aliases: []string{"g"},
	Short:   "Generate code",
	Long:    `Generate code`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("An argument must be provided to the gen command (e.g. model, input, migration, handlers, etc.)")
	},
}
