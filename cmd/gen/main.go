package main

import (
	"pressebo/api/cmd"
)

// Create a new command
// var rootCmd = &cobra.Command{
// 	Use:   "",
// 	Short: "pressebo",
// 	Long:  `p is a CLI tool for managing the pressebo application.`,
// }

// Execute the command
// func Execute() error {
// 	return cmd.Execute()
// }

func main() {
	cmd.Execute()
}
