package rabbitmqx

import (
	"fmt"
	"github.com/wagslane/go-rabbitmq"
	"testing"
	"time"
)

func TestNewDelayedQueue(t *testing.T) {
	d, err := NewDelayedQueueConsumer("amqp://admin:123456@192.168.33.10:5672/app", "test.deleayed.ex", "test.deleayed.rk", "test.deleayed.queue")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(1111)
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
