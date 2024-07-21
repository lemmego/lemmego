package main

import (
	"lemmego/api/cli"
	"lemmego/internal/plugins/auth/cmd"
)

func main() {
	cli.RootCmd.AddCommand(cmd.AuthCmd)
	err := cli.Execute()
	if err != nil {
		panic(err)
	}
}
