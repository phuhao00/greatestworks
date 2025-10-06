// Package monitoring Prometheus监控实现
// Author: MMO Server Team
// Created: 2024

package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	dto "github.com/prometheus/client_model/go"
)

// PrometheusRegistry Prometheus注册表实现
type PrometheusRegistry struct {
	registry *prometheus.Registry
	metrics  map[string]Metric
	mutex    sync.RWMutex
}

// NewPrometheusRegistry 创建Prometheus注册表
func NewPrometheusRegistry() *PrometheusRegistry {
	return &PrometheusRegistry{
		registry: prometheus.NewRegistry(),
		metrics:  make(map[string]Metric),
	}
}

// Register 注册指标
func (pr *PrometheusRegistry) Register(metric Metric) error {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	name := metric.Name()
	if _, exists := pr.metrics[name]; exists {
		return ErrMetricExists
	}

	// 转换为Prometheus指标
	promMetric, err := pr.convertToPrometheusMetric(metric)
	if err != nil {
		return err
	}

	// 注册到Prometheus
	if err := pr.registry.Register(promMetric); err != nil {
		return fmt.Errorf("failed to register prometheus metric: %w", err)
	}

	pr.metrics[name] = metric
	return nil
}

// Unregister 注销指标
func (pr *PrometheusRegistry) Unregister(name string) error {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	metric, exists := pr.metrics[name]
	if !exists {
		return ErrMetricNotFound
	}

	// 从Prometheus注销
	promMetric, err := pr.convertToPrometheusMetric(metric)
	if err != nil {
		return err
	}

	if !pr.registry.Unregister(promMetric) {
		return fmt.Errorf("failed to unregister prometheus metric")
	}

	delete(pr.metrics, name)
	return nil
}

// Get 获取指标
func (pr *PrometheusRegistry) Get(name string) (Metric, bool) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()
	metric, exists := pr.metrics[name]
	return metric, exists
}

// GetAll 获取所有指标
func (pr *PrometheusRegistry) GetAll() map[string]Metric {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	result := make(map[string]Metric)
	for k, v := range pr.metrics {
		result[k] = v
	}
	return result
}

// Clear 清空所有指标
func (pr *PrometheusRegistry) Clear() {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	// 创建新的注册表
	pr.registry = prometheus.NewRegistry()
	pr.metrics = make(map[string]Metric)
}

// Gather 收集指标数据
func (pr *PrometheusRegistry) Gather() ([]*MetricFamily, error) {
	metricFamilies, err := pr.registry.Gather()
	if err != nil {
		return nil, err
	}

	// 转换为内部格式
	result := make([]*MetricFamily, 0, len(metricFamilies))
	for _, mf := range metricFamilies {
		converted := pr.convertFromPrometheusMetricFamily(mf)
		result = append(result, converted)
	}

	return result, nil
}

// GetPrometheusRegistry 获取底层Prometheus注册表
func (pr *PrometheusRegistry) GetPrometheusRegistry() *prometheus.Registry {
	return pr.registry
}

// convertToPrometheusMetric 转换为Prometheus指标
func (pr *PrometheusRegistry) convertToPrometheusMetric(metric Metric) (prometheus.Collector, error) {
	switch m := metric.(type) {
	case *PrometheusCounter:
		return m.counter, nil
	case *PrometheusGauge:
		return m.gauge, nil
	case *PrometheusHistogram:
		return m.histogram, nil
	case *PrometheusSummary:
		return m.summary, nil
	default:
		return nil, fmt.Errorf("unsupported metric type: %T", metric)
	}
}

