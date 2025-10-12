package handlers

import (
	"context"
	"fmt"
	"reflect"
)

// Command 命令接口
type Command interface {
	CommandType() string
	Validate() error
}

// CommandHandler 命令处理器接口
type CommandHandler[T Command, R any] interface {
	Handle(ctx context.Context, cmd T) (R, error)
}

// CommandBus 命令总线
type CommandBus struct {
	handlers map[string]interface{}
}

// NewCommandBus 创建命令总线
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[string]interface{}),
	}
}

// RegisterHandler 注册命令处理器
func (bus *CommandBus) RegisterHandler(commandType string, handler interface{}) {
	bus.handlers[commandType] = handler
}

// Execute 执行命令
func (bus *CommandBus) Execute(ctx context.Context, cmd Command) (interface{}, error) {
	// 验证命令
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("command validation failed: %w", err)
	}

	// 获取处理器
	handler, exists := bus.handlers[cmd.CommandType()]
	if !exists {
		return nil, fmt.Errorf("no handler registered for command type: %s", cmd.CommandType())
	}

	// 使用反射调用处理器
	handlerValue := reflect.ValueOf(handler)
	handlerType := reflect.TypeOf(handler)

	// 查找Handle方法
	handleMethod, exists := handlerType.MethodByName("Handle")
	if !exists {
		return nil, fmt.Errorf("handler does not have Handle method")
	}

	// 调用Handle方法
	args := []reflect.Value{
		handlerValue,
		reflect.ValueOf(ctx),
		reflect.ValueOf(cmd),
	}

	results := handleMethod.Func.Call(args)
	if len(results) != 2 {
		return nil, fmt.Errorf("handler Handle method should return (result, error)")
	}

	// 检查错误
	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}

	return results[0].Interface(), nil
}

// ExecuteTyped 执行类型化命令
func ExecuteTyped[T Command, R any](ctx context.Context, bus *CommandBus, cmd T) (R, error) {
	result, err := bus.Execute(ctx, cmd)
	if err != nil {
		var zero R
		return zero, err
	}

	if typedResult, ok := result.(R); ok {
		return typedResult, nil
	}

	var zero R
	return zero, fmt.Errorf("unexpected result type")
}

// Middleware 中间件接口
type Middleware interface {
	Execute(ctx context.Context, cmd Command, next func(context.Context, Command) (interface{}, error)) (interface{}, error)
}

// ValidationMiddleware 验证中间件
type ValidationMiddleware struct{}

// Execute 执行验证中间件
func (m *ValidationMiddleware) Execute(ctx context.Context, cmd Command, next func(context.Context, Command) (interface{}, error)) (interface{}, error) {
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	return next(ctx, cmd)
}

// LoggingMiddleware 日志中间件
type LoggingMiddleware struct {
	logger Logger
}

// Logger 日志接口
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
}

// NewLoggingMiddleware 创建日志中间件
func NewLoggingMiddleware(logger Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

// Execute 执行日志中间件
func (m *LoggingMiddleware) Execute(ctx context.Context, cmd Command, next func(context.Context, Command) (interface{}, error)) (interface{}, error) {
	m.logger.Info("executing command", "type", cmd.CommandType())

	result, err := next(ctx, cmd)
	if err != nil {
		m.logger.Error("command execution failed", err, "type", cmd.CommandType())
		return nil, err
	}

	m.logger.Info("command executed successfully", "type", cmd.CommandType())
	return result, nil
}
