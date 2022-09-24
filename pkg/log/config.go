package log

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Config struct {
	Level      string `json:"level" yaml:"level"`           // 输出的日志级别, 默认debug
	Target     string `json:"target" yaml:"target"`         // 日志标识，一般使用服务名，必须
	ID         string `json:"id" yaml:"id"`                 // 与target组成唯一标识
	OutPath    string `json:"outPath" yaml:"outPath"`       // 日志文件输出路径，默认logs
	MaxSize    int    `json:"maxSize" yaml:"maxSize"`       // 文件最大大小 MB，默认50M
	CallerSkip int    `json:"callerSkip" yaml:"callerSkip"` // 跳过多少层次，找caller
}

// NewConfigByMap 从一个map中创建Config
func NewConfigByMap(m map[string]interface{}) Config {
	c := Config{}
	for k, v := range m {
		c.Set(k, v)
	}
	c.defaultChange()
	return c
}

// Set 动态赋值Config属性
func (c *Config) Set(k string, v interface{}) {
	k = strings.ToLower(k)
	switch k {
	case "level":
		c.Level = v.(string)
	case "target":
		c.Target = v.(string)
	case "id":
		c.ID = v.(string)
	case "outpath":
		c.OutPath = v.(string)
	case "maxsize":
		c.MaxSize = v.(int)
	case "callerskip":
		c.CallerSkip = v.(int)
	}
}

// LogPath 日志文件路径
func (c *Config) LogPath() string {
	// 替换掉可能导致路径问题的字符
	rep := strings.NewReplacer(
		`\`, "",
		"/", "",
		":", "_",
		"*", "_",
		"?", "_",
		"=", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return filepath.Join(c.OutPath, rep.Replace(c.Target), rep.Replace(c.ServiceID()+".log"))
}

// ServiceID 服务ID, 使用Target+ID生成
func (c *Config) ServiceID() string {
	if c.ID == "" {
		return c.Target
	}
	return fmt.Sprintf("%s_%s", c.Target, c.ID)
}

func (c *Config) defaultChange() {
	if c.Level == "" {
		c.Level = "debug"
	}
	if c.Target == "" {
		panic("target is must be")
	}
	if c.OutPath == "" {
		c.OutPath = "logs"
	}
	if c.MaxSize == 0 {
		c.MaxSize = 50
	}
	if c.CallerSkip == 0 {
		c.CallerSkip = 1
	}
}
