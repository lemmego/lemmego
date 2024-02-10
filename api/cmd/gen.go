package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// genCmd represents the generator command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate code",
	Long:  `Generate code`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("gen called")
	},
}
