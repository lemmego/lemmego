package commands

import (
	"github.com/lemmego/api/app"
)

func Load() []app.Command {
	return []app.Command{
		InspireCommand,
	}
}
