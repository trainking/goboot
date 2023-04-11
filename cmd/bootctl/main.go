package main

import (
	"log"
	"os"

	"github.com/trainking/goboot/cmd/bootctl/initcmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "bootctl",
		Usage: "A goboot ctl",
		Commands: []*cli.Command{
			initcmd.CMD(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
