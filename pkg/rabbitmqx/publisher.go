package rabbitmqx

import (
	"strconv"

	"github.com/wagslane/go-rabbitmq"
)

type Publisher struct {
	publisher *rabbitmq.Publisher

	exchange   string
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

// SetExchange 设置交换机
func (p *Publisher) SetExchange(exchange string) {
	p.exchange = exchange
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
	if p.exchange != "" {
		p.options = append(p.options, rabbitmq.WithPublishOptionsExchange(p.exchange))
	}

	return p.publisher.Publish(
		msg,
		p.routingKey,
		p.options...,
	)
}

// PublishDelayed 发送延迟消息，ttl是毫秒为单位
func (p *Publisher) PublishDelayed(msg []byte, ttl int64) error {
	p.options = append(p.options, rabbitmq.WithPublishOptionsExpiration(strconv.FormatInt(ttl, 10)))
	return p.Publish(msg)
}

// Close 关闭消息推送者
func (p *Publisher) Close() error {
	return p.publisher.Close()
}
