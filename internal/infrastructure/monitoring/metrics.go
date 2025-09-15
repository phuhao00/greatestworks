// Package monitoring 统一监控系统
// Author: MMO Server Team
// Created: 2024

package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MetricType 指标类型
type MetricType string

const (
	// CounterType 计数器类型
	CounterType MetricType = "counter"
	// GaugeType 仪表盘类型
	GaugeType MetricType = "gauge"
	// HistogramType 直方图类型
	HistogramType MetricType = "histogram"
	// SummaryType 摘要类型
	SummaryType MetricType = "summary"
)

// Labels 标签映射
type Labels map[string]string

// Metric 指标接口
type Metric interface {
	// Name 获取指标名称
	Name() string
	// Type 获取指标类型
	Type() MetricType
	// Help 获取帮助信息
	Help() string
	// Labels 获取标签
	Labels() Labels
	// Value 获取当前值
	Value() interface{}
	// Reset 重置指标
	Reset()
	// String 字符串表示
	String() string
}

// Counter 计数器接口
type Counter interface {
	Metric
	// Inc 增加计数（默认增加1）
	Inc()
	// Add 增加指定值
	Add(value float64)
	// Get 获取当前计数
	Get() float64
}

// Gauge 仪表盘接口
type Gauge interface {
	Metric
	// Set 设置值
	Set(value float64)
	// Inc 增加值（默认增加1）
	Inc()
	// Dec 减少值（默认减少1）
	Dec()
	// Add 增加指定值
	Add(value float64)
	// Sub 减少指定值
	Sub(value float64)
	// Get 获取当前值
	Get() float64
}

// Histogram 直方图接口
type Histogram interface {
	Metric
	// Observe 观察值
	Observe(value float64)
	// ObserveWithLabels 带标签观察值
	ObserveWithLabels(value float64, labels Labels)
	// GetBuckets 获取桶信息
	GetBuckets() []float64
	// GetCounts 获取计数信息
	GetCounts() []uint64
	// GetSum 获取总和
	GetSum() float64
	// GetCount 获取总计数
	GetCount() uint64
}

// Summary 摘要接口
type Summary interface {
	Metric
	// Observe 观察值
	Observe(value float64)
	// ObserveWithLabels 带标签观察值
	ObserveWithLabels(value float64, labels Labels)
	// GetQuantiles 获取分位数
	GetQuantiles() map[float64]float64
	// GetSum 获取总和
	GetSum() float64
	// GetCount 获取总计数
	GetCount() uint64
}

// Timer 计时器接口
type Timer interface {
	// Start 开始计时
	Start() TimerContext
	// Time 计时函数执行时间
	Time(fn func())
	// TimeContext 计时函数执行时间（带上下文）
	TimeContext(ctx context.Context, fn func(context.Context))
	// ObserveDuration 观察持续时间
	ObserveDuration(duration time.Duration)
}

// TimerContext 计时器上下文
type TimerContext interface {
	// Stop 停止计时并记录
	Stop()
	// Duration 获取已经过的时间
	Duration() time.Duration
}

// Registry 指标注册表接口
type Registry interface {
	// Register 注册指标
	Register(metric Metric) error
	// Unregister 注销指标
	Unregister(name string) error
	// Get 获取指标
	Get(name string) (Metric, bool)
	// GetAll 获取所有指标
	GetAll() map[string]Metric
	// Clear 清空所有指标
	Clear()
	// Gather 收集指标数据
	Gather() ([]*MetricFamily, error)
}

// MetricFamily 指标族
type MetricFamily struct {
	Name    string      `json:"name"`
	Help    string      `json:"help"`
	Type    MetricType  `json:"type"`
	Metrics []*Sample   `json:"metrics"`
}

// Sample 指标样本
type Sample struct {
	Labels    Labels      `json:"labels"`
	Value     float64     `json:"value"`
	Timestamp time.Time   `json:"timestamp"`
	Buckets   []Bucket    `json:"buckets,omitempty"`
	Quantiles []Quantile  `json:"quantiles,omitempty"`
}

