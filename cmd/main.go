package main

import (
	baseCmd "lemmego/api/cmd"
	"lemmego/internal/plugins/auth/cmd"
)

func main() {
	baseCmd.RootCmd.AddCommand(cmd.AuthCmd)
	err := baseCmd.Execute()
	if err != nil {
		panic(err)
	}
}
