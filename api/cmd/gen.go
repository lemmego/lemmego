package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

var DataTypeMap = map[string]string{
	"text":     reflect.String.String(),
	"textarea": reflect.String.String(),
	"integer":  reflect.Int.String(),
	"decimal":  reflect.Float64.String(),
	"boolean":  reflect.Bool.String(),
	"radio":    reflect.Array.String(),
	"checkbox": reflect.Array.String(),
	"dropdown": reflect.Array.String(),
	"date":     "time.Time",
	"time":     "time.Time",
	"image":    reflect.String.String(),
}

var DBTypeMap = map[string]string{
	"text":     "String",
	"textarea": "Text",
	"integer":  "UnsignedBigInt",
	"decimal":  "Decimal",
	"boolean":  "Boolean",
	"radio":    "String",
	"checkbox": "String",
	"dropdown": "String",
	"date":     "DateTime",
	"time":     "Time",
	"image":    "String",
}

type Replacable struct {
	Placeholder string
	Value       string
}

type Generator interface {
	GetStub() []byte
	GetPackagePath() string
	GetReplacables() []*Replacable
}

type GeneratorCommand struct {
	Generator
}

func (gc *GeneratorCommand) Generate() {
	fmt.Println("Generating...")
}

func NewGeneratorCommand(generator Generator) *GeneratorCommand {
	return &GeneratorCommand{generator}
}

func ParseTemplate(tmplData map[string]string, fileContents string) (string, error) {
	var out bytes.Buffer
	tx := template.New("template")
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
	Use:   "gen",
	Short: "Generate code",
	Long:  `Generate code`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("An argument must be provided to the gen command (e.g. model, input, migration, handlers, etc.)")
	},
}
