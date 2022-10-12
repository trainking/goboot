package rabbitmqx

import (
	"strings"

	"github.com/wagslane/go-rabbitmq"
)

// DeleayedQueueConfig 延迟队列配置
type DeleayedQueueConfig struct {
	MqUrl      string // 消息队列连接地址配置
	Exchange   string // 交换机
	RoutingKey string // 路由键
	Queue      string // 队列名
}

// DeleayedQueueConfig 设置延迟队列配置
func NewDeleayedQueueConfigByMap(setting map[string]interface{}) *DeleayedQueueConfig {
	config := new(DeleayedQueueConfig)

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
		}
	}

	return config
}

// NewDelayedQueuePulisher 创建一个publisher，使用DeleayedQueueConfig的配置
func NewDelayedQueuePulisher(config *DeleayedQueueConfig) (*Publisher, error) {
	p, err := NewPublisher(config.MqUrl)
	if err != nil {
		return nil, err
	}

	p.SetExchange(config.Exchange)
	p.SetRoutingKey([]string{config.RoutingKey})

	return p, nil
}

// NewDelayedQueueConsumer 新建一个延迟队列的消费者。使用rabbitmq_delayed_message_exchange插件实现的延迟队列
// 因此需要安装对应版本的插件
func NewDelayedQueueConsumer(config *DeleayedQueueConfig) (*Consumer, error) {
	consumer, err := NewConsumer(config.MqUrl)
	if err != nil {
		return nil, err
	}
	consumer.SetExchange(config.Exchange, "x-delayed-message")
	consumer.SetRoutingKey([]string{config.RoutingKey})
	consumer.SetQueue(config.Queue)

	consumer.SetOptions([]func(*rabbitmq.ConsumeOptions){
		rabbitmq.WithConsumeOptionsBindingExchangeArgs(rabbitmq.Table{
			"x-delayed-type": "direct",
		}),
	})

	return consumer, nil
}
