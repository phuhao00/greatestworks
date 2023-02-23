package rabbitmq

import (
	pubsub "github.com/samber/go-amqp-pubsub"
	"github.com/samber/mo"
)

func example() {
	// `err` can be ignored since it will connect lazily to rabbitmq
	conn, _ := pubsub.NewConnection("connection-1", pubsub.ConnectionOptions{
		URI:            "amqp://dev:dev@localhost:5672",
		LazyConnection: mo.Some(true),
	})

	producer := pubsub.NewProducer(conn, "producer-1", pubsub.ProducerOptions{
		Exchange: pubsub.ProducerOptionsExchange{
			Name: "product.event",
			Kind: pubsub.ExchangeKindTopic,
		},
	})

	//err := producer.Publish(routingKey, false, false, amqp.Publishing{
	//	ContentType:  "application/json",
	//	DeliveryMode: amqp.Persistent,
	//	Body:         []byte(`{"hello": "world"}`),
	//})

	producer.Close()
	conn.Close()

}