// Bucket 直方图桶
type Bucket struct {
	UpperBound float64 `json:"upper_bound"`
	Count      uint64  `json:"count"`
}

// Quantile 分位数
type Quantile struct {
	Quantile float64 `json:"quantile"`
	Value    float64 `json:"value"`
}

// Factory 指标工厂接口
type Factory interface {
	// NewCounter 创建计数器
	NewCounter(name, help string, labels Labels) Counter
	// NewGauge 创建仪表盘
	NewGauge(name, help string, labels Labels) Gauge
	// NewHistogram 创建直方图
	NewHistogram(name, help string, labels Labels, buckets []float64) Histogram
	// NewSummary 创建摘要
	NewSummary(name, help string, labels Labels, quantiles map[float64]float64) Summary
	// NewTimer 创建计时器
	NewTimer(name, help string, labels Labels) Timer
}

// Collector 收集器接口
type Collector interface {
	// Describe 描述指标
	Describe(ch chan<- *MetricDesc)
	// Collect 收集指标
	Collect(ch chan<- Metric)
}

// MetricDesc 指标描述
type MetricDesc struct {
	Name      string
	Help      string
	Type      MetricType
	Labels    []string
	ConstLabels Labels
}

// Exporter 导出器接口
type Exporter interface {
	// Export 导出指标
	Export(ctx context.Context, metrics []*MetricFamily) error
	// Format 获取导出格式
	Format() string
}

// Server 监控服务器接口
type Server interface {
	// Start 启动服务器
	Start(ctx context.Context) error
	// Stop 停止服务器
	Stop(ctx context.Context) error
	// RegisterHandler 注册处理器
	RegisterHandler(path string, handler func() ([]byte, error))
	// GetAddr 获取监听地址
	GetAddr() string
}

