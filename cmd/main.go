package main

import (
	baseCmd "lemmego/api/cmd"
)

func main() {
	// baseCmd.RootCmd.AddCommand(cmd.AuthCmd)
	baseCmd.Execute()
}
