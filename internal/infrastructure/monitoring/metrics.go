package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"greatestworks/internal/infrastructure/logger"
)

// CounterInterface 计数器指标接口
type CounterInterface interface {
	Inc()
	Add(value int64)
	GetName() string
	GetType() MetricType
	GetValue() interface{}
	GetLabels() map[string]string
	GetTimestamp() time.Time
	Help() string
	Type() MetricType
	Name() string
	Get() int64
}

// GaugeInterface 仪表盘指标接口
type GaugeInterface interface {
	Set(value float64)
	Inc()
	Dec()
	Add(value float64)
	Sub(value float64)
	GetName() string
	GetType() MetricType
	GetValue() interface{}
	GetLabels() map[string]string
	GetTimestamp() time.Time
	Help() string
	Type() MetricType
	Name() string
}

// HistogramInterface 直方图指标接口
type HistogramInterface interface {
	Observe(value float64)
	GetName() string
	GetType() MetricType
	GetValue() interface{}
	GetLabels() map[string]string
	GetTimestamp() time.Time
	Help() string
	Type() MetricType
	Name() string
}

// Summary 摘要指标接口
type Summary interface {
	Observe(value float64)
	GetName() string
	GetType() MetricType
	GetValue() interface{}
	GetLabels() map[string]string
	GetTimestamp() time.Time
}

// Timer 计时器接口
type Timer interface {
	Start() TimerContext
	Time(fn func())
	TimeContext(ctx context.Context, fn func(context.Context))
	ObserveDuration(duration time.Duration)
}

// TimerContext 计时器上下文接口
type TimerContext interface {
	Stop()
	Duration() time.Duration
}

// Config 监控配置
type Config struct {
	Enabled              bool              `json:"enabled"`
	Port                 int               `json:"port"`
	Path                 string            `json:"path"`
	Namespace            string            `json:"namespace"`
	Subsystem            string            `json:"subsystem"`
	Host                 string            `json:"host"`
	Labels               map[string]string `json:"labels"`
	EnableRuntimeMetrics bool              `json:"enable_runtime_metrics"`
	EnableProcessMetrics bool              `json:"enable_process_metrics"`
}

// Factory 指标工厂接口
type Factory interface {
	NewCounter(name, help string, labels Labels) CounterInterface
	NewGauge(name, help string, labels Labels) GaugeInterface
	NewHistogram(name, help string, buckets []float64, labels Labels) HistogramInterface
	NewSummary(name, help string, objectives map[float64]float64, labels Labels) Summary
}

// Manager 监控管理器接口
type Manager interface {
	GetRegistry() Registry
	GetFactory() Factory
	RegisterCollector(collector Collector) error
}

// Registry 指标注册表接口
type Registry interface {
	Register(collector Collector) error
	Unregister(collector Collector) bool
	MustRegister(collectors ...Collector)
	Gather() ([]*MetricFamily, error)
}

// Collector 收集器接口
type Collector interface {
	Describe(chan<- *MetricDesc)
	Collect(chan<- Metric)
}

// PrometheusCollector Prometheus收集器接口
type PrometheusCollector interface {
	Describe(chan<- *prometheus.Desc)
	Collect(chan<- prometheus.Metric)
}

// Labels 标签类型
type Labels map[string]string

// MetricDesc 指标描述
type MetricDesc struct {
	Name string
	Help string
	Type MetricType
}

// MetricFamily 指标族
type MetricFamily struct {
	Name    string
	Help    string
	Type    MetricType
	Metrics []*Sample
}

// Sample 样本
type Sample struct {
	Labels    Labels
	Value     float64
	Timestamp time.Time
	Buckets   []Bucket
	Quantiles []Quantile
}

// Bucket 直方图桶
type Bucket struct {
	UpperBound float64
	Count      uint64
}

// Quantile 分位数
type Quantile struct {
	Quantile float64
	Value    float64
}

// 常量定义
const (
	CounterType   MetricType = "counter"
	GaugeType     MetricType = "gauge"
	HistogramType MetricType = "histogram"
	SummaryType   MetricType = "summary"
)

// 错误定义
var (
	ErrMetricExists   = fmt.Errorf("metric already exists")
	ErrMetricNotFound = fmt.Errorf("metric not found")
)

// MetricType 指标类型
type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
	MetricTypeSummary   MetricType = "summary"
)

// Metric 指标接口
type Metric interface {
	GetName() string
	GetType() MetricType
	GetValue() interface{}
	GetLabels() map[string]string
	GetTimestamp() time.Time
	Help() string
	Type() MetricType
	Name() string
}

// Counter 计数器指标
type Counter struct {
	name      string
	help      string
	value     int64
	labels    map[string]string
	timestamp time.Time
	mutex     sync.RWMutex
}

