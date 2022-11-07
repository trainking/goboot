package rabbitmqx

import (
	"github.com/wagslane/go-rabbitmq"
)

// NewDelayedQueuePulisher 创建一个publisher，使用DeleayedQueueConfig的配置
func NewDelayedQueuePulisher(config *Config) (*Publisher, error) {
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
func NewDelayedQueueConsumer(config *Config) (*Consumer, error) {
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