// convertFromPrometheusMetricFamily 从Prometheus指标族转换
func (pr *PrometheusRegistry) convertFromPrometheusMetricFamily(mf *dto.MetricFamily) *MetricFamily {
	result := &MetricFamily{
		Name:    mf.GetName(),
		Help:    mf.GetHelp(),
		Type:    pr.convertMetricType(mf.GetType()),
		Metrics: make([]*Sample, 0, len(mf.GetMetric())),
	}

	for _, m := range mf.GetMetric() {
		sample := &Sample{
			Labels:    pr.convertLabels(m.GetLabel()),
			Timestamp: time.Now(),
		}

		switch mf.GetType() {
		case dto.MetricType_COUNTER:
			sample.Value = m.GetCounter().GetValue()
		case dto.MetricType_GAUGE:
			sample.Value = m.GetGauge().GetValue()
		case dto.MetricType_HISTOGRAM:
			hist := m.GetHistogram()
			sample.Value = hist.GetSampleSum()
			sample.Buckets = pr.convertBuckets(hist.GetBucket())
		case dto.MetricType_SUMMARY:
			summ := m.GetSummary()
			sample.Value = summ.GetSampleSum()
			sample.Quantiles = pr.convertQuantiles(summ.GetQuantile())
		}

		result.Metrics = append(result.Metrics, sample)
	}

	return result
}

// convertMetricType 转换指标类型
func (pr *PrometheusRegistry) convertMetricType(promType dto.MetricType) MetricType {
	switch promType {
	case dto.MetricType_COUNTER:
		return CounterType
	case dto.MetricType_GAUGE:
		return GaugeType
	case dto.MetricType_HISTOGRAM:
		return HistogramType
	case dto.MetricType_SUMMARY:
		return SummaryType
	default:
		return GaugeType
	}
}

// convertLabels 转换标签
func (pr *PrometheusRegistry) convertLabels(promLabels []*dto.LabelPair) Labels {
	labels := make(Labels)
	for _, label := range promLabels {
		labels[label.GetName()] = label.GetValue()
	}
	return labels
}

// convertBuckets 转换桶
func (pr *PrometheusRegistry) convertBuckets(promBuckets []*dto.Bucket) []Bucket {
	buckets := make([]Bucket, 0, len(promBuckets))
	for _, bucket := range promBuckets {
		buckets = append(buckets, Bucket{
			UpperBound: bucket.GetUpperBound(),
			Count:      bucket.GetCumulativeCount(),
		})
	}
	return buckets
}

// convertQuantiles 转换分位数
func (pr *PrometheusRegistry) convertQuantiles(promQuantiles []*dto.Quantile) []Quantile {
	quantiles := make([]Quantile, 0, len(promQuantiles))
	for _, quantile := range promQuantiles {
		quantiles = append(quantiles, Quantile{
			Quantile: quantile.GetQuantile(),
			Value:    quantile.GetValue(),
		})
	}
	return quantiles
}

// PrometheusFactory Prometheus指标工厂
type PrometheusFactory struct {
	namespace string
	subsystem string
	labels    Labels
	registry  *PrometheusRegistry
}

// NewPrometheusFactory 创建Prometheus工厂
func NewPrometheusFactory(namespace, subsystem string, labels Labels, registry *PrometheusRegistry) *PrometheusFactory {
	return &PrometheusFactory{
		namespace: namespace,
		subsystem: subsystem,
		labels:    labels,
		registry:  registry,
	}
}

// NewCounter 创建计数器
func (pf *PrometheusFactory) NewCounter(name, help string, labels Labels) Counter {
	opts := prometheus.CounterOpts{
		Namespace:   pf.namespace,
		Subsystem:   pf.subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: pf.mergeLabels(labels),
	}

	counter := prometheus.NewCounter(opts)
	return &PrometheusCounter{
		counter: counter,
		name:    pf.buildMetricName(name),
		help:    help,
		labels:  labels,
	}
}

// NewGauge 创建仪表盘
func (pf *PrometheusFactory) NewGauge(name, help string, labels Labels) GaugeInterface {
	opts := prometheus.GaugeOpts{
		Namespace:   pf.namespace,
		Subsystem:   pf.subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: pf.mergeLabels(labels),
	}

	gauge := prometheus.NewGauge(opts)
	return &PrometheusGauge{
		gauge:  gauge,
		name:   pf.buildMetricName(name),
		help:   help,
		labels: labels,
	}
}

