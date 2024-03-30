package cmd

import (
	"bytes"
	"errors"
	"html/template"
	"log"

	"github.com/spf13/cobra"
)

type GeneratorCommand interface {
	GetStubPath() string
	GetPackagePath() string
	GetReplacables() []*Replacable
	Generate() error
}

func ParseTemplate(tmplData any, fileContents string) (string, error) {
	var out bytes.Buffer
	tx := template.New("template")
	t := template.Must(tx.Parse(fileContents))
	err := t.Execute(&out, tmplData)
	if err != nil {
		return "", errors.New("Unable to execute template:" + err.Error())
	}
	return out.String(), nil
}

// genCmd represents the generator command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate code",
	Long:  `Generate code`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("gen called")
	},
}