// NewCounter 创建计数器
func NewCounter(name string, labels map[string]string) *Counter {
	return &Counter{
		name:      name,
		help:      "",
		value:     0,
		labels:    labels,
		timestamp: time.Now(),
	}
}

// NewCounterWithHelp 创建带帮助信息的计数器
func NewCounterWithHelp(name, help string, labels map[string]string) *Counter {
	return &Counter{
		name:      name,
		help:      help,
		value:     0,
		labels:    labels,
		timestamp: time.Now(),
	}
}

// Inc 增加计数
func (c *Counter) Inc() {
	c.Add(1)
}

// Add 增加指定值
func (c *Counter) Add(value int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value += value
	c.timestamp = time.Now()
}

// GetName 获取名称
func (c *Counter) GetName() string {
	return c.name
}

// GetType 获取类型
func (c *Counter) GetType() MetricType {
	return MetricTypeCounter
}

// GetValue 获取值
func (c *Counter) GetValue() interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.value
}

// GetLabels 获取标签
func (c *Counter) GetLabels() map[string]string {
	return c.labels
}

// GetTimestamp 获取时间戳
func (c *Counter) GetTimestamp() time.Time {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.timestamp
}

// Help 获取帮助信息
func (c *Counter) Help() string {
	return c.help
}

// Type 获取类型
func (c *Counter) Type() MetricType {
	return MetricTypeCounter
}

// Name 获取名称
func (c *Counter) Name() string {
	return c.name
}

// Get 获取当前值
func (c *Counter) Get() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.value
}

// Gauge 仪表指标
type Gauge struct {
	name      string
	help      string
	value     float64
	labels    map[string]string
	timestamp time.Time
	mutex     sync.RWMutex
}

// NewGauge 创建仪表
func NewGauge(name string, labels map[string]string) *Gauge {
	return &Gauge{
		name:      name,
		help:      "",
		value:     0,
		labels:    labels,
		timestamp: time.Now(),
	}
}

// NewGaugeWithHelp 创建带帮助信息的仪表
func NewGaugeWithHelp(name, help string, labels map[string]string) *Gauge {
	return &Gauge{
		name:      name,
		help:      help,
		value:     0,
		labels:    labels,
		timestamp: time.Now(),
	}
}

// Set 设置值
func (g *Gauge) Set(value float64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.value = value
	g.timestamp = time.Now()
}

// Inc 增加1
func (g *Gauge) Inc() {
	g.Add(1)
}

// Dec 减少1
func (g *Gauge) Dec() {
	g.Add(-1)
}

// Add 增加指定值
func (g *Gauge) Add(value float64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.value += value
	g.timestamp = time.Now()
}

// Sub 减少指定值
func (g *Gauge) Sub(value float64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.value -= value
	g.timestamp = time.Now()
}

// GetName 获取名称
func (g *Gauge) GetName() string {
	return g.name
}

// GetType 获取类型
func (g *Gauge) GetType() MetricType {
	return MetricTypeGauge
}

// GetValue 获取值
func (g *Gauge) GetValue() interface{} {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.value
}

// GetLabels 获取标签
func (g *Gauge) GetLabels() map[string]string {
	return g.labels
}

// GetTimestamp 获取时间戳
func (g *Gauge) GetTimestamp() time.Time {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.timestamp
}

// Help 获取帮助信息
func (g *Gauge) Help() string {
	return g.help
}

// Type 获取类型
func (g *Gauge) Type() MetricType {
	return MetricTypeGauge
}

// Name 获取名称
func (g *Gauge) Name() string {
	return g.name
}

// Histogram 直方图指标
type Histogram struct {
	name      string
	help      string
	buckets   []float64
	counts    []int64
	sum       float64
	total     int64
	labels    map[string]string
	timestamp time.Time
	mutex     sync.RWMutex
}

// NewHistogram 创建直方图
func NewHistogram(name string, buckets []float64, labels map[string]string) *Histogram {
	if buckets == nil {
		buckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
	}

	return &Histogram{
		name:      name,
		help:      "",
		buckets:   buckets,
		counts:    make([]int64, len(buckets)+1), // +1 for +Inf bucket
		sum:       0,
		total:     0,
		labels:    labels,
		timestamp: time.Now(),
	}
}

// NewHistogramWithHelp 创建带帮助信息的直方图
func NewHistogramWithHelp(name, help string, buckets []float64, labels map[string]string) *Histogram {
	if buckets == nil {
		buckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
	}

	return &Histogram{
		name:      name,
		help:      help,
		buckets:   buckets,
		counts:    make([]int64, len(buckets)+1), // +1 for +Inf bucket
		sum:       0,
		total:     0,
		labels:    labels,
		timestamp: time.Now(),
	}
}

