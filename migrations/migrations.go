package main

import (
	"github.com/joho/godotenv"
	"github.com/lemmego/migration/cmd"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	cmd.Execute()
}
