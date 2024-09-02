package main

import (
	"github.com/lemmego/api/cli"
	"github.com/lemmego/lemmego/internal/plugins/auth/cmd"
)

func main() {
	cli.RootCmd.AddCommand(cmd.AuthCmd)
	err := cli.Execute()
	if err != nil {
		panic(err)
	}
}
