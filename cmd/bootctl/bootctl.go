package main

import (
	"log"
	"os"

	"github.com/trainking/goboot/cmd/bootctl/conf"
	"github.com/trainking/goboot/cmd/bootctl/internal/gameapi"
	"github.com/trainking/goboot/cmd/bootctl/internal/httpapi"
	"github.com/trainking/goboot/cmd/bootctl/internal/initcmd"
	"github.com/trainking/goboot/cmd/bootctl/internal/service"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "bootctl",
		Usage: "A goboot ctl",
		Commands: []*cli.Command{
			initcmd.CMD(),
			httpapi.CMD(),
			gameapi.CMD(),
			service.CMD(),
		},
	}

	// 加载配置
	conf.InitConf()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
