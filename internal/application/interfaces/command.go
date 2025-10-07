// Package interfaces 定义应用层接口
package interfaces

import "context"

// Command 命令接口
type Command interface {
	CommandType() string
	Validate() error
}

// CommandHandler 命令处理器接口
type CommandHandler[T Command, R any] interface {
	Handle(ctx context.Context, cmd T) (R, error)
}

// CommandBus 命令总线接口
type CommandBus interface {
	RegisterHandler(commandType string, handler interface{})
	Execute(ctx context.Context, cmd Command) (interface{}, error)
}

// Query 查询接口
type Query interface {
	QueryType() string
	Validate() error
}

// QueryHandler 查询处理器接口
type QueryHandler[T Query, R any] interface {
	Handle(ctx context.Context, query T) (R, error)
}

// QueryBus 查询总线接口
type QueryBus interface {
	RegisterHandler(queryType string, handler interface{})
	Execute(ctx context.Context, query Query) (interface{}, error)
}

// Event 事件接口
type Event interface {
	EventType() string
	AggregateID() string
	Version() int64
	OccurredAt() int64
}

// EventHandler 事件处理器接口
type EventHandler[T Event] interface {
	Handle(ctx context.Context, event T) error
}

// EventBus 事件总线接口
type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string, handler EventHandler[Event]) error
	Unsubscribe(eventType string, handler EventHandler[Event]) error
}

// Middleware 中间件接口
type Middleware interface {
	Execute(ctx context.Context, cmd Command, next func(context.Context, Command) (interface{}, error)) (interface{}, error)
}

// QueryMiddleware 查询中间件接口
type QueryMiddleware interface {
	Execute(ctx context.Context, query Query, next func(context.Context, Query) (interface{}, error)) (interface{}, error)
}
