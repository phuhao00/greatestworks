package pubsub

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

type ProducerOptionsExchange struct {
	Name string
	Kind ExchangeKind

	// optional arguments
	Durable    mo.Option[bool]       // default true
	AutoDelete mo.Option[bool]       // default false
	Internal   mo.Option[bool]       // default false
	NoWait     mo.Option[bool]       // default false
	Args       mo.Option[amqp.Table] // default nil
}

type ProducerOptions struct {
	Exchange ProducerOptionsExchange
}

type Producer struct {
	conn    *Connection
	name    string
	options ProducerOptions

	channel *amqp.Channel
	done    *rpc[struct{}, struct{}]
}

func NewProducer(conn *Connection, name string, opt ProducerOptions) *Producer {
	doneCh := make(chan struct{})

	p := &Producer{
		conn:    conn,
		name:    name,
		options: opt,

		channel: nil,
		done:    newRPC[struct{}, struct{}](doneCh),
	}

	go p.lifecycle()

	return p
}

func (p *Producer) lifecycle() {
	cancel, ch := p.conn.ListenConnection()
	onConnect := make(chan *amqp.Connection, 42)
	onDisconnect := make(chan struct{}, 42)

	for {
		select {
		case conn := <-ch:
			if conn != nil {
				onConnect <- conn
			} else {
				onDisconnect <- struct{}{}
			}

		case conn := <-onConnect:
			err := p.setupProducer(conn)
			if err != nil {
				logger("AMQP producer '%s': %s", p.name, err.Error())
				onDisconnect <- struct{}{}
			}

		case <-onDisconnect:
			if p.channel != nil {
				lo.Try0(func() { p.channel.Close() })

				p.channel = nil
			}

		case req := <-p.done.C:
			cancel()
			req.B(struct{}{})
			return
		}
	}
}

func (p *Producer) Close() error {
	_ = p.done.Send(struct{}{})
	safeClose(p.done.C)
	return nil
}

func (p *Producer) setupProducer(conn *amqp.Connection) error {
	// create a channel dedicated to this producer
	channel, err := conn.Channel()
	if err != nil {
		return err
	}

	// create exchange if not exist
	err = channel.ExchangeDeclare(
		p.options.Exchange.Name,
		string(p.options.Exchange.Kind),
		p.options.Exchange.Durable.OrElse(true),
		p.options.Exchange.AutoDelete.OrElse(false),
		p.options.Exchange.Internal.OrElse(false),
		p.options.Exchange.NoWait.OrElse(false),
		p.options.Exchange.Args.OrElse(nil),
	)
	if err != nil {
		_ = channel.Close()
		return err
	}

	p.channel = channel

	go p.handleCancel(conn, channel)

	return nil
}

func (p *Producer) handleCancel(conn *amqp.Connection, channel *amqp.Channel) {
	onClose := channel.NotifyClose(make(chan *amqp.Error))
	onCancel := channel.NotifyCancel(make(chan string))

	select {
	case err := <-onClose:
		if err != nil {
			logger("AMQP channel '%s': %s", p.name, err.Error())
		}
	case msg := <-onCancel:
		logger("AMQP channel '%s': %v", p.name, msg)

		lo.Try0(func() { channel.Close() })

		err := p.setupProducer(conn)
		if err != nil {
			logger("AMQP producer '%s': %s", p.name, err.Error())
		}
	}
}

/**
 * API
 */

func (p *Producer) PublishWithContext(ctx context.Context, routingKey string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	if p.channel == nil {
		return fmt.Errorf("AMQP: channel '%s' not available", p.name)
	}

	return p.channel.PublishWithContext(
		ctx,
		p.options.Exchange.Name,
		routingKey,
		mandatory,
		immediate,
		msg,
	)
}

func (p *Producer) PublishWithDeferredConfirmWithContext(ctx context.Context, routingKey string, mandatory bool, immediate bool, msg amqp.Publishing) (*amqp.DeferredConfirmation, error) {
	if p.channel == nil {
		return nil, fmt.Errorf("AMQP: channel '%s' not available", p.name)
	}

	return p.channel.PublishWithDeferredConfirmWithContext(
		ctx,
		p.options.Exchange.Name,
		routingKey,
		mandatory,
		immediate,
		msg,
	)
}

func (p *Producer) Publish(routingKey string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	return p.PublishWithContext(
		context.Background(),
		routingKey,
		mandatory,
		immediate,
		msg,
	)
}

func (p *Producer) PublishWithDeferredConfirm(routingKey string, mandatory bool, immediate bool, msg amqp.Publishing) (*amqp.DeferredConfirmation, error) {
	return p.PublishWithDeferredConfirmWithContext(
		context.Background(),
		routingKey,
		mandatory,
		immediate,
		msg,
	)
}
