package main

import (
	"log"
	"os"

	"github.com/trainking/goboot/cmd/bootctl/internal/httpapi"
	"github.com/trainking/goboot/cmd/bootctl/internal/initcmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "bootctl",
		Usage: "A goboot ctl",
		Commands: []*cli.Command{
			initcmd.CMD(),
			httpapi.CMD(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
