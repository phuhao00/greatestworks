package handlers

import (
	"context"
	"fmt"
	"reflect"
)

// Query 查询接口
type Query interface {
	QueryType() string
	Validate() error
}

// QueryHandler 查询处理器接口
type QueryHandler[T Query, R any] interface {
	Handle(ctx context.Context, query T) (R, error)
}

// QueryBus 查询总线
type QueryBus struct {
	handlers map[string]interface{}
}

// NewQueryBus 创建查询总线
func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[string]interface{}),
	}
}

// RegisterHandler 注册查询处理器
func (bus *QueryBus) RegisterHandler(queryType string, handler interface{}) {
	bus.handlers[queryType] = handler
}

// Execute 执行查询
func (bus *QueryBus) Execute(ctx context.Context, query Query) (interface{}, error) {
	// 验证查询
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("query validation failed: %w", err)
	}
	
	// 获取处理器
	handler, exists := bus.handlers[query.QueryType()]
	if !exists {
		return nil, fmt.Errorf("no handler registered for query type: %s", query.QueryType())
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
		reflect.ValueOf(query),
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

// ExecuteTyped 执行类型化查询
func ExecuteQueryTyped[T Query, R any](ctx context.Context, bus *QueryBus, query T) (R, error) {
	result, err := bus.Execute(ctx, query)
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

// QueryMiddleware 查询中间件接口
type QueryMiddleware interface {
	Execute(ctx context.Context, query Query, next func(context.Context, Query) (interface{}, error)) (interface{}, error)
}

// CachingMiddleware 缓存中间件
type CachingMiddleware struct {
	cache Cache
}

// Cache 缓存接口
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl int) error
	Delete(key string) error
}

// NewCachingMiddleware 创建缓存中间件
func NewCachingMiddleware(cache Cache) *CachingMiddleware {
	return &CachingMiddleware{cache: cache}
}

// Execute 执行缓存中间件
func (m *CachingMiddleware) Execute(ctx context.Context, query Query, next func(context.Context, Query) (interface{}, error)) (interface{}, error) {
	// 生成缓存键
	cacheKey := fmt.Sprintf("%s:%v", query.QueryType(), query)
	
	// 尝试从缓存获取
	if cached, found := m.cache.Get(cacheKey); found {
		return cached, nil
	}
	
	// 执行查询
	result, err := next(ctx, query)
	if err != nil {
		return nil, err
	}
	
	// 缓存结果（TTL 5分钟）
	m.cache.Set(cacheKey, result, 300)
	
	return result, nil
}

// QueryLoggingMiddleware 查询日志中间件
type QueryLoggingMiddleware struct {
	logger Logger
}

// NewQueryLoggingMiddleware 创建查询日志中间件
func NewQueryLoggingMiddleware(logger Logger) *QueryLoggingMiddleware {
	return &QueryLoggingMiddleware{logger: logger}
}

// Execute 执行查询日志中间件
func (m *QueryLoggingMiddleware) Execute(ctx context.Context, query Query, next func(context.Context, Query) (interface{}, error)) (interface{}, error) {
	m.logger.Info("executing query", "type", query.QueryType())
	
	result, err := next(ctx, query)
	if err != nil {
		m.logger.Error("query execution failed", err, "type", query.QueryType())
		return nil, err
	}
	
	m.logger.Info("query executed successfully", "type", query