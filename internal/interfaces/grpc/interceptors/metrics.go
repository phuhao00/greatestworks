package interceptors

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"greatestworks/internal/infrastructure/logger"
)

// MetricsCollector 指标收集器接口
type MetricsCollector interface {
	IncRequestCount(method string, code codes.Code)
	ObserveRequestDuration(method string, duration time.Duration)
	IncActiveConnections()
	DecActiveConnections()
	IncStreamMessages(method string, direction string)
	SetGaugeValue(name string, value float64)
}

// DefaultMetricsCollector 默认指标收集器实现
type DefaultMetricsCollector struct {
	requestCounts    map[string]map[codes.Code]int64
	requestDurations map[string][]time.Duration
	activeConns      int64
	streamMessages   map[string]map[string]int64
	gauges           map[string]float64
	mutex            sync.RWMutex
	logger           logger.Logger
}

// NewDefaultMetricsCollector 创建默认指标收集器
func NewDefaultMetricsCollector(logger logger.Logger) *DefaultMetricsCollector {
	return &DefaultMetricsCollector{
		requestCounts:    make(map[string]map[codes.Code]int64),
		requestDurations: make(map[string][]time.Duration),
		activeConns:      0,
		streamMessages:   make(map[string]map[string]int64),
		gauges:           make(map[string]float64),
		logger:           logger,
	}
}

// IncRequestCount 增加请求计数
func (c *DefaultMetricsCollector) IncRequestCount(method string, code codes.Code) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.requestCounts[method] == nil {
		c.requestCounts[method] = make(map[codes.Code]int64)
	}
	c.requestCounts[method][code]++

	c.logger.Debug("Request count incremented", 
		"method", method, 
		"code", code.String(), 
		"count", c.requestCounts[method][code])
}

// ObserveRequestDuration 记录请求持续时间
func (c *DefaultMetricsCollector) ObserveRequestDuration(method string, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.requestDurations[method] == nil {
		c.requestDurations[method] = make([]time.Duration, 0)
	}
	c.requestDurations[method] = append(c.requestDurations[method], duration)

	// 保持最近1000个记录
	if len(c.requestDurations[method]) > 1000 {
		c.requestDurations[method] = c.requestDurations[method][1:]
	}

	c.logger.Debug("Request duration observed", 
		"method", method, 
		"duration", duration.String())
}

// IncActiveConnections 增加活跃连接数
func (c *DefaultMetricsCollector) IncActiveConnections() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.activeConns++
	c.logger.Debug("Active connections incremented", "count", c.activeConns)
}

// DecActiveConnections 减少活跃连接数
func (c *DefaultMetricsCollector) DecActiveConnections() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.activeConns > 0 {
		c.activeConns--
	}
	c.logger.Debug("Active connections decremented", "count", c.activeConns)
}

// IncStreamMessages 增加流消息计数
func (c *DefaultMetricsCollector) IncStreamMessages(method string, direction string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.streamMessages[method] == nil {
		c.streamMessages[method] = make(map[string]int64)
	}
	c.streamMessages[method][direction]++

	c.logger.Debug("Stream message count incremented", 
		"method", method, 
		"direction", direction, 
		"count", c.streamMessages[method][direction])
}

// SetGaugeValue 设置仪表值
func (c *DefaultMetricsCollector) SetGaugeValue(name string, value float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.gauges[name] = value
	c.logger.Debug("Gauge value set", "name", name, "value", value)
}

// GetMetrics 获取所有指标
func (c *DefaultMetricsCollector) GetMetrics() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	metrics := make(map[string]interface{})

	// 请求计数
	requestCounts := make(map[string]map[string]int64)
	for method, codes := range c.requestCounts {
		requestCounts[method] = make(map[string]int64)
		for code, count := range codes {
			requestCounts[method][code.String()] = count
		}
	}
	metrics["request_counts"] = requestCounts

	// 请求持续时间统计
	durationStats := make(map[string]map[string]interface{})
	for method, durations := range c.requestDurations {
		if len(durations) > 0 {
			stats := calculateDurationStats(durations)
			durationStats[method] = stats
		}
	}
	metrics["request_durations"] = durationStats

	// 活跃连接数
	metrics["active_connections"] = c.activeConns

	// 流消息计数
	metrics["stream_messages"] = c.streamMessages

	// 仪表值
	metrics["gauges"] = c.gauges

	return metrics
}

// calculateDurationStats 计算持续时间统计
func calculateDurationStats(durations []time.Duration) map[string]interface{} {
	if len(durations) == 0 {
		return nil
	}

	// 计算总和、最小值、最大值
	var total time.Duration
	min := durations[0]
	max := durations[0]

	for _, d := range durations {
		total += d
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}

	avg := total / time.Duration(len(durations))

	// 计算百分位数（简化实现）
	p50, p90, p95, p99 := calculatePercentiles(durations)

	return map[string]interface{}{
		"count":       len(durations),
		"total_ms":    total.Milliseconds(),
		"avg_ms":      avg.Milliseconds(),
		"min_ms":      min.Milliseconds(),
		"max_ms":      max.Milliseconds(),
		"p50_ms":      p50.Milliseconds(),
		"p90_ms":      p90.Milliseconds(),
		"p95_ms":      p95.Milliseconds(),
		"p99_ms":      p99.Milliseconds(),
	}
}

