package httapicmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var CMD = func() *cli.Command {

	var name string // 名称
	var addr string // 监听地址
	var instanceID int

	return &cli.Command{
		Name:  "httpapi",
		Usage: "create a http api",
		Action: func(ctx *cli.Context) error {
			if name == "" {
				return errors.New("must a name")
			}

			// 创建api的目录
			if err := os.MkdirAll(filepath.Join("api", name), os.ModePerm); err != nil {
				return err
			}

			mainFile, err := os.Create(filepath.Join("api", name, fmt.Sprintf("%s.api.go", name)))
			if err != nil {
				return err
			}
			defer mainFile.Close()

			// TODO 写入main文件

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name, n",
				Usage:       "http api name",
				Destination: &name,
			},
			&cli.StringFlag{
				Name:        "addr, a",
				Usage:       "http api listen address",
				Destination: &addr,
				DefaultText: "127.0.0.1:8080",
			},
			&cli.IntFlag{
				Name:        "id, i",
				Usage:       "http api intance id",
				Destination: &instanceID,
				DefaultText: "1",
			},
		},
	}
}
