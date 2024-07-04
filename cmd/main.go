package main

import (
	baseCmd "pressebo/api/cmd"
)

func main() {
	// baseCmd.RootCmd.AddCommand(cmd.AuthCmd)
	baseCmd.Execute()
}
