package main

import (
	baseCmd "pressebo/api/cmd"
	"pressebo/internal/plugins/auth/cmd"
)

func main() {
	baseCmd.RootCmd.AddCommand(cmd.AuthCmd)
	baseCmd.Execute()
}
