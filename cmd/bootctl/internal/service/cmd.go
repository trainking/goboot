package service

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/trainking/goboot/cmd/bootctl/conf"
	"github.com/trainking/goboot/cmd/bootctl/temptools"
	"github.com/urfave/cli/v2"
)

var CMD = func() *cli.Command {

	var name string    // 名称
	var addr string    // 监听地址
	var instanceID int // 实例ID

	return &cli.Command{
		Name:  "service",
		Usage: "create a gPRC service",
		Action: func(ctx *cli.Context) error {
			if name == "" {
				return errors.New("must a name")
			}

			// 创建protobuf文件
			if err := os.MkdirAll(filepath.Join("internal", "pb", "proto"), os.ModePerm); err != nil {
				return err
			}

			// 创建pb文件
			if err := createProtoFile(name); err != nil {
				return err
			}

			// 执行protoc
			protocCmd := exec.Command("protoc", "--go_out=./internal/", "--go-grpc_out=./internal/", fmt.Sprintf("./internal/pb/proto/%s.service.proto", name))
			_, err := protocCmd.Output()
			if err != nil {
				return err
			}

			// 创建service的目录
			if err := os.MkdirAll(filepath.Join("internal", "service", name), os.ModePerm); err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Join("internal", "service", name, "server"), os.ModePerm); err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Join("internal", "service", name, "client"), os.ModePerm); err != nil {
				return err
			}

			// 创建main文件
			if err := createMainFile(name, addr, instanceID); err != nil {
				return err
			}

			if err := createServerFile(name); err != nil {
				return err
			}

			if err := createClientFile(name); err != nil {
				return err
			}

			// mod tidy
			gomodCmd := exec.Command("go", "mod", "tidy")
			_, err = gomodCmd.Output()
			if err != nil {
				return err
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name, n",
				Usage:       "game server name",
				Destination: &name,
			},
			&cli.StringFlag{
				Name:        "addr, a",
				Usage:       "game server listen address",
				Destination: &addr,
				Value:       "127.0.0.1:8080",
			},
			&cli.IntFlag{
				Name:        "id, i",
				Usage:       "game server intance id",
				Destination: &instanceID,
				Value:       1,
			},
		},
	}
}

// createProtoFile 创建protobuf文件
func createProtoFile(name string) error {
	f, err := os.Create(filepath.Join("internal", "pb", "proto", fmt.Sprintf("%s.service.proto", name)))
	if err != nil {
		return err
	}
	defer f.Close()

	text := temptools.ReplaceTemplate(protoText, map[string]interface{}{
		"name": temptools.ToUpperFistring(name),
	})
	if _, err := f.WriteString(text); err != nil {
		return err
	}

	return nil
}

// createMainFile 创建main文件
func createMainFile(name, addr string, id int) error {
	f, err := os.Create(filepath.Join("internal", "service", name, fmt.Sprintf("%s.service.go", name)))
	if err != nil {
		return err
	}
	defer f.Close()

	text := temptools.ReplaceTemplate(mainText, map[string]interface{}{
		"project": conf.GetConf().ModName,
		"name":    name,
		"Name":    temptools.ToUpperFistring(name),
		"addr":    addr,
		"id":      strconv.Itoa(id),
	})
	if _, err := f.WriteString(text); err != nil {
		return err
	}

	return nil
}

// createServerFile 创建服务端文件
func createServerFile(name string) error {
	f, err := os.Create(filepath.Join("internal", "service", name, "server", "server.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	text := temptools.ReplaceTemplate(serverText, map[string]interface{}{
		"project": conf.GetConf().ModName,
		"Name":    temptools.ToUpperFistring(name),
	})

	if _, err := f.WriteString(text); err != nil {
		return err
	}
	return nil
}

// createClientFile 创建客户端文件
func createClientFile(name string) error {
	f, err := os.Create(filepath.Join("internal", "service", name, "client", "client.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	text := temptools.ReplaceTemplate(clientText, map[string]interface{}{
		"project": conf.GetConf().ModName,
		"Name":    temptools.ToUpperFistring(name),
	})

	if _, err := f.WriteString(text); err != nil {
		return err
	}
	return nil
}
