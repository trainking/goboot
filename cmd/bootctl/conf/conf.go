package conf

import (
	"os"

	"gopkg.in/yaml.v3"
)

const BOOT_RC = "./.bootrc"

var _conf *Conf

type Conf struct {
	ModName string `yaml:"mod_name"`
}

// InitConf 初始化配置
func InitConf() {
	data, err := os.ReadFile(BOOT_RC)
	if err != nil {
		return
	}

	var conf Conf
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		panic(err)
	}

	_conf = &conf
}

// WriteConf 写入配置文件
func WriteConf(c Conf) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(BOOT_RC, data, os.ModePerm)
	return err
}

// GetConf 获取配置
func GetConf() *Conf {
	return _conf
}
