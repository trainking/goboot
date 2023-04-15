package gameapi

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/trainking/goboot/cmd/bootctl/temptools"
	"github.com/urfave/cli/v2"
)

var CMD = func() *cli.Command {

	var name string    // 名称
	var option string  // 操作符
	var addr string    // 监听地址
	var instanceID int // 实例ID
	var network string // 使用协议
	var natsUrl string // nats地址
	var etcdUrl string // etcd地址，多实例用 `,`分割

	return &cli.Command{
		Name:  "game",
		Usage: "create a game api",
		Action: func(ctx *cli.Context) error {
			if name == "" {
				return errors.New("must a name")
			}

			switch option {
			case "c":
				create(name, addr, instanceID, network, natsUrl, etcdUrl)
			case "g":
				generate(name)
			default:
				return errors.New("wrong option")
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
				Name:        "option, o",
				Usage:       "option code, eg. c create; g gerante proto;",
				Destination: &option,
				Value:       "c",
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
			&cli.StringFlag{
				Name:        "network",
				Usage:       "game server use network",
				Destination: &network,
				Value:       "kcp",
			},
			&cli.StringFlag{
				Name:        "nats",
				Usage:       "nats address",
				Destination: &natsUrl,
				Value:       "nats://127.0.0.1:4222",
			},
			&cli.StringFlag{
				Name:        "etcd",
				Usage:       "etcd address, use , splite",
				Destination: &etcdUrl,
				Value:       "127.0.0.1:2379",
			},
		},
	}
}

// create 创建
func create(name string, addr string, instanceID int, network string, natsUrl string, etcdUrl string) error {
	// 创建api的目录
	if err := os.MkdirAll(filepath.Join("internal", "api", name), os.ModePerm); err != nil {
		return err
	}

	// 创建main文件
	if err := createMainFile(name, addr, instanceID); err != nil {
		return err
	}

	// 创建配置文件
	if err := createConfigYaml(name, network, natsUrl, etcdUrl); err != nil {
		return err
	}

	// 创建opcode的proto
	if err := createOpProto(name, nil); err != nil {
		return err
	}

	// 创建game.proto
	if err := createGameProto(name); err != nil {
		return err
	}

	// mod tidy
	gomodCmd := exec.Command("go", "mod", "tidy")
	_, err := gomodCmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func readGameProto(name string) ([]string, error) {
	var result []string
	f, err := os.Open(filepath.Join("internal", "pb", "proto", fmt.Sprintf("%s.game.proto", name)))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		// 读取到消息标志
		if strings.Index(line, "@Message") > 0 {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			msgR := strings.Split(line, " ")
			if len(msgR) > 2 {
				result = append(result, msgR[1])
			}
		}
	}

	return result, nil
}

// generate 生成op.proto
func generate(name string) error {
	opcodArr, err := readGameProto(name)
	if err != nil {
		return err
	}

	var codeSB strings.Builder
	for i, code := range opcodArr {
		codeSB.WriteString(`  `)
		codeSB.WriteString("OP_")
		codeSB.WriteString(code)
		codeSB.WriteString(" = ")
		codeSB.WriteString(strconv.Itoa(i + 1))
		codeSB.WriteString(";\n")
	}

	if err := createOpProto(name, map[string]interface{}{
		"Codes": codeSB.String(),
	}); err != nil {
		return err
	}

	protocCmd1 := exec.Command("protoc", "--go_out=./internal/", fmt.Sprintf("./internal/pb/proto/%s.op.proto", name))
	if _, err := protocCmd1.Output(); err != nil {
		return err
	}

	protocCmd2 := exec.Command("protoc", "--go_out=./internal/", fmt.Sprintf("./internal/pb/proto/%s.game.proto", name))
	if _, err := protocCmd2.Output(); err != nil {
		return err
	}

	return nil
}

// createMainFile 创建main入口go文件
func createMainFile(name, addr string, instanceID int) error {
	mainFile, err := os.Create(filepath.Join("internal", "api", name, fmt.Sprintf("%s.api.go", name)))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	// 写入main文件
	text := temptools.ReplaceTemplate(mainText, map[string]interface{}{
		"name": name,
		"addr": addr,
		"id":   strconv.Itoa(instanceID),
	})
	if _, err := mainFile.WriteString(text); err != nil {
		return err
	}

	return nil
}

// createConfigYaml 创建配置文件yaml
func createConfigYaml(name string, network string, natsUrl string, etcdUrl string) error {
	f, err := os.Create(filepath.Join("configs", fmt.Sprintf("%s.api.yml", name)))
	if err != nil {
		return err
	}
	defer f.Close()

	var etcdArr = strings.Split(etcdUrl, ",")
	var etcdSB strings.Builder
	for _, url := range etcdArr {
		etcdSB.WriteString(`  - "`)
		etcdSB.WriteString(url)
		etcdSB.WriteString(`"`)
		etcdSB.WriteByte(10) // 换行
	}

	text := temptools.ReplaceTemplate(ymlText, map[string]interface{}{
		"network": network,
		"nats":    natsUrl,
		"etcd":    etcdSB.String(),
	})
	if _, err := f.WriteString(text); err != nil {
		return err
	}

	return nil
}

// createOpProto 创建opcode的protobuf
func createOpProto(name string, excm map[string]interface{}) error {
	f, err := os.Create(filepath.Join("internal", "pb", "proto", fmt.Sprintf("%s.op.proto", name)))
	if err != nil {
		return err
	}
	defer f.Close()

	text := temptools.ReplaceTemplate(opPorotText, excm)
	if _, err := f.WriteString(text); err != nil {
		return err
	}

	return nil
}

// createGameProto 创建game.proto
func createGameProto(name string) error {
	f, err := os.Create(filepath.Join("internal", "pb", "proto", fmt.Sprintf("%s.game.proto", name)))
	if err != nil {
		return err
	}
	defer f.Close()

	text := temptools.ReplaceTemplate(gameProtoText, map[string]interface{}{
		"name": name,
	})
	if _, err := f.WriteString(text); err != nil {
		return err
	}

	return nil
}