// calculatePercentiles 计算百分位数（简化实现）
func calculatePercentiles(durations []time.Duration) (p50, p90, p95, p99 time.Duration) {
	if len(durations) == 0 {
		return 0, 0, 0, 0
	}

	// 简化实现：使用索引计算百分位数
	len := len(durations)
	p50Index := len * 50 / 100
	p90Index := len * 90 / 100
	p95Index := len * 95 / 100
	p99Index := len * 99 / 100

	if p50Index >= len {
		p50Index = len - 1
	}
	if p90Index >= len {
		p90Index = len - 1
	}
	if p95Index >= len {
		p95Index = len - 1
	}
	if p99Index >= len {
		p99Index = len - 1
	}

	return durations[p50Index], durations[p90Index], durations[p95Index], durations[p99Index]
}

// 全局指标收集器
var globalMetricsCollector MetricsCollector
var metricsOnce sync.Once

// GetGlobalMetricsCollector 获取全局指标收集器
func GetGlobalMetricsCollector(logger logger.Logger) MetricsCollector {
	metricsOnce.Do(func() {
		globalMetricsCollector = NewDefaultMetricsCollector(logger)
	})
	return globalMetricsCollector
}

// MetricsUnaryInterceptor 指标拦截器（一元RPC）
func MetricsUnaryInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	collector := GetGlobalMetricsCollector(logger)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// 增加活跃连接数
		collector.IncActiveConnections()
		defer collector.DecActiveConnections()

		// 执行处理器
		resp, err := handler(ctx, req)

		// 计算持续时间
		duration := time.Since(start)

		// 获取状态码
		statusCode := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
			} else {
				statusCode = codes.Internal
			}
		}

		// 记录指标
		collector.IncRequestCount(info.FullMethod, statusCode)
		collector.ObserveRequestDuration(info.FullMethod, duration)

		// 记录慢请求指标
		if duration > 5*time.Second {
			collector.IncRequestCount(info.FullMethod+":slow", statusCode)
		}

		return resp, err
	}
}

// MetricsStreamInterceptor 指标拦截器（流式RPC）
func MetricsStreamInterceptor(logger logger.Logger) grpc.StreamServerInterceptor {
	collector := GetGlobalMetricsCollector(logger)

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// 增加活跃连接数
		collector.IncActiveConnections()
		defer collector.DecActiveConnections()

		// 创建包装的流用于统计消息
		wrappedStream := &metricsServerStream{
			ServerStream: stream,
			collector:    collector,
			methodName:   info.FullMethod,
		}

		// 执行处理器
		err := handler(srv, wrappedStream)

		// 计算持续时间
		duration := time.Since(start)

		// 获取状态码
		statusCode := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
			} else {
				statusCode = codes.Internal
			}
		}

		// 记录指标
		collector.IncRequestCount(info.FullMethod, statusCode)
		collector.ObserveRequestDuration(info.FullMethod, duration)

		// 记录长时间运行的流
		if duration > 30*time.Second {
			collector.IncRequestCount(info.FullMethod+":long_running", statusCode)
		}

		return err
	}
}

// metricsServerStream 包装的服务器流用于指标收集
type metricsServerStream struct {
	grpc.ServerStream
	collector  MetricsCollector
	methodName string
}

// SendMsg 发送消息
func (s *metricsServerStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.collector.IncStreamMessages(s.methodName, "sent")
	} else {
		s.collector.IncStreamMessages(s.methodName, "send_error")
	}
	return err
}

// RecvMsg 接收消息
func (s *metricsServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		s.collector.IncStreamMessages(s.methodName, "received")
	} else if err.Error() != "EOF" {
		s.collector.IncStreamMessages(s.methodName, "recv_error")
	}
	return err
}

// GetMetricsHandler 获取指标处理器（用于HTTP暴露指标）
func GetMetricsHandler(logger logger.Logger) func() map[string]interface{} {
	collector := GetGlobalMetricsCollector(logger)
	return func() map[string]interface{} {
		if defaultCollector, ok := collector.(*DefaultMetricsCollector); ok {
			return defaultCollector.GetMetrics()
		}
		return map[string]interface{}{
			"error": "metrics collector not available",
		}
	}
}

// RecordCustomMetric 记录自定义指标
func RecordCustomMetric(logger logger.Logger, name string, value float64) {
	collector := GetGlobalMetricsCollector(logger)
	collector.SetGaugeValue(name, value)
}

// RecordMethodCall 记录方法调用指标
func RecordMethodCall(logger logger.Logger, method string, duration time.Duration, success bool) {
	collector := GetGlobalMetricsCollector(logger)

	code := codes.OK
	if !success {
		code = codes.Internal
	}

	collector.IncRequestCount(method, code)
	collector.ObserveRequestDuration(method, duration)
}

// GetMethodStats 获取方法统计信息
func GetMethodStats(logger logger.Logger, method string) map[string]interface{} {
	collector := GetGlobalMetricsCollector(logger)
	if defaultCollector, ok := collector.(*DefaultMetricsCollector); ok {
		metrics := defaultCollector.GetMetrics()
		
		result := make(map[string]interface{})
		
		// 请求计数
		if requestCounts, ok := metrics["request_counts"].(map[string]map[string]int64); ok {
			if methodCounts, exists := requestCounts[method]; exists {
				result["request_counts"] = methodCounts
			}
		}
		
		// 请求持续时间
		if durations, ok := metrics["request_durations"].(map[string]map[string]interface{}); ok {
			if methodDurations, exists := durations[method]; exists {
				result["request_durations"] = methodDurations
			}
		}
		
		// 流消息
		if streamMessages, ok := metrics["stream_messages"].(map[string]map[string]int64); ok {
			if methodMessages, exists := streamMessages[method]; exists {
				result["stream_messages"] = methodMessages
			}
		}
		
		return result
	}
	
	return map[string]interface{}{
		"error": "method stats not available",
	}
}