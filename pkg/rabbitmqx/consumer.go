package rabbitmqx

import (
	"runtime"

	"github.com/wagslane/go-rabbitmq"
)

// Consumer 消费者
type Consumer struct {
	consumer rabbitmq.Consumer

	exchange   string
	exType     string
	queue      string
	routingKey []string
	options    []func(*rabbitmq.ConsumeOptions)
}

// NewConsumer 新建自定义的消费者
func NewConsumer(mqUrl string) (*Consumer, error) {
	consumer, err := rabbitmq.NewConsumer(
		mqUrl, rabbitmq.Config{}, rabbitmq.WithConsumerOptionsLogging)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: consumer, options: []func(*rabbitmq.ConsumeOptions){
		rabbitmq.WithConsumeOptionsBindingExchangeDurable,
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsConcurrency(runtime.NumCPU()),
	}}, nil
}

// NewConsumerByConfig 通过Config
func NewConsumerByConfig(c *Config) (*Consumer, error) {
	consumer, err := NewConsumer(c.MqUrl)
	if err != nil {
		return nil, err
	}

	if c.Exchange != "" && c.ExType != "" {
		consumer.SetExchange(c.Exchange, c.ExType)
	}
	if c.RoutingKey != "" {
		consumer.SetRoutingKey([]string{c.RoutingKey})
	}
	if c.Queue != "" {
		consumer.SetQueue(c.Queue)
	}

	return consumer, nil
}

// SetExchange 设置交换机
func (c *Consumer) SetExchange(exchange string, exType string) {
	c.exchange = exchange
	c.exType = exType
}

// SetQueue 设置消费队列
func (c *Consumer) SetQueue(queue string) {
	c.queue = queue
}

// SetRoutingKey 设置routingKey
func (c *Consumer) SetRoutingKey(routingKey []string) {
	c.routingKey = routingKey
}

// SetOptions 设置Options
func (c *Consumer) SetOptions(options []func(*rabbitmq.ConsumeOptions)) {
	c.options = options
}

// Consume 消费
func (c *Consumer) Consume(handler rabbitmq.Handler) error {
	if c.exchange != "" && c.exType != "" {
		c.options = append(c.options, rabbitmq.WithConsumeOptionsBindingExchangeName(c.exchange), rabbitmq.WithConsumeOptionsBindingExchangeKind(c.exType))
	}
	err := c.consumer.StartConsuming(
		handler, c.queue, c.routingKey, c.options...)
	return err
}

// Close 关闭
func (c *Consumer) Close() error {
	return c.consumer.Close()
}
