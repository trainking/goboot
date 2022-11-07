package rabbitmqx

import "strings"

// Config rabiitmq配置的抽象
type Config struct {
	MqUrl      string // 消息队列连接地址配置
	Exchange   string // 交换机
	RoutingKey string // 路由键
	Queue      string // 队列名
	ExType     string // 交换机类型
}

// NewConfigByMap 从map[string]interface{}中新建一个Config
func NewConfigByMap(setting map[string]interface{}) *Config {
	config := new(Config)

	for k, v := range setting {
		switch strings.ToLower(k) {
		case "mqurl":
			config.MqUrl = v.(string)
		case "exchange":
			config.Exchange = v.(string)
		case "routingkey":
			config.RoutingKey = v.(string)
		case "queue":
			config.Queue = v.(string)
		case "extype":
			config.ExType = v.(string)
		}
	}

	if config.ExType == "" {
		config.ExType = "direct"
	}

	return config
}
