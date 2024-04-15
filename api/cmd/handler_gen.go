package cmd

import (
	_ "embed"
	"fmt"
	"pressebo/api/cmder"
	"pressebo/api/fsys"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

//go:embed handler.txt
var handlerStub string

type HandlerField struct {
	Name string
}

type HandlerConfig struct {
	Name string
}

type HandlerGenerator struct {
	name string
}

func NewHandlerGenerator(mc *HandlerConfig) *HandlerGenerator {
	return &HandlerGenerator{mc.Name}
}

func (mg *HandlerGenerator) GetReplacables() []*Replacable {
	return []*Replacable{
		{Placeholder: "HandlerName", Value: strcase.ToCamel(mg.name)},
	}
}

func (mg *HandlerGenerator) GetPackagePath() string {
	return "internal/handlers"
}

func (mg *HandlerGenerator) GetStub() string {
	return handlerStub
}

func (mg *HandlerGenerator) Generate() error {
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

	err = fs.Write(mg.GetPackagePath()+"/"+mg.name+"_handlers.go", []byte(output))

	if err != nil {
		return err
	}

	return nil
}

var handlerCmd = &cobra.Command{
	Use:   "handlers",
	Short: "Generate a handler set",
	Long:  `Generate a handler set`,
	Run: func(cmd *cobra.Command, args []string) {
		var handlerName string

		cmder.Ask("Enter the resource name in snake_case", cmder.SnakeCase).Fill(&handlerName)

		mg := NewHandlerGenerator(&HandlerConfig{Name: handlerName})
		err := mg.Generate()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Handler generated successfully.")
	},
}
