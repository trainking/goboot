package generate

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/trainking/goboot/cmd/bootctl/temptools"
	"github.com/urfave/cli/v2"
)

var CMD = func() *cli.Command {

	// 输出codex包的路径
	var path string
	// code配置文件的路径
	var codeFile string

	return &cli.Command{
		Name:  "generate",
		Usage: "generate a codex.",
		Action: func(ctx *cli.Context) error {
			if path == "" {
				return errors.New("must a path")
			}

			if codeFile == "" {
				return errors.New("must a code file")
			}

			return Generate(path, codeFile)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Aliases:     []string{"p"},
				Usage:       "codex package path",
				Destination: &path,
			},
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "code file path",
				Destination: &codeFile,
			},
		},
	}
}

// Generate 生成codex
func Generate(path, codeFile string) error {
	rows, err := readCodeFile(codeFile)
	if err != nil {
		return err
	}

	length := len(rows)
	var codeText strings.Builder
	var msgText strings.Builder
	var keys strings.Builder
	for i, row := range rows {
		codeText.WriteString(row.Code)
		msgText.WriteString("\"" + row.Msg + "\"")
		if i == 0 {
			keys.WriteString("\t" + row.Key + " Code = iota" + "\t" + "// " + row.Code + " " + row.Desc)
		} else {
			keys.WriteString("\n\t" + row.Key + "\t" + "// " + row.Code + " " + row.Desc)
		}

		if i < length-1 {
			codeText.WriteString(", ")
			msgText.WriteString(", ")
		}
	}

	fmt.Printf("Length: %d Code: %s Msg: %s Keys: %s\n", length, codeText.String(), msgText.String(), keys.String())

	if err := os.MkdirAll(filepath.Join(path, "codex"), os.ModePerm); err != nil {
		return err
	}

	codexFile, err := os.Create(filepath.Join(path, "codex", "codex.go"))
	if err != nil {
		return err
	}
	defer codexFile.Close()

	text := temptools.ReplaceTemplate(template, map[string]interface{}{
		"length": strconv.Itoa(length),
		"code":   codeText.String(),
		"msg":    msgText.String(),
		"keys":   keys.String(),
	})
	if _, err := codexFile.WriteString(text); err != nil {
		return err
	}

	return nil
}

type Row struct {
	Key  string
	Code string
	Msg  string
	Desc string
}

// readCodeFile 读取code.csv文件
func readCodeFile(path string) ([]Row, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rows []Row
	reader := csv.NewReader(file)
	// 跳过第一行（标题行）
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		rows = append(rows, Row{
			Key:  row[0],
			Code: row[1],
			Msg:  row[2],
			Desc: row[3],
		})
	}

	return rows, nil
}
