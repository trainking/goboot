package main

import (
	"log"
	"os"

	"github.com/trainking/goboot/cmd/codexctl/internal/generate"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "codexctl",
		Usage: "A tool for managing Codex",
		Commands: []*cli.Command{
			generate.CMD(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