// NewHistogram 创建直方图
func (pf *PrometheusFactory) NewHistogram(name, help string, labels Labels, buckets []float64) HistogramInterface {
	opts := prometheus.HistogramOpts{
		Namespace:   pf.namespace,
		Subsystem:   pf.subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: pf.mergeLabels(labels),
		Buckets:     buckets,
	}

	histogram := prometheus.NewHistogram(opts)
	return &PrometheusHistogram{
		histogram: histogram,
		name:      pf.buildMetricName(name),
		help:      help,
		labels:    labels,
		buckets:   buckets,
	}
}

// NewSummary 创建摘要
func (pf *PrometheusFactory) NewSummary(name, help string, labels Labels, quantiles map[float64]float64) Summary {
	objectives := make(map[float64]float64)
	for q, err := range quantiles {
		objectives[q] = err
	}

	opts := prometheus.SummaryOpts{
		Namespace:   pf.namespace,
		Subsystem:   pf.subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: pf.mergeLabels(labels),
		Objectives:  objectives,
	}

	summary := prometheus.NewSummary(opts)
	return &PrometheusSummary{
		summary:   summary,
		name:      pf.buildMetricName(name),
		help:      help,
		labels:    labels,
		quantiles: quantiles,
	}
}

// NewTimer 创建计时器
func (pf *PrometheusFactory) NewTimer(name, help string, labels Labels) Timer {
	// 使用直方图实现计时器
	buckets := prometheus.DefBuckets
	histogram := pf.NewHistogram(name+"_duration_seconds", help, labels, buckets)
	return &PrometheusTimer{
		histogram: histogram,
		name:      pf.buildMetricName(name),
		help:      help,
		labels:    labels,
	}
}

// mergeLabels 合并标签
func (pf *PrometheusFactory) mergeLabels(labels Labels) prometheus.Labels {
	result := make(prometheus.Labels)

	// 添加工厂标签
	for k, v := range pf.labels {
		result[k] = v
	}

	// 添加指标标签
	for k, v := range labels {
		result[k] = v
	}

	return result
}

// buildMetricName 构建指标名称
func (pf *PrometheusFactory) buildMetricName(name string) string {
	parts := []string{}
	if pf.namespace != "" {
		parts = append(parts, pf.namespace)
	}
	if pf.subsystem != "" {
		parts = append(parts, pf.subsystem)
	}
	parts = append(parts, name)
	return strings.Join(parts, "_")
}

// PrometheusCounter Prometheus计数器实现
type PrometheusCounter struct {
	counter prometheus.Counter
	name    string
	help    string
	labels  Labels
}

func (pc *PrometheusCounter) Name() string     { return pc.name }
func (pc *PrometheusCounter) Type() MetricType { return CounterType }
func (pc *PrometheusCounter) Help() string     { return pc.help }
func (pc *PrometheusCounter) Labels() Labels   { return pc.labels }
func (pc *PrometheusCounter) Value() interface{} {
	metric := &dto.MetricFamily{}
	pc.counter.Write(metric)
	return metric.GetMetric()[0].GetCounter().GetValue()
}
func (pc *PrometheusCounter) Reset()        { /* Prometheus计数器不支持重置 */ }
func (pc *PrometheusCounter) String() string { return fmt.Sprintf("%s{%v} = %v", pc.name, pc.labels, pc.Value()) }
func (pc *PrometheusCounter) Inc()           { pc.counter.Inc() }
func (pc *PrometheusCounter) Add(value float64) { pc.counter.Add(value) }
func (pc *PrometheusCounter) Get() float64 {
	metric := &dto.MetricFamily{}
	pc.counter.Write(metric)
	return metric.GetMetric()[0].GetCounter().GetValue()
}

// PrometheusGauge Prometheus仪表盘实现
type PrometheusGauge struct {
	gauge  prometheus.Gauge
	name   string
	help   string
	labels Labels
}

