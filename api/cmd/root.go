package cmd

// Import the cobra library
import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"pressebo/internal/plugins/auth/cmd"
)

// Create a new command
var RootCmd = &cobra.Command{
	Use:   "",
	Short: fmt.Sprintf("%s", os.Getenv("APP_NAME")),
}

// Execute the command
func Execute() error {
	RootCmd.AddCommand(genCmd)

	genCmd.AddCommand(cmd.AuthCmd)
	return RootCmd.Execute()
}