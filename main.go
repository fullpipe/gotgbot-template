package main

import (
	"bm/cmd/bot"
	"bm/cmd/graph"
	"bm/cmd/migrate"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "bot",
		Usage: "good bot",
		Commands: []*cli.Command{
			bot.NewCommand(),
			graph.NewCommand(),
			migrate.NewCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