func (pg *PrometheusGauge) Name() string     { return pg.name }
func (pg *PrometheusGauge) Type() MetricType { return GaugeType }
func (pg *PrometheusGauge) Help() string     { return pg.help }
func (pg *PrometheusGauge) Labels() Labels   { return pg.labels }
func (pg *PrometheusGauge) Value() interface{} {
	metric := &dto.MetricFamily{}
	pg.gauge.Write(metric)
	return metric.GetMetric()[0].GetGauge().GetValue()
}
func (pg *PrometheusGauge) Reset()                { pg.gauge.Set(0) }
func (pg *PrometheusGauge) String() string        { return fmt.Sprintf("%s{%v} = %v", pg.name, pg.labels, pg.Value()) }
func (pg *PrometheusGauge) Set(value float64)     { pg.gauge.Set(value) }
func (pg *PrometheusGauge) Inc()                  { pg.gauge.Inc() }
func (pg *PrometheusGauge) Dec()                  { pg.gauge.Dec() }
func (pg *PrometheusGauge) Add(value float64)     { pg.gauge.Add(value) }
func (pg *PrometheusGauge) Sub(value float64)     { pg.gauge.Sub(value) }
func (pg *PrometheusGauge) Get() float64 {
	metric := &dto.MetricFamily{}
	pg.gauge.Write(metric)
	return metric.GetMetric()[0].GetGauge().GetValue()
}

// PrometheusHistogram Prometheus直方图实现
type PrometheusHistogram struct {
	histogram prometheus.Histogram
	name      string
	help      string
	labels    Labels
	buckets   []float64
}

func (ph *PrometheusHistogram) Name() string     { return ph.name }
func (ph *PrometheusHistogram) Type() MetricType { return HistogramType }
func (ph *PrometheusHistogram) Help() string     { return ph.help }
func (ph *PrometheusHistogram) Labels() Labels   { return ph.labels }
func (ph *PrometheusHistogram) Value() interface{} {
	metric := &dto.MetricFamily{}
	ph.histogram.Write(metric)
	return metric.GetMetric()[0].GetHistogram().GetSampleSum()
}
func (ph *PrometheusHistogram) Reset() { /* Prometheus直方图不支持重置 */ }
func (ph *PrometheusHistogram) String() string {
	return fmt.Sprintf("%s{%v} = %v", ph.name, ph.labels, ph.Value())
}
func (ph *PrometheusHistogram) Observe(value float64) { ph.histogram.Observe(value) }
func (ph *PrometheusHistogram) ObserveWithLabels(value float64, labels Labels) {
	// Prometheus直方图不支持动态标签
	ph.histogram.Observe(value)
}
func (ph *PrometheusHistogram) GetBuckets() []float64 { return ph.buckets }
func (ph *PrometheusHistogram) GetCounts() []uint64 {
	metric := &dto.MetricFamily{}
	ph.histogram.Write(metric)
	buckets := metric.GetMetric()[0].GetHistogram().GetBucket()
	counts := make([]uint64, len(buckets))
	for i, bucket := range buckets {
		counts[i] = bucket.GetCumulativeCount()
	}
	return counts
}
func (ph *PrometheusHistogram) GetSum() float64 {
	metric := &dto.MetricFamily{}
	ph.histogram.Write(metric)
	return metric.GetMetric()[0].GetHistogram().GetSampleSum()
}
func (ph *PrometheusHistogram) GetCount() uint64 {
	metric := &dto.MetricFamily{}
	ph.histogram.Write(metric)
	return metric.GetMetric()[0].GetHistogram().GetSampleCount()
}

// PrometheusSummary Prometheus摘要实现
type PrometheusSummary struct {
	summary   prometheus.Summary
	name      string
	help      string
	labels    Labels
	quantiles map[float64]float64
}