// Config 监控配置
type Config struct {
	// 是否启用监控
	Enabled bool `yaml:"enabled" json:"enabled"`
	// 监听端口
	Port int `yaml:"port" json:"port"`
	// 监听地址
	Host string `yaml:"host" json:"host"`
	// 指标路径
	Path string `yaml:"path" json:"path"`
	// 命名空间
	Namespace string `yaml:"namespace" json:"namespace"`
	// 子系统
	Subsystem string `yaml:"subsystem" json:"subsystem"`
	// 标签
	Labels map[string]string `yaml:"labels" json:"labels"`
	// 收集间隔
	CollectInterval time.Duration `yaml:"collect_interval" json:"collect_interval"`
	// 是否启用运行时指标
	EnableRuntimeMetrics bool `yaml:"enable_runtime_metrics" json:"enable_runtime_metrics"`
	// 是否启用进程指标
	EnableProcessMetrics bool `yaml:"enable_process_metrics" json:"enable_process_metrics"`
	// 是否启用Go指标
	EnableGoMetrics bool `yaml:"enable_go_metrics" json:"enable_go_metrics"`
	// 自定义收集器
	CustomCollectors []string `yaml:"custom_collectors" json:"custom_collectors"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled:              true,
		Port:                 9090,
		Host:                 "0.0.0.0",
		Path:                 "/metrics",
		Namespace:            "mmo",
		Subsystem:            "server",
		Labels:               make(map[string]string),
		CollectInterval:      15 * time.Second,
		EnableRuntimeMetrics: true,
		EnableProcessMetrics: true,
		EnableGoMetrics:      true,
		CustomCollectors:     make([]string, 0),
	}
}

// Manager 监控管理器接口
type Manager interface {
	// GetRegistry 获取注册表
	GetRegistry() Registry
	// GetFactory 获取工厂
	GetFactory() Factory
	// RegisterCollector 注册收集器
	RegisterCollector(collector Collector) error
	// UnregisterCollector 注销收集器
	UnregisterCollector(collector Collector) error
	// StartServer 启动监控服务器
	StartServer(ctx context.Context) error
	// StopServer 停止监控服务器
	StopServer(ctx context.Context) error
	// Export 导出指标
	Export(ctx context.Context, format string) ([]byte, error)
	// GetMetrics 获取指标
	GetMetrics() ([]*MetricFamily, error)
}

// 预定义指标名称
const (
	// HTTP相关指标
	HTTPRequestsTotal     = "http_requests_total"
	HTTPRequestDuration   = "http_request_duration_seconds"
	HTTPRequestSize       = "http_request_size_bytes"
	HTTPResponseSize      = "http_response_size_bytes"
	HTTPRequestsInFlight  = "http_requests_in_flight"

	// 数据库相关指标
	DBConnectionsOpen     = "db_connections_open"
	DBConnectionsIdle     = "db_connections_idle"
	DBConnectionsInUse    = "db_connections_in_use"
	DBQueryDuration       = "db_query_duration_seconds"
	DBQueriesTotal        = "db_queries_total"

	// 缓存相关指标
	CacheHitsTotal        = "cache_hits_total"
	CacheMissesTotal      = "cache_misses_total"
	CacheOperationDuration = "cache_operation_duration_seconds"
	CacheSize             = "cache_size_bytes"
	CacheEntries          = "cache_entries"

	// 游戏相关指标
	PlayersOnline         = "players_online"
	PlayerActions         = "player_actions_total"
	GameEvents            = "game_events_total"
	BattlesActive         = "battles_active"
	GuildMembers          = "guild_members_total"

	// 系统相关指标
	CPUUsage              = "cpu_usage_percent"
	MemoryUsage           = "memory_usage_bytes"
	DiskUsage             = "disk_usage_bytes"
	NetworkBytesReceived  = "network_bytes_received_total"
	NetworkBytesSent      = "network_bytes_sent_total"

	// 错误相关指标
	ErrorsTotal           = "errors_total"
	PanicsTotal           = "panics_total"
	TimeoutsTotal         = "timeouts_total"
)

// 预定义标签名称
const (
	// HTTP标签
	LabelMethod     = "method"
	LabelPath       = "path"
	LabelStatusCode = "status_code"
	LabelHandler    = "handler"

	// 数据库标签
	LabelDatabase   = "database"
	LabelTable      = "table"
	LabelOperation  = "operation"
	LabelQuery      = "query"

	// 缓存标签
	LabelCacheType  = "cache_type"
	LabelCacheKey   = "cache_key"

	// 游戏标签
	LabelPlayerID   = "player_id"
	LabelAction     = "action"
	LabelEventType  = "event_type"
	LabelBattleType = "battle_type"
	LabelGuildID    = "guild_id"

	// 系统标签
	LabelComponent  = "component"
	LabelModule     = "module"
	LabelService    = "service"
	LabelInstance   = "instance"
	LabelVersion    = "version"

	// 错误标签
	LabelErrorType  = "error_type"
	LabelErrorCode  = "error_code"
)

// 预定义错误
var (
	ErrMetricNotFound     = fmt.Errorf("metric not found")
	ErrMetricExists       = fmt.Errorf("metric already exists")
	ErrInvalidMetricType  = fmt.Errorf("invalid metric type")
	ErrInvalidMetricName  = fmt.Errorf("invalid metric name")
	ErrCollectorExists    = fmt.Errorf("collector already exists")
	ErrCollectorNotFound  = fmt.Errorf("collector not found")
	ErrServerNotStarted   = fmt.Errorf("server not started")
	ErrServerAlreadyStarted = fmt.Errorf("server already started")
)

// MetricOptions 指标选项
type MetricOptions struct {
	Name        string
	Help        string
	Labels      Labels
	ConstLabels Labels
	Buckets     []float64
	Quantiles   map[float64]float64
	MaxAge      time.Duration
	AgeBuckets  int
	BufCap      int
}

// NewMetricOptions 创建指标选项
func NewMetricOptions(name, help string) *MetricOptions {
	return &MetricOptions{
		Name:   name,
		Help:   help,
		Labels: make(Labels),
		ConstLabels: make(Labels),
	}
}

// WithLabels 添加标签
func (mo *MetricOptions) WithLabels(labels Labels) *MetricOptions {
	for k, v := range labels {
		mo.Labels[k] = v
	}
	return mo
}

// WithConstLabels 添加常量标签
func (mo *MetricOptions) WithConstLabels(labels Labels) *MetricOptions {
	for k, v := range labels {
		mo.ConstLabels[k] = v
	}
	return mo
}

// WithBuckets 设置直方图桶
func (mo *MetricOptions) WithBuckets(buckets []float64) *MetricOptions {
	mo.Buckets = buckets
	return mo
}

// WithQuantiles 设置摘要分位数
func (mo *MetricOptions) WithQuantiles(quantiles map[float64]float64) *MetricOptions {
	mo.Quantiles = quantiles
	return mo
}

// 全局变量
var (
	defaultRegistry Registry
	defaultFactory  Factory
	defaultManager  Manager
	registryMutex   sync.RWMutex
)

// SetDefaultRegistry 设置默认注册表
func SetDefaultRegistry(registry Registry) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	defaultRegistry = registry
}

// GetDefaultRegistry 获取默认注册表
func GetDefaultRegistry() Registry {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	return defaultRegistry
}

// SetDefaultFactory 设置默认工厂
func SetDefaultFactory(factory Factory) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	defaultFactory = factory
}

// GetDefaultFactory 获取默认工厂
func GetDefaultFactory() Factory {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	return defaultFactory
}

// SetDefaultManager 设置默认管理器
func SetDefaultManager(manager Manager) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	defaultManager = manager
}

// GetDefaultManager 获取默认管理器
func GetDefaultManager() Manager {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	return defaultManager
}

// 便捷函数

// NewCounter 创建计数器
func NewCounter(name, help string, labels Labels) Counter {
	if defaultFactory != nil {
		return defaultFactory.NewCounter(name, help, labels)
	}
	return nil
}

// NewGauge 创建仪表盘
func NewGauge(name, help string, labels Labels) Gauge {
	if defaultFactory != nil {
		return defaultFactory.NewGauge(name, help, labels)
	}
	return nil
}

// NewHistogram 创建直方图
func NewHistogram(name, help string, labels Labels, buckets []float64) Histogram {
	if defaultFactory != nil {
		return defaultFactory.NewHistogram(name, help, labels, buckets)
	}
	return nil
}

// NewSummary 创建摘要
func NewSummary(name, help string, labels Labels, quantiles map[float64]float64) Summary {
	if defaultFactory != nil {
		return defaultFactory.NewSummary(name, help, labels, quantiles)
	}
	return nil
}

// NewTimer 创建计时器
func NewTimer(name, help string, labels Labels) Timer {
	if defaultFactory != nil {
		return defaultFactory.NewTimer(name, help, labels)
	}
	return nil
}

// Register 注册指标
func Register(metric Metric) error {
	if defaultRegistry != nil {
		return defaultRegistry.Register(metric)
	}
	return ErrMetricNotFound
}

// Unregister 注销指标
func Unregister(name string) error {
	if defaultRegistry != nil {
		return defaultRegistry.Unregister(name)
	}
	return ErrMetricNotFound
}

// Get 获取指标
func Get(name string) (Metric, bool) {
	if defaultRegistry != nil {
		return defaultRegistry.Get(name)
	}
	return nil, false
}

// GetAll 获取所有指标
func GetAll() map[string]Metric {
	if defaultRegistry != nil {
		return defaultRegistry.GetAll()
	}
	return make(map[string]Metric)
}

// Gather 收集指标数据
func Gather() ([]*MetricFamily, error) {
	if defaultRegistry != nil {
		return defaultRegistry.Gather()
	}
	return nil, ErrMetricNotFound
}