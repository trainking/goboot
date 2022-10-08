package utils

import (
	"github.com/spf13/viper"
)

// LoadConfigFileViper 从配置文件加载配置到viper
func LoadConfigFileViper(path string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		return v, err
	}
	return v, nil
}