func (ps *PrometheusSummary) Name() string     { return ps.name }
func (ps *PrometheusSummary) Type() MetricType { return SummaryType }
func (ps *PrometheusSummary) Help() string     { return ps.help }
func (ps *PrometheusSummary) Labels() Labels   { return ps.labels }
func (ps *PrometheusSummary) Value() interface{} {
	metric := &dto.MetricFamily{}
	ps.summary.Write(metric)
	return metric.GetMetric()[0].GetSummary().GetSampleSum()
}
func (ps *PrometheusSummary) Reset() { /* Prometheus摘要不支持重置 */ }
func (ps *PrometheusSummary) String() string {
	return fmt.Sprintf("%s{%v} = %v", ps.name, ps.labels, ps.Value())
}
func (ps *PrometheusSummary) Observe(value float64) { ps.summary.Observe(value) }
func (ps *PrometheusSummary) ObserveWithLabels(value float64, labels Labels) {
	// Prometheus摘要不支持动态标签
	ps.summary.Observe(value)
}
func (ps *PrometheusSummary) GetQuantiles() map[float64]float64 {
	metric := &dto.MetricFamily{}
	ps.summary.Write(metric)
	quantiles := make(map[float64]float64)
	for _, q := range metric.GetMetric()[0].GetSummary().GetQuantile() {
		quantiles[q.GetQuantile()] = q.GetValue()
	}
	return quantiles
}
func (ps *PrometheusSummary) GetSum() float64 {
	metric := &dto.MetricFamily{}
	ps.summary.Write(metric)
	return metric.GetMetric()[0].GetSummary().GetSampleSum()
}
func (ps *PrometheusSummary) GetCount() uint64 {
	metric := &dto.MetricFamily{}
	ps.summary.Write(metric)
	return metric.GetMetric()[0].GetSummary().GetSampleCount()
}

// PrometheusTimer Prometheus计时器实现
type PrometheusTimer struct {
	histogram HistogramInterface
	name      string
	help      string
	labels    Labels
}

func (pt *PrometheusTimer) Start() TimerContext {
	return &PrometheusTimerContext{
		timer: pt,
		start: time.Now(),
	}
}

func (pt *PrometheusTimer) Time(fn func()) {
	start := time.Now()
	fn()
	pt.ObserveDuration(time.Since(start))
}

func (pt *PrometheusTimer) TimeContext(ctx context.Context, fn func(context.Context)) {
	start := time.Now()
	fn(ctx)
	pt.ObserveDuration(time.Since(start))
}

func (pt *PrometheusTimer) ObserveDuration(duration time.Duration) {
	pt.histogram.Observe(duration.Seconds())
}

// PrometheusTimerContext Prometheus计时器上下文
type PrometheusTimerContext struct {
	timer *PrometheusTimer
	start time.Time
}

func (ptc *PrometheusTimerContext) Stop() {
	duration := time.Since(ptc.start)
	ptc.timer.ObserveDuration(duration)
}

func (ptc *PrometheusTimerContext) Duration() time.Duration {
	return time.Since(ptc.start)
}

// PrometheusExporter Prometheus导出器
type PrometheusExporter struct {
	registry *PrometheusRegistry
	format   expfmt.Format
}

// NewPrometheusExporter 创建Prometheus导出器
func NewPrometheusExporter(registry *PrometheusRegistry) *PrometheusExporter {
	return &PrometheusExporter{
		registry: registry,
		format:   expfmt.FmtText,
	}
}

// Export 导出指标
func (pe *PrometheusExporter) Export(ctx context.Context, metrics []*MetricFamily) error {
	// Prometheus导出器通过HTTP端点导出，这里不需要实现
	return nil
}

// Format 获取导出格式
func (pe *PrometheusExporter) Format() string {
	return "prometheus"
}

// PrometheusServer Prometheus监控服务器
type PrometheusServer struct {
	registry *PrometheusRegistry
	server   *http.Server
	config   *Config
	mutex    sync.RWMutex
	running  bool
}

// NewPrometheusServer 创建Prometheus服务器
func NewPrometheusServer(registry *PrometheusRegistry, config *Config) *PrometheusServer {
	return &PrometheusServer{
		registry: registry,
		config:   config,
	}
}

// Start 启动服务器
func (ps *PrometheusServer) Start(ctx context.Context) error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if ps.running {
		return ErrServerAlreadyStarted
	}

	// 创建HTTP服务器
	mux := http.NewServeMux()

	// 注册指标端点
	handler := promhttp.HandlerFor(ps.registry.GetPrometheusRegistry(), promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
	mux.Handle(ps.config.Path, handler)

	// 注册健康检查端点
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf("%s:%d", ps.config.Host, ps.config.Port)
	ps.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// 启动服务器
	go func() {
		if err := ps.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Prometheus server error: %v\n", err)
		}
	}()

	ps.running = true
	return nil
}

