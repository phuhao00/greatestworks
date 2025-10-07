package monitoring

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// CounterInterface 计数器指标接口
type CounterInterface interface {
	Inc()
	Add(value int64)
	GetName() string
}

// GaugeInterface 仪表指标接口
type GaugeInterface interface {
	Set(value float64)
	Inc()
	Dec()
	GetName() string
}

// HistogramInterface 直方图指标接口
type HistogramInterface interface {
	Observe(value float64)
	GetName() string
}

// SummaryInterface 摘要指标接口
type SummaryInterface interface {
	Observe(value float64)
	GetName() string
}

// MetricsCollector 指标收集器
type MetricsCollector struct {
	counters   map[string]CounterInterface
	gauges     map[string]GaugeInterface
	histograms map[string]HistogramInterface
	summaries  map[string]SummaryInterface
	mutex      sync.RWMutex
}

// NewMetricsCollector 创建指标收集器
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		counters:   make(map[string]CounterInterface),
		gauges:     make(map[string]GaugeInterface),
		histograms: make(map[string]HistogramInterface),
		summaries:  make(map[string]SummaryInterface),
	}
}

// RegisterCounter 注册计数器
func (mc *MetricsCollector) RegisterCounter(name string, counter CounterInterface) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.counters[name] = counter
}

// RegisterGauge 注册仪表
func (mc *MetricsCollector) RegisterGauge(name string, gauge GaugeInterface) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.gauges[name] = gauge
}

// RegisterHistogram 注册直方图
func (mc *MetricsCollector) RegisterHistogram(name string, histogram HistogramInterface) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.histograms[name] = histogram
}

// RegisterSummary 注册摘要
func (mc *MetricsCollector) RegisterSummary(name string, summary SummaryInterface) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.summaries[name] = summary
}

// GetCounter 获取计数器
func (mc *MetricsCollector) GetCounter(name string) (CounterInterface, bool) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	counter, exists := mc.counters[name]
	return counter, exists
}

// GetGauge 获取仪表
func (mc *MetricsCollector) GetGauge(name string) (GaugeInterface, bool) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	gauge, exists := mc.gauges[name]
	return gauge, exists
}

// GetHistogram 获取直方图
func (mc *MetricsCollector) GetHistogram(name string) (HistogramInterface, bool) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	histogram, exists := mc.histograms[name]
	return histogram, exists
}

// GetSummary 获取摘要
func (mc *MetricsCollector) GetSummary(name string) (SummaryInterface, bool) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	summary, exists := mc.summaries[name]
	return summary, exists
}

// PrometheusCounter Prometheus计数器实现
type PrometheusCounter struct {
	counter prometheus.Counter
	name    string
}

// NewPrometheusCounter 创建Prometheus计数器
func NewPrometheusCounter(name, help string) *PrometheusCounter {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})

	prometheus.MustRegister(counter)

	return &PrometheusCounter{
		counter: counter,
		name:    name,
	}
}

// Inc 增加计数
func (pc *PrometheusCounter) Inc() {
	pc.counter.Inc()
}

// Add 增加指定值
func (pc *PrometheusCounter) Add(value int64) {
	pc.counter.Add(float64(value))
}

// GetName 获取名称
func (pc *PrometheusCounter) GetName() string {
	return pc.name
}

// PrometheusGauge Prometheus仪表实现
type PrometheusGauge struct {
	gauge prometheus.Gauge
	name  string
}

// NewPrometheusGauge 创建Prometheus仪表
func NewPrometheusGauge(name, help string) *PrometheusGauge {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})

	prometheus.MustRegister(gauge)

	return &PrometheusGauge{
		gauge: gauge,
		name:  name,
	}
}

// Set 设置值
func (pg *PrometheusGauge) Set(value float64) {
	pg.gauge.Set(value)
}

// Inc 增加
func (pg *PrometheusGauge) Inc() {
	pg.gauge.Inc()
}

// Dec 减少
func (pg *PrometheusGauge) Dec() {
	pg.gauge.Dec()
}

// GetName 获取名称
func (pg *PrometheusGauge) GetName() string {
	return pg.name
}

// PrometheusHistogram Prometheus直方图实现
type PrometheusHistogram struct {
	histogram prometheus.Histogram
	name      string
}

