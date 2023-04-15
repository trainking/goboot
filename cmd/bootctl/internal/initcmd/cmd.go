package initcmd

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/trainking/goboot/cmd/bootctl/conf"
	"github.com/urfave/cli/v2"
)

var dirs []string = []string{"bin", "cmd", "configs", "docs", "internal", "pkg"}

var CMD = func() *cli.Command {

	// 项目名
	var name string

	return &cli.Command{
		Name:  "init",
		Usage: "init a project.",
		Action: func(ctx *cli.Context) error {

			if name == "" {
				return errors.New("must a name")
			}

			// 创建bootrc配置文件
			if err := conf.WriteConf(conf.Conf{
				ModName: name,
			}); err != nil {
				return err
			}

			// 创建基础目录
			for _, dir := range dirs {
				if err := os.MkdirAll(dir, os.ModePerm); err != nil {
					return err
				}
			}

			if err := os.MkdirAll(filepath.Join("internal", "pb", "proto"), os.ModePerm); err != nil {
				return err
			}

			gitignFile, err := os.Create(".gitignore")
			if err != nil {
				return err
			}
			defer gitignFile.Close()

			if _, err := gitignFile.WriteString(gitignoreText); err != nil {
				return err
			}

			gomodCmd := exec.Command("go", "mod", "init", name)
			_, err = gomodCmd.Output()
			if err != nil {
				return err
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Usage:       "project name",
				Destination: &name,
			},
		},
	}
}
