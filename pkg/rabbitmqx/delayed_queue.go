package rabbitmqx

import "github.com/wagslane/go-rabbitmq"

// NewDelayedQueueConsumer 新建一个延迟队列的消费者。使用rabbitmq_delayed_message_exchange插件实现的延迟队列
// 因此需要安装对应版本的插件
func NewDelayedQueueConsumer(mqUrl string, exchange string, routingKey string, queue string) (*Consumer, error) {
	consumer, err := NewConsumer(mqUrl)
	if err != nil {
		return nil, err
	}
	consumer.SetExchange(exchange, "x-delayed-message")
	consumer.SetRoutingKey([]string{routingKey})
	consumer.SetQueue(queue)

	consumer.SetOptions([]func(*rabbitmq.ConsumeOptions){
		rabbitmq.WithConsumeOptionsBindingExchangeArgs(rabbitmq.Table{
			"x-delayed-type": "direct",
		}),
	})

	return consumer, nil
}