// NewPrometheusHistogram 创建Prometheus直方图
func NewPrometheusHistogram(name, help string, buckets []float64) *PrometheusHistogram {
	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    name,
		Help:    help,
		Buckets: buckets,
	})

	prometheus.MustRegister(histogram)

	return &PrometheusHistogram{
		histogram: histogram,
		name:      name,
	}
}

// Observe 观察值
func (ph *PrometheusHistogram) Observe(value float64) {
	ph.histogram.Observe(value)
}

// GetName 获取名称
func (ph *PrometheusHistogram) GetName() string {
	return ph.name
}

// PrometheusSummary Prometheus摘要实现
type PrometheusSummary struct {
	summary prometheus.Summary
	name    string
}

// NewPrometheusSummary 创建Prometheus摘要
func NewPrometheusSummary(name, help string) *PrometheusSummary {
	summary := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: name,
		Help: help,
	})

	prometheus.MustRegister(summary)

	return &PrometheusSummary{
		summary: summary,
		name:    name,
	}
}

// Observe 观察值
func (ps *PrometheusSummary) Observe(value float64) {
	ps.summary.Observe(value)
}

// GetName 获取名称
func (ps *PrometheusSummary) GetName() string {
	return ps.name
}

// MetricsHandler 指标处理器
type MetricsHandler struct {
	collector *MetricsCollector
}

// NewMetricsHandler 创建指标处理器
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		collector: NewMetricsCollector(),
	}
}

// RegisterDefaultMetrics 注册默认指标
func (mh *MetricsHandler) RegisterDefaultMetrics() {
	// 注册请求计数器
	requestCounter := NewPrometheusCounter("http_requests_total", "Total number of HTTP requests")
	mh.collector.RegisterCounter("http_requests_total", requestCounter)

	// 注册响应时间直方图
	responseTimeHistogram := NewPrometheusHistogram("http_request_duration_seconds", "HTTP request duration in seconds",
		[]float64{0.1, 0.5, 1.0, 2.5, 5.0, 10.0})
	mh.collector.RegisterHistogram("http_request_duration_seconds", responseTimeHistogram)

	// 注册活跃连接数仪表
	activeConnectionsGauge := NewPrometheusGauge("active_connections", "Number of active connections")
	mh.collector.RegisterGauge("active_connections", activeConnectionsGauge)

	// 注册内存使用量仪表
	memoryUsageGauge := NewPrometheusGauge("memory_usage_bytes", "Memory usage in bytes")
	mh.collector.RegisterGauge("memory_usage_bytes", memoryUsageGauge)
}

// HandleMetrics 处理指标请求
func (mh *MetricsHandler) HandleMetrics(c *gin.Context) {
	// 设置内容类型
	c.Header("Content-Type", "text/plain")

	// 获取Prometheus指标
	registry := prometheus.NewRegistry()

	// 收集所有指标
	mh.collector.mutex.RLock()
	for _, counter := range mh.collector.counters {
		if pc, ok := counter.(*PrometheusCounter); ok {
			registry.MustRegister(pc.counter)
		}
	}
	for _, gauge := range mh.collector.gauges {
		if pg, ok := gauge.(*PrometheusGauge); ok {
			registry.MustRegister(pg.gauge)
		}
	}
	for _, histogram := range mh.collector.histograms {
		if ph, ok := histogram.(*PrometheusHistogram); ok {
			registry.MustRegister(ph.histogram)
		}
	}
	for _, summary := range mh.collector.summaries {
		if ps, ok := summary.(*PrometheusSummary); ok {
			registry.MustRegister(ps.summary)
		}
	}
	mh.collector.mutex.RUnlock()

	// 生成指标数据
	gatherer := prometheus.Gatherer(registry)
	metricFamilies, err := gatherer.Gather()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error gathering metrics: %v", err)
		return
	}

	// 输出指标数据
	for _, mf := range metricFamilies {
		c.String(http.StatusOK, fmt.Sprintf("# HELP %s %s\n", mf.GetName(), mf.GetHelp()))
		c.String(http.StatusOK, fmt.Sprintf("# TYPE %s %s\n", mf.GetName(), mf.GetType().String()))

		for _, metric := range mf.GetMetric() {
			c.String(http.StatusOK, fmt.Sprintf("%s %v\n", mf.GetName(), metric.GetCounter().GetValue()))
		}
	}
}

// GetCollector 获取指标收集器
func (mh *MetricsHandler) GetCollector() *MetricsCollector {
	return mh.collector
}
