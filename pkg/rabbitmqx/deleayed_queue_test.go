package rabbitmqx

import (
	"fmt"
	"testing"
	"time"

	"github.com/wagslane/go-rabbitmq"
)

func TestNewDelayedQueue(t *testing.T) {
	config := NewConfigByMap(map[string]interface{}{
		"mqurl":      "amqp://admin:123456@192.168.33.10:5672/app",
		"exchange":   "test.deleayed.ex",
		"routingkey": "test.deleayed.rk",
		"queue":      "test.deleayed.queue",
	})
	d, err := NewDelayedQueueConsumer(config)
	if err != nil {
		t.Error(err)
	}

	d.Consume(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		fmt.Println(d.Body)
		return rabbitmq.Ack
	})

	time.Sleep(20 * time.Second)
}

func TestPublishDelayed(t *testing.T) {
	p, err := NewPublisher("amqp://admin:123456@192.168.33.10:5672/app")
	if err != nil {
		t.Error(err)
	}
	p.SetExchange("test.deleayed.ex")
	p.SetRoutingKey([]string{"test.deleayed.rk"})
	if err != nil {
		t.Error(err)
	}

	if err := p.PublishDelayed([]byte("11111"), 3000); err != nil {
		t.Error(err)
	}
}
