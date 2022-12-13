package rabbitmqx

import (
	"github.com/wagslane/go-rabbitmq"
)

type Publisher struct {
	publisher *rabbitmq.Publisher

	routingKey []string
	options    []func(*rabbitmq.PublishOptions)
}

// NewPublisher 新建自定义的发送者
func NewPublisher(mqUrl string) (*Publisher, error) {
	publisher, err := rabbitmq.NewPublisher(mqUrl, rabbitmq.Config{}, rabbitmq.WithPublisherOptionsLogging)
	if err != nil {
		return nil, err
	}

	return &Publisher{publisher: publisher, options: []func(*rabbitmq.PublishOptions){
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
	}}, nil
}

// NewPublisherByConfig 通过Config创建Publisher
func NewPublisherByConfig(c *Config) (*Publisher, error) {
	publisher, err := NewPublisher(c.MqUrl)
	if err != nil {
		return nil, err
	}

	if c.Exchange != "" {
		publisher.SetExchange(c.Exchange)
	}
	if c.RoutingKey != "" {
		publisher.SetRoutingKey([]string{c.RoutingKey})
	}

	return publisher, nil
}

// SetExchange 设置交换机
func (p *Publisher) SetExchange(exchange string) {
	p.options = append(p.options, rabbitmq.WithPublishOptionsExchange(exchange))
}

// SetRoutingKey 设置路由键
func (p *Publisher) SetRoutingKey(routingKey []string) {
	p.routingKey = routingKey
}

// SetOptions 设置Options
func (p *Publisher) SetOptions(options []func(*rabbitmq.PublishOptions)) {
	p.options = options
}

// Publish 推送消息
func (p *Publisher) Publish(msg []byte) error {

	return p.publisher.Publish(
		msg,
		p.routingKey,
		p.options...,
	)
}

// PublishDelayed 发送延迟消息，ttl是毫秒为单位
func (p *Publisher) PublishDelayed(msg []byte, ttl int64) error {
	_delayedPoption := rabbitmq.WithPublishOptionsHeaders(rabbitmq.Table{
		"x-delay": ttl,
	})
	options := append(p.options, _delayedPoption)

	return p.publisher.Publish(
		msg,
		p.routingKey,
		options...,
	)
}

// Close 关闭消息推送者
func (p *Publisher) Close() error {
	return p.publisher.Close()
}