// Stop 停止服务器
func (ps *PrometheusServer) Stop(ctx context.Context) error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if !ps.running {
		return ErrServerNotStarted
	}

	if ps.server != nil {
		if err := ps.server.Shutdown(ctx); err != nil {
			return err
		}
	}

	ps.running = false
	return nil
}

// RegisterHandler 注册处理器
func (ps *PrometheusServer) RegisterHandler(path string, handler func() ([]byte, error)) {
	// Prometheus服务器使用固定的处理器
}

// GetAddr 获取监听地址
func (ps *PrometheusServer) GetAddr() string {
	if ps.server != nil {
		return ps.server.Addr
	}
	return fmt.Sprintf("%s:%d", ps.config.Host, ps.config.Port)
}

// 便捷函数

// NewPrometheusManager 创建Prometheus管理器
func NewPrometheusManager(config *Config) Manager {
	registry := NewPrometheusRegistry()
	factory := NewPrometheusFactory(config.Namespace, config.Subsystem, config.Labels, registry)
	server := NewPrometheusServer(registry, config)

	// 注册默认收集器
	if config.EnableRuntimeMetrics {
		registry.registry.MustRegister(collectors.NewGoCollector())
	}
	if config.EnableProcessMetrics {
		registry.registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}

	return &PrometheusManager{
		registry: registry,
		factory:  factory,
		server:   server,
		config:   config,
	}
}

// PrometheusManager Prometheus管理器
type PrometheusManager struct {
	registry   *PrometheusRegistry
	factory    *PrometheusFactory
	server     *PrometheusServer
	config     *Config
	collectors []Collector
	mutex      sync.RWMutex
}

func (pm *PrometheusManager) GetRegistry() Registry { return pm.registry }
func (pm *PrometheusManager) GetFactory() Factory   { return pm.factory }

func (pm *PrometheusManager) RegisterCollector(collector Collector) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.collectors = append(pm.collectors, collector)
	return nil
}

func (pm *PrometheusManager) UnregisterCollector(collector Collector) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	for i, c := range pm.collectors {
		if c == collector {
			pm.collectors = append(pm.collectors[:i], pm.collectors[i+1:]...)
			return nil
		}
	}
	return ErrCollectorNotFound
}

func (pm *PrometheusManager) StartServer(ctx context.Context) error {
	return pm.server.Start(ctx)
}

func (pm *PrometheusManager) StopServer(ctx context.Context) error {
	return pm.server.Stop(ctx)
}

func (pm *PrometheusManager) Export(ctx context.Context, format string) ([]byte, error) {
	metrics, err := pm.registry.Gather()
	if err != nil {
		return nil, err
	}

	// 简单的文本格式导出
	var result strings.Builder
	for _, mf := range metrics {
		result.WriteString(fmt.Sprintf("# HELP %s %s\n", mf.Name, mf.Help))
		result.WriteString(fmt.Sprintf("# TYPE %s %s\n", mf.Name, mf.Type))
		for _, sample := range mf.Metrics {
			labelStr := pm.formatLabels(sample.Labels)
			result.WriteString(fmt.Sprintf("%s%s %g\n", mf.Name, labelStr, sample.Value))
		}
	}

	return []byte(result.String()), nil
}

func (pm *PrometheusManager) GetMetrics() ([]*MetricFamily, error) {
	return pm.registry.Gather()
}

// formatLabels 格式化标签
func (pm *PrometheusManager) formatLabels(labels Labels) string {
	if len(labels) == 0 {
		return ""
	}

	var pairs []string
	for k, v := range labels {
		pairs = append(pairs, fmt.Sprintf(`%s="%s"`, k, v))
	}
	sort.Strings(pairs)
	return "{" + strings.Join(pairs, ",") + "}"
}