// Observe 观察值
func (h *Histogram) Observe(value float64) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.sum += value
	h.total++
	h.timestamp = time.Now()

	// 找到对应的桶
	for i, bucket := range h.buckets {
		if value <= bucket {
			h.counts[i]++
			return
		}
	}
	// +Inf bucket
	h.counts[len(h.buckets)]++
}

// GetName 获取名称
func (h *Histogram) GetName() string {
	return h.name
}

// GetType 获取类型
func (h *Histogram) GetType() MetricType {
	return MetricTypeHistogram
}

// GetValue 获取值
func (h *Histogram) GetValue() interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return map[string]interface{}{
		"buckets": h.buckets,
		"counts":  h.counts,
		"sum":     h.sum,
		"count":   h.total,
	}
}

// GetLabels 获取标签
func (h *Histogram) GetLabels() map[string]string {
	return h.labels
}

// GetTimestamp 获取时间戳
func (h *Histogram) GetTimestamp() time.Time {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.timestamp
}

// Help 获取帮助信息
func (h *Histogram) Help() string {
	return h.help
}

// Type 获取类型
func (h *Histogram) Type() MetricType {
	return MetricTypeHistogram
}

// Name 获取名称
func (h *Histogram) Name() string {
	return h.name
}

// MetricsRegistry 指标注册表
type MetricsRegistry struct {
	metrics map[string]Metric
	mutex   sync.RWMutex
	logger  logger.Logger
}

// NewMetricsRegistry 创建指标注册表
func NewMetricsRegistry(logger logger.Logger) *MetricsRegistry {
	return &MetricsRegistry{
		metrics: make(map[string]Metric),
		logger:  logger,
	}
}

// Register 注册指标
func (r *MetricsRegistry) Register(metric Metric) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	name := metric.GetName()
	if _, exists := r.metrics[name]; exists {
		return fmt.Errorf("metric %s already registered", name)
	}

	r.metrics[name] = metric
	r.logger.Debug("Metric registered", "name", name, "type", metric.GetType())
	return nil
}

// Unregister 注销指标
func (r *MetricsRegistry) Unregister(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.metrics, name)
	r.logger.Debug("Metric unregistered", "name", name)
}

// GetMetric 获取指标
func (r *MetricsRegistry) GetMetric(name string) (Metric, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	metric, exists := r.metrics[name]
	return metric, exists
}

// GetAllMetrics 获取所有指标
func (r *MetricsRegistry) GetAllMetrics() map[string]Metric {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(map[string]Metric)
	for name, metric := range r.metrics {
		result[name] = metric
	}
	return result
}

// GetMetricsData 获取指标数据
func (r *MetricsRegistry) GetMetricsData() map[string]interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	data := make(map[string]interface{})
	for name, metric := range r.metrics {
		data[name] = map[string]interface{}{
			"name":      metric.GetName(),
			"type":      metric.GetType(),
			"value":     metric.GetValue(),
			"labels":    metric.GetLabels(),
			"timestamp": metric.GetTimestamp().Unix(),
		}
	}
	return data
}

// MonitoringService 监控服务
type MonitoringService struct {
	registry *MetricsRegistry
	logger   logger.Logger

	// 预定义指标
	httpRequestsTotal    *Counter
	httpRequestDuration  *Histogram
	httpActiveConnections *Gauge
	tcpConnectionsTotal   *Counter
	tcpActiveConnections  *Gauge
	systemMemoryUsage     *Gauge
	systemCPUUsage        *Gauge
	databaseConnections   *Gauge
	errorCount            *Counter
}

// NewMonitoringService 创建监控服务
func NewMonitoringService(logger logger.Logger) *MonitoringService {
	registry := NewMetricsRegistry(logger)

	service := &MonitoringService{
		registry: registry,
		logger:   logger,
	}

	// 初始化预定义指标
	service.initPredefinedMetrics()

	return service
}

