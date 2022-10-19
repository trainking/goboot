package main

import (
	"flag"
	"fmt"

	"github.com/trainking/goboot/pkg/rabbitmqx"
	"github.com/wagslane/go-rabbitmq"
)

var instanceID = flag.Int("ins", 0, "queue instace id")

func main() {
	flag.Parse()

	consumer, err := rabbitmqx.NewConsumer("amqp://admin:123456@192.168.33.10:5672/app")
	if err != nil {
		panic(err)
	}

	consumer.SetExchange("test", "direct")
	consumer.SetQueue(fmt.Sprintf("test.Q.%d", *instanceID))
	consumer.SetRoutingKey([]string{"test.rk"})

	consumer.Consume(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		fmt.Println(string(d.Body))
		return rabbitmq.Ack
	})

	select {}
}
