
# Pub/Sub framework for RabbitMQ and Go

- Based on github.com/rabbitmq/amqp091-go driver
- Resilient to network failure
- Auto reconnect: recreate channels, bindings, producers, consumers...
- Hot update of queue bindings (thread-safe)
- Retry
- Dead letter queue on message rejection

## How to

During your tests, feel free to restart Rabbitmq. This library will reconnect automatically.

### Connect

```go
import pubsub "github.com/samber/go-amqp-pubsub"

conn, err := pubsub.NewConnection("connection-1", pubsub.ConnectionOptions{
    URI: "amqp://dev:dev@localhost:5672",
    Config: amqp.Config{
        Dial:      amqp.DefaultDial(time.Second),
    },
})

// ...

conn.Close()
```

### Producer

```go
import (
    pubsub "github.com/samber/go-amqp-pubsub"
    "github.com/samber/lo"
    "github.com/samber/mo"
)

// `err` can be ignored since it will connect lazily to rabbitmq
conn, err := pubsub.NewConnection("connection-1", pubsub.ConnectionOptions{
    URI: "amqp://dev:dev@localhost:5672",
    LazyConnection: mo.Some(true),
})

producer := pubsub.NewProducer(conn, "producer-1", pubsub.ProducerOptions{
    Exchange: pubsub.ProducerOptionsExchange{
        Name: "product.event",
        Kind: pubsub.ExchangeKindTopic,
    },
})

err := producer.Publish(routingKey, false, false, amqp.Publishing{
    ContentType:  "application/json",
    DeliveryMode: amqp.Persistent,
    Body:         []byte(`{"hello": "world"}`),
})

producer.Close()
conn.Close()
```

### Consumer

```go
import (
    pubsub "github.com/samber/go-amqp-pubsub"
    "github.com/samber/lo"
    "github.com/samber/mo"
)

// `err` can be ignore since it will connect lazily to rabbitmq
conn, err := pubsub.NewConnection("connection-1", pubsub.ConnectionOptions{
    URI: "amqp://dev:dev@localhost:5672",
    LazyConnection: mo.Some(true),
})

consumer := pubsub.NewConsumer(conn, "consumer-1", pubsub.ConsumerOptions{
    Queue: pubsub.ConsumerOptionsQueue{
        Name: "product.onEdit",
    },
    Bindings: []pubsub.ConsumerOptionsBinding{
        {ExchangeName: "product.event", RoutingKey: "product.created"},
        {ExchangeName: "product.event", RoutingKey: "product.updated"},
    },
    Message: pubsub.ConsumerOptionsMessage{
        PrefetchCount: mo.Some(100),
    },
    EnableDeadLetter: mo.Some(true),     // will create a "product.onEdit.deadLetter" DL queue
})

for msg := range consumer.Consume() {
    lo.Try0(func() { // handle exceptions
        // ...
        msg.Ack(false)
    })
}

consumer.Close()
conn.Close()
```

### Consumer with pooling and batching

See [examples/consumer-with-pool-and-batch](examples/consumer-with-pool-and-batch/main.go).

### Consumer with retry strategy

![Retry architecture](doc/retry.png)

See [examples/consumer-with-retry.md](examples/consumer-with-retry/main.go).

2 retry strategies are available:
- Exponential backoff
- Constant interval

```go
consumer := pubsub.NewConsumer(conn, "example-consumer-1", pubsub.ConsumerOptions{
    Queue: pubsub.ConsumerOptionsQueue{
        Name: "product.onEdit",
    },
    // ...
    RetryStrategy:    mo.Some(pubsub.NewExponentialRetryStrategy(3, 3*time.Second, 2)), // will create a "product.onEdit.retry" queue
})

for msg := range consumer.Consume() {
    // ...
    msg.Reject(false)   // will retry 3 times with exponential backoff
}
```

#### Custom retry strategy

Custom strategies can be provided to the consumer.

```go
type MyCustomRetryStrategy struct {}

func NewMyCustomRetryStrategy() RetryStrategy {
	return &MyCustomRetryStrategy{}
}

func (rs *MyCustomRetryStrategy) NextBackOff(msg *amqp.Delivery, attempts int) (time.Duration, bool) {
    // retries every 10 seconds, until message get older than 5 minutes
    if msg.Timestamp.Add(5*time.Minute).After(time.Now()) {
        return 10 * time.Second, true
    }

    return time.Duration{}, false
}
```

#### Consistency

On retry, the message is published into the retry queue then is acked from the initial queue. This 2 phases delivery is unsafe, since connection could drop during operation. With the `ConsistentRetry` policy, the steps will be embbeded into a transaction. Use it carefully because the delivery rate will be reduced by an order of magnitude.

```go
consumer := pubsub.NewConsumer(conn, "example-consumer-1", pubsub.ConsumerOptions{
    Queue: pubsub.ConsumerOptionsQueue{
        Name: "product.onEdit",
    },
    // ...
    RetryStrategy:    mo.Some(pubsub.NewExponentialRetryStrategy(3, 3*time.Second, 2)),
    RetryConsistency: mo.Some(pubsub.ConsistentRetry),
})
```

## Run examples

```sh
# run rabbitmq
docker-compose up rabbitmq
```

```sh
# run producer
cd examples/producer/
go mod download
go run main.go --rabbitmq-uri amqp://dev:dev@localhost:5672
```

```sh
# run consumer
cd examples/consumer/
go mod download
go run main.go --rabbitmq-uri amqp://dev:dev@localhost:5672
```

Then trigger network failure, by restarting rabbitmq:

```sh
docker-compose restart rabbitmq
```

## Todo

- Connection pooling (eg: 10 connections, 100 channels per connections)
- Better documentation
- Testing + CI
- BatchPublish + PublishWithConfirmation + BatchPublishWithConfirmation