// initPredefinedMetrics 初始化预定义指标
func (s *MonitoringService) initPredefinedMetrics() {
	// HTTP指标
	s.httpRequestsTotal = NewCounter("http_requests_total", map[string]string{"service": "greatestworks"})
	s.httpRequestDuration = NewHistogram("http_request_duration_seconds", nil, map[string]string{"service": "greatestworks"})
	s.httpActiveConnections = NewGauge("http_active_connections", map[string]string{"service": "greatestworks"})

	// TCP指标
	s.tcpConnectionsTotal = NewCounter("tcp_connections_total", map[string]string{"service": "greatestworks"})
	s.tcpActiveConnections = NewGauge("tcp_active_connections", map[string]string{"service": "greatestworks"})



	// 系统指标
	s.systemMemoryUsage = NewGauge("system_memory_usage_bytes", map[string]string{"service": "greatestworks"})
	s.systemCPUUsage = NewGauge("system_cpu_usage_percent", map[string]string{"service": "greatestworks"})
	s.databaseConnections = NewGauge("database_connections", map[string]string{"service": "greatestworks"})
	s.errorCount = NewCounter("errors_total", map[string]string{"service": "greatestworks"})

	// 注册指标
	s.registry.Register(s.httpRequestsTotal)
	s.registry.Register(s.httpRequestDuration)
	s.registry.Register(s.httpActiveConnections)
	s.registry.Register(s.tcpConnectionsTotal)
	s.registry.Register(s.tcpActiveConnections)

	s.registry.Register(s.systemMemoryUsage)
	s.registry.Register(s.systemCPUUsage)
	s.registry.Register(s.databaseConnections)
	s.registry.Register(s.errorCount)

	s.logger.Info("Predefined metrics initialized")
}

// RecordHTTPRequest 记录HTTP请求
func (s *MonitoringService) RecordHTTPRequest(duration time.Duration) {
	s.httpRequestsTotal.Inc()
	s.httpRequestDuration.Observe(duration.Seconds())
}

// RecordTCPConnection 记录TCP连接
func (s *MonitoringService) RecordTCPConnection() {
	s.tcpConnectionsTotal.Inc()
	s.tcpActiveConnections.Inc()
}

// RecordTCPDisconnection 记录TCP断开连接
func (s *MonitoringService) RecordTCPDisconnection() {
	s.tcpActiveConnections.Dec()
}



// RecordError 记录错误
func (s *MonitoringService) RecordError() {
	s.errorCount.Inc()
}

// UpdateSystemMetrics 更新系统指标
func (s *MonitoringService) UpdateSystemMetrics(memoryUsage, cpuUsage float64, dbConnections int) {
	s.systemMemoryUsage.Set(memoryUsage)
	s.systemCPUUsage.Set(cpuUsage)
	s.databaseConnections.Set(float64(dbConnections))
}

// GetRegistry 获取指标注册表
func (s *MonitoringService) GetRegistry() *MetricsRegistry {
	return s.registry
}

// GetMetricsHandler 获取指标HTTP处理器
func (s *MonitoringService) GetMetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := s.registry.GetMetricsData()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    data,
			"timestamp": time.Now().Unix(),
		})
	}
}

// HTTPMetricsMiddleware HTTP指标中间件
func (s *MonitoringService) HTTPMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 增加活跃连接数
		s.httpActiveConnections.Inc()
		defer s.httpActiveConnections.Dec()

		// 处理请求
		c.Next()

		// 记录指标
		duration := time.Since(start)
		s.RecordHTTPRequest(duration)

		// 如果有错误，记录错误指标
		if len(c.Errors) > 0 {
			s.RecordError()
		}
	}
}

// StartSystemMetricsCollection 启动系统指标收集
func (s *MonitoringService) StartSystemMetricsCollection(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("System metrics collection stopped")
			return
		case <-ticker.C:
			s.collectSystemMetrics()
		}
	}
}

// collectSystemMetrics 收集系统指标
func (s *MonitoringService) collectSystemMetrics() {
	// TODO: 实现实际的系统指标收集
	// 这里应该使用系统调用或第三方库来获取实际的系统指标
	// 例如：内存使用量、CPU使用率、磁盘使用量等

	// 示例数据
	memoryUsage := float64(1024 * 1024 * 512) // 512MB
	cpuUsage := float64(25.5)                 // 25.5%
	dbConnections := 10

	s.UpdateSystemMetrics(memoryUsage, cpuUsage, dbConnections)
	s.logger.Debug("System metrics collected", 
		"memory_usage", memoryUsage,
		"cpu_usage", cpuUsage,
		"db_connections", dbConnections)
}

// GetStats 获取监控统计信息
func (s *MonitoringService) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"metrics_count":      len(s.registry.GetAllMetrics()),
		"http_requests":      s.httpRequestsTotal.GetValue(),
		"tcp_connections":    s.tcpConnectionsTotal.GetValue(),

		"error_count":        s.errorCount.GetValue(),
		"active_http_conns":  s.httpActiveConnections.GetValue(),
		"active_tcp_conns":   s.tcpActiveConnections.GetValue(),
		"memory_usage":       s.systemMemoryUsage.GetValue(),
		"cpu_usage":          s.systemCPUUsage.GetValue(),
		"db_connections":     s.databaseConnections.GetValue(),
	}
}