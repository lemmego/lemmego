package bootstrap

import (
	"github.com/lemmego/api/app"
	"github.com/lemmego/lemmego/internal/commands"
)

func LoadCommands() []app.Command {
	return []app.Command{
		commands.AppKeyCommand,
		commands.InspireCommand,
	}
}
