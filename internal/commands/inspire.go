package commands

import (
	"fmt"
	"github.com/lemmego/api/app"
	"github.com/spf13/cobra"
	"math/rand"
	"time"
)

func init() {
	app.RegisterCommands(InspireCommand)
}

var InspireCommand = func(a app.App) *cobra.Command {
	return &cobra.Command{
		Use: "inspire",
		Run: func(cmd *cobra.Command, args []string) {
			rand.Seed(time.Now().UnixNano())
			randomIndex := rand.Intn(len(quotes()))
			randomItem := quotes()[randomIndex]
			fmt.Println(randomItem)
		},
	}
}

func quotes() []string {
	return []string{
		"\"It takes courage to grow up and become who you really are.\" — E.E. Cummings",
		"\"Your self-worth is determined by you. You don't have to depend on someone telling you who you are.\" — Beyoncé",
		"\"Nothing is impossible. The word itself says 'I'm possible!'\" — Audrey Hepburn",
		"\"Keep your face always toward the sunshine, and shadows will fall behind you.\" — Walt Whitman",
		"\"Attitude is a little thing that makes a big difference.\" — Winston Churchill",
	}
}
