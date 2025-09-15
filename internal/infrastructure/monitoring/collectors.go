// Package monitoring 指标收集器
// Author: MMO Server Team
// Created: 2024

package monitoring

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// SystemCollector 系统指标收集器
type SystemCollector struct {
	factory Factory
	metrics map[string]Metric
	mutex   sync.RWMutex
	enabled bool
}

// NewSystemCollector 创建系统收集器
func NewSystemCollector(factory Factory) *SystemCollector {
	sc := &SystemCollector{
		factory: factory,
		metrics: make(map[string]Metric),
		enabled: true,
	}

	sc.initMetrics()
	return sc
}

// initMetrics 初始化指标
func (sc *SystemCollector) initMetrics() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// CPU指标
	sc.metrics["cpu_usage_percent"] = sc.factory.NewGauge(
		"cpu_usage_percent",
		"CPU usage percentage",
		Labels{"type": "total"},
	)

	// 内存指标
	sc.metrics["memory_usage_bytes"] = sc.factory.NewGauge(
		"memory_usage_bytes",
		"Memory usage in bytes",
		Labels{"type": "used"},
	)

	sc.metrics["memory_total_bytes"] = sc.factory.NewGauge(
		"memory_total_bytes",
		"Total memory in bytes",
		Labels{},
	)

	sc.metrics["memory_available_bytes"] = sc.factory.NewGauge(
		"memory_available_bytes",
		"Available memory in bytes",
		Labels{},
	)

	// 磁盘指标
	sc.metrics["disk_usage_bytes"] = sc.factory.NewGauge(
		"disk_usage_bytes",
		"Disk usage in bytes",
		Labels{"device": "root", "type": "used"},
	)

	sc.metrics["disk_total_bytes"] = sc.factory.NewGauge(
		"disk_total_bytes",
		"Total disk space in bytes",
		Labels{"device": "root"},
	)

	// 网络指标
	sc.metrics["network_bytes_received_total"] = sc.factory.NewCounter(
		"network_bytes_received_total",
		"Total bytes received",
		Labels{},
	)

	sc.metrics["network_bytes_sent_total"] = sc.factory.NewCounter(
		"network_bytes_sent_total",
		"Total bytes sent",
		Labels{},
	)
}

// Describe 描述指标
func (sc *SystemCollector) Describe(ch chan<- *MetricDesc) {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	for name, metric := range sc.metrics {
		ch <- &MetricDesc{
			Name: name,
			Help: metric.Help(),
			Type: metric.Type(),
		}
	}
}

// Collect 收集指标
func (sc *SystemCollector) Collect(ch chan<- Metric) {
	if !sc.enabled {
		return
	}

	sc.collectCPUMetrics()
	sc.collectMemoryMetrics()
	sc.collectDiskMetrics()
	sc.collectNetworkMetrics()

	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	for _, metric := range sc.metrics {
		ch <- metric
	}
}

// collectCPUMetrics 收集CPU指标
func (sc *SystemCollector) collectCPUMetrics() {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil || len(percent) == 0 {
		return
	}

	if gauge, ok := sc.metrics["cpu_usage_percent"].(Gauge); ok {
		gauge.Set(percent[0])
	}
}

// collectMemoryMetrics 收集内存指标
func (sc *SystemCollector) collectMemoryMetrics() {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	if gauge, ok := sc.metrics["memory_usage_bytes"].(Gauge); ok {
		gauge.Set(float64(vmem.Used))
	}

	if gauge, ok := sc.metrics["memory_total_bytes"].(Gauge); ok {
		gauge.Set(float64(vmem.Total))
	}

	if gauge, ok := sc.metrics["memory_available_bytes"].(Gauge); ok {
		gauge.Set(float64(vmem.Available))
	}
}

// collectDiskMetrics 收集磁盘指标
func (sc *SystemCollector) collectDiskMetrics() {
	usage, err := disk.Usage("/")
	if err != nil {
		return
	}

	if gauge, ok := sc.metrics["disk_usage_bytes"].(Gauge); ok {
		gauge.Set(float64(usage.Used))
	}

	if gauge, ok := sc.metrics["disk_total_bytes"].(Gauge); ok {
		gauge.Set(float64(usage.Total))
	}
}

// collectNetworkMetrics 收集网络指标
func (sc *SystemCollector) collectNetworkMetrics() {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return
	}

	if counter, ok := sc.metrics["network_bytes_received_total"].(Counter); ok {
		// 注意：这里应该计算增量，但为了简化，直接设置总值
		counter.Add(float64(stats[0].BytesRecv) - counter.Get())
	}

	if counter, ok := sc.metrics["network_bytes_sent_total"].(Counter); ok {
		counter.Add(float64(stats[0].BytesSent) - counter.Get())
	}
}

// Enable 启用收集器
func (sc *SystemCollector) Enable() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.enabled = true
}

// Disable 禁用收集器
func (sc *SystemCollector) Disable() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.enabled = false
}

// RuntimeCollector 运行时指标收集器
type RuntimeCollector struct {
	factory Factory
	metrics map[string]Metric
	mutex   sync.RWMutex
	enabled bool
}

// NewRuntimeCollector 创建运行时收集器
func NewRuntimeCollector(factory Factory) *RuntimeCollector {
	rc := &RuntimeCollector{
		factory: factory,
		metrics: make(map[string]Metric),
		enabled: true,
	}

	rc.initMetrics()
	return rc
}

// initMetrics 初始化指标
func (rc *RuntimeCollector) initMetrics() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	// Go运行时指标
	rc.metrics["go_goroutines"] = rc.factory.NewGauge(
		"go_goroutines",
		"Number of goroutines",
		Labels{},
	)

	rc.metrics["go_threads"] = rc.factory.NewGauge(
		"go_threads",
		"Number of OS threads",
		Labels{},
	)

	rc.metrics["go_memstats_alloc_bytes"] = rc.factory.NewGauge(
		"go_memstats_alloc_bytes",
		"Number of bytes allocated and still in use",
		Labels{},
	)

	rc.metrics["go_memstats_total_alloc_bytes"] = rc.factory.NewCounter(
		"go_memstats_total_alloc_bytes_total",
		"Total number of bytes allocated",
		Labels{},
	)

	rc.metrics["go_memstats_sys_bytes"] = rc.factory.NewGauge(
		"go_memstats_sys_bytes",
		"Number of bytes obtained from system",
		Labels{},
	)

	rc.metrics["go_memstats_heap_alloc_bytes"] = rc.factory.NewGauge(
		"go_memstats_heap_alloc_bytes",
		"Number of heap bytes allocated and still in use",
		Labels{},
	)

	rc.metrics["go_memstats_heap_sys_bytes"] = rc.factory.NewGauge(
		"go_memstats_heap_sys_bytes",
		"Number of heap bytes obtained from system",
		Labels{},
	)

	rc.metrics["go_memstats_heap_objects"] = rc.factory.NewGauge(
		"go_memstats_heap_objects",
		"Number of allocated objects",
		Labels{},
	)

	rc.metrics["go_memstats_gc_total"] = rc.factory.NewCounter(
		"go_memstats_gc_total",
		"Total number of GC runs",
		Labels{},
	)

	rc.metrics["go_gc_duration_seconds"] = rc.factory.NewHistogram(
		"go_gc_duration_seconds",
		"Time spent in garbage collection",
		Labels{},
		[]float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0},
	)
}

// Describe 描述指标
func (rc *RuntimeCollector) Describe(ch chan<- *MetricDesc) {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	for name, metric := range rc.metrics {
		ch <- &MetricDesc{
			Name: name,
			Help: metric.Help(),
			Type: metric.Type(),
		}
	}
}

// Collect 收集指标
func (rc *RuntimeCollector) Collect(ch chan<- Metric) {
	if !rc.enabled {
		return
	}

	rc.collectRuntimeMetrics()

	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	for _, metric := range rc.metrics {
		ch <- metric
	}
}

// collectRuntimeMetrics 收集运行时指标
func (rc *RuntimeCollector) collectRuntimeMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Goroutines
	if gauge, ok := rc.metrics["go_goroutines"].(Gauge); ok {
		gauge.Set(float64(runtime.NumGoroutine()))
	}

	// Threads
	if gauge, ok := rc.metrics["go_threads"].(Gauge); ok {
		var numThreads int
		runtime.SetFinalizer(&numThreads, nil)
		gauge.Set(float64(runtime.GOMAXPROCS(0)))
	}

	// Memory stats
	if gauge, ok := rc.metrics["go_memstats_alloc_bytes"].(Gauge); ok {
		gauge.Set(float64(m.Alloc))
	}

	if counter, ok := rc.metrics["go_memstats_total_alloc_bytes"].(Counter); ok {
		current := counter.Get()
		if float64(m.TotalAlloc) > current {
			counter.Add(float64(m.TotalAlloc) - current)
		}
	}

	if gauge, ok := rc.metrics["go_memstats_sys_bytes"].(Gauge); ok {
		gauge.Set(float64(m.Sys))
	}

	if gauge, ok := rc.metrics["go_memstats_heap_alloc_bytes"].(Gauge); ok {
		gauge.Set(float64(m.HeapAlloc))
	}

	if gauge, ok := rc.metrics["go_memstats_heap_sys_bytes"].(Gauge); ok {
		gauge.Set(float64(m.HeapSys))
	}

	if gauge, ok := rc.metrics["go_memstats_heap_objects"].(Gauge); ok {
		gauge.Set(float64(m.HeapObjects))
	}

	if counter, ok := rc.metrics["go_memstats_gc_total"].(Counter); ok {
		current := counter.Get()
		if float64(m.NumGC) > current {
			counter.Add(float64(m.NumGC) - current)
		}
	}

	// GC duration (简化实现)
	if histogram, ok := rc.metrics["go_gc_duration_seconds"].(Histogram); ok {
		// 这里应该记录实际的GC时间，但为了简化，使用PauseNs的平均值
		if m.NumGC > 0 {
			avgPause := float64(m.PauseTotalNs) / float64(m.NumGC) / 1e9
			histogram.Observe(avgPause)
		}
	}
}

// Enable 启用收集器
func (rc *RuntimeCollector) Enable() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rc.enabled = true
}

// Disable 禁用收集器
func (rc *RuntimeCollector) Disable() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rc.enabled = false
}

// ProcessCollector 进程指标收集器
type ProcessCollector struct {
	factory Factory
	metrics map[string]Metric
	mutex   sync.RWMutex
	enabled bool
	process *process.Process
}

// NewProcessCollector 创建进程收集器
func NewProcessCollector(factory Factory) *ProcessCollector {
	proc, _ := process.NewProcess(int32(runtime.GOMAXPROCS(0)))

	pc := &ProcessCollector{
		factory: factory,
		metrics: make(map[string]Metric),
		enabled: true,
		process: proc,
	}

	pc.initMetrics()
	return pc
}

// initMetrics 初始化指标
func (pc *ProcessCollector) initMetrics() {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	// 进程指标
	pc.metrics["process_cpu_seconds_total"] = pc.factory.NewCounter(
		"process_cpu_seconds_total",
		"Total user and system CPU time spent in seconds",
		Labels{},
	)

	pc.metrics["process_resident_memory_bytes"] = pc.factory.NewGauge(
		"process_resident_memory_bytes",
		"Resident memory size in bytes",
		Labels{},
	)

	pc.metrics["process_virtual_memory_bytes"] = pc.factory.NewGauge(
		"process_virtual_memory_bytes",
		"Virtual memory size in bytes",
		Labels{},
	)

	pc.metrics["process_open_fds"] = pc.factory.NewGauge(
		"process_open_fds",
		"Number of open file descriptors",
		Labels{},
	)

	pc.metrics["process_start_time_seconds"] = pc.factory.NewGauge(
		"process_start_time_seconds",
		"Start time of the process since unix epoch in seconds",
		Labels{},
	)
}

// Describe 描述指标
func (pc *ProcessCollector) Describe(ch chan<- *MetricDesc) {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	for name, metric := range pc.metrics {
		ch <- &MetricDesc{
			Name: name,
			Help: metric.Help(),
			Type: metric.Type(),
		}
	}
}

// Collect 收集指标
func (pc *ProcessCollector) Collect(ch chan<- Metric) {
	if !pc.enabled || pc.process == nil {
		return
	}

	pc.collectProcessMetrics()

	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	for _, metric := range pc.metrics {
		ch <- metric
	}
}

// collectProcessMetrics 收集进程指标
func (pc *ProcessCollector) collectProcessMetrics() {
	// CPU时间
	if times, err := pc.process.Times(); err == nil {
		if counter, ok := pc.metrics["process_cpu_seconds_total"].(Counter); ok {
			totalCPU := times.User + times.System
			current := counter.Get()
			if totalCPU > current {
				counter.Add(totalCPU - current)
			}
		}
	}

	// 内存信息
	if memInfo, err := pc.process.MemoryInfo(); err == nil {
		if gauge, ok := pc.metrics["process_resident_memory_bytes"].(Gauge); ok {
			gauge.Set(float64(memInfo.RSS))
		}

		if gauge, ok := pc.metrics["process_virtual_memory_bytes"].(Gauge); ok {
			gauge.Set(float64(memInfo.VMS))
		}
	}

	// 文件描述符
	if fds, err := pc.process.NumFDs(); err == nil {
		if gauge, ok := pc.metrics["process_open_fds"].(Gauge); ok {
			gauge.Set(float64(fds))
		}
	}

	// 启动时间
	if createTime, err := pc.process.CreateTime(); err == nil {
		if gauge, ok := pc.metrics["process_start_time_seconds"].(Gauge); ok {
			gauge.Set(float64(createTime) / 1000) // 转换为秒
		}
	}
}

// Enable 启用收集器
func (pc *ProcessCollector) Enable() {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	pc.enabled = true
}

// Disable 禁用收集器
func (pc *ProcessCollector) Disable() {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	pc.enabled = false
}

// GameCollector 游戏指标收集器
type GameCollector struct {
	factory Factory
	metrics map[string]Metric
	mutex   sync.RWMutex
	enabled bool
}

// NewGameCollector 创建游戏收集器
func NewGameCollector(factory Factory) *GameCollector {
	gc := &GameCollector{
		factory: factory,
		metrics: make(map[string]Metric),
		enabled: true,
	}

	gc.initMetrics()
	return gc
}

// initMetrics 初始化指标
func (gc *GameCollector) initMetrics() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()

	// 游戏指标
	gc.metrics["players_online"] = gc.factory.NewGauge(
		"players_online",
		"Number of players currently online",
		Labels{},
	)

	gc.metrics["player_actions_total"] = gc.factory.NewCounter(
		"player_actions_total",
		"Total number of player actions",
		Labels{"action": "unknown"},
	)

	gc.metrics["game_events_total"] = gc.factory.NewCounter(
		"game_events_total",
		"Total number of game events",
		Labels{"event_type": "unknown"},
	)

	gc.metrics["battles_active"] = gc.factory.NewGauge(
		"battles_active",
		"Number of active battles",
		Labels{},
	)

	gc.metrics["guild_members_total"] = gc.factory.NewGauge(
		"guild_members_total",
		"Total number of guild members",
		Labels{},
	)

	gc.metrics["player_level_distribution"] = gc.factory.NewHistogram(
		"player_level_distribution",
		"Distribution of player levels",
		Labels{},
		[]float64{1, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
	)

	gc.metrics["session_duration_seconds"] = gc.factory.NewHistogram(
		"session_duration_seconds",
		"Player session duration in seconds",
		Labels{},
		[]float64{60, 300, 900, 1800, 3600, 7200, 14400, 28800},
	)
}

// Describe 描述指标
func (gc *GameCollector) Describe(ch chan<- *MetricDesc) {
	gc.mutex.RLock()
	defer gc.mutex.RUnlock()

	for name, metric := range gc.metrics {
		ch <- &MetricDesc{
			Name: name,
			Help: metric.Help(),
			Type: metric.Type(),
		}
	}
}

// Collect 收集指标
func (gc *GameCollector) Collect(ch chan<- Metric) {
	if !gc.enabled {
		return
	}

	gc.mutex.RLock()
	defer gc.mutex.RUnlock()

	for _, metric := range gc.metrics {
		ch <- metric
	}
}

// RecordPlayerAction 记录玩家行为
func (gc *GameCollector) RecordPlayerAction(action string) {
	if counter, ok := gc.metrics["player_actions_total"].(Counter); ok {
		counter.Inc()
	}
}

// RecordGameEvent 记录游戏事件
func (gc *GameCollector) RecordGameEvent(eventType string) {
	if counter, ok := gc.metrics["game_events_total"].(Counter); ok {
		counter.Inc()
	}
}

// SetPlayersOnline 设置在线玩家数
func (gc *GameCollector) SetPlayersOnline(count int) {
	if gauge, ok := gc.metrics["players_online"].(Gauge); ok {
		gauge.Set(float64(count))
	}
}

// SetActiveBattles 设置活跃战斗数
func (gc *GameCollector) SetActiveBattles(count int) {
	if gauge, ok := gc.metrics["battles_active"].(Gauge); ok {
		gauge.Set(float64(count))
	}
}

// SetGuildMembers 设置公会成员总数
func (gc *GameCollector) SetGuildMembers(count int) {
	if gauge, ok := gc.metrics["guild_members_total"].(Gauge); ok {
		gauge.Set(float64(count))
	}
}

// RecordPlayerLevel 记录玩家等级
func (gc *GameCollector) RecordPlayerLevel(level int) {
	if histogram, ok := gc.metrics["player_level_distribution"].(Histogram); ok {
		histogram.Observe(float64(level))
	}
}

// RecordSessionDuration 记录会话持续时间
func (gc *GameCollector) RecordSessionDuration(duration time.Duration) {
	if histogram, ok := gc.metrics["session_duration_seconds"].(Histogram); ok {
		histogram.Observe(duration.Seconds())
	}
}

// Enable 启用收集器
func (gc *GameCollector) Enable() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	gc.enabled = true
}

// Disable 禁用收集器
func (gc *GameCollector) Disable() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	gc.enabled = false
}

// HTTPCollector HTTP指标收集器
type HTTPCollector struct {
	factory Factory
	metrics map[string]Metric
	mutex   sync.RWMutex
	enabled bool
}

// NewHTTPCollector 创建HTTP收集器
func NewHTTPCollector(factory Factory) *HTTPCollector {
	hc := &HTTPCollector{
		factory: factory,
		metrics: make(map[string]Metric),
		enabled: true,
	}

	hc.initMetrics()
	return hc
}

// initMetrics 初始化指标
func (hc *HTTPCollector) initMetrics() {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()

	// HTTP指标
	hc.metrics["http_requests_total"] = hc.factory.NewCounter(
		"http_requests_total",
		"Total number of HTTP requests",
		Labels{"method": "unknown", "status_code": "unknown"},
	)

	hc.metrics["http_request_duration_seconds"] = hc.factory.NewHistogram(
		"http_request_duration_seconds",
		"HTTP request duration in seconds",
		Labels{"method": "unknown", "status_code": "unknown"},
		[]float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 5.0, 10.0},
	)

	hc.metrics["http_request_size_bytes"] = hc.factory.NewHistogram(
		"http_request_size_bytes",
		"HTTP request size in bytes",
		Labels{"method": "unknown"},
		[]float64{100, 1000, 10000, 100000, 1000000},
	)

	hc.metrics["http_response_size_bytes"] = hc.factory.NewHistogram(
		"http_response_size_bytes",
		"HTTP response size in bytes",
		Labels{"method": "unknown", "status_code": "unknown"},
		[]float64{100, 1000, 10000, 100000, 1000000},
	)

	hc.metrics["http_requests_in_flight"] = hc.factory.NewGauge(
		"http_requests_in_flight",
		"Number of HTTP requests currently being processed",
		Labels{},
	)
}

// Describe 描述指标
func (hc *HTTPCollector) Describe(ch chan<- *MetricDesc) {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()

	for name, metric := range hc.metrics {
		ch <- &MetricDesc{
			Name: name,
			Help: metric.Help(),
			Type: metric.Type(),
		}
	}
}

// Collect 收集指标
func (hc *HTTPCollector) Collect(ch chan<- Metric) {
	if !hc.enabled {
		return
	}

	hc.mutex.RLock()
	defer hc.mutex.RUnlock()

	for _, metric := range hc.metrics {
		ch <- metric
	}
}

// RecordRequest 记录HTTP请求
func (hc *HTTPCollector) RecordRequest(method, statusCode string, duration time.Duration, requestSize, responseSize int64) {
	// 增加请求计数
	if counter, ok := hc.metrics["http_requests_total"].(Counter); ok {
		counter.Inc()
	}

	// 记录请求持续时间
	if histogram, ok := hc.metrics["http_request_duration_seconds"].(Histogram); ok {
		histogram.Observe(duration.Seconds())
	}

	// 记录请求大小
	if histogram, ok := hc.metrics["http_request_size_bytes"].(Histogram); ok {
		histogram.Observe(float64(requestSize))
	}

	// 记录响应大小
	if histogram, ok := hc.metrics["http_response_size_bytes"].(Histogram); ok {
		histogram.Observe(float64(responseSize))
	}
}

// IncInFlightRequests 增加正在处理的请求数
func (hc *HTTPCollector) IncInFlightRequests() {
	if gauge, ok := hc.metrics["http_requests_in_flight"].(Gauge); ok {
		gauge.Inc()
	}
}

// DecInFlightRequests 减少正在处理的请求数
func (hc *HTTPCollector) DecInFlightRequests() {
	if gauge, ok := hc.metrics["http_requests_in_flight"].(Gauge); ok {
		gauge.Dec()
	}
}

// Enable 启用收集器
func (hc *HTTPCollector) Enable() {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	hc.enabled = true
}

// Disable 禁用收集器
func (hc *HTTPCollector) Disable() {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	hc.enabled = false
}

// CollectorManager 收集器管理器
type CollectorManager struct {
	collectors map[string]Collector
	mutex      sync.RWMutex
	registry   Registry
	factory    Factory
}

// NewCollectorManager 创建收集器管理器
func NewCollectorManager(registry Registry, factory Factory) *CollectorManager {
	return &CollectorManager{
		collectors: make(map[string]Collector),
		registry:   registry,
		factory:    factory,
	}
}

// RegisterCollector 注册收集器
func (cm *CollectorManager) RegisterCollector(name string, collector Collector) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if _, exists := cm.collectors[name]; exists {
		return ErrCollectorExists
	}

	cm.collectors[name] = collector
	return nil
}

// UnregisterCollector 注销收集器
func (cm *CollectorManager) UnregisterCollector(name string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if _, exists := cm.collectors[name]; !exists {
		return ErrCollectorNotFound
	}

	delete(cm.collectors, name)
	return nil
}

// GetCollector 获取收集器
func (cm *CollectorManager) GetCollector(name string) (Collector, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	collector, exists := cm.collectors[name]
	return collector, exists
}

// GetAllCollectors 获取所有收集器
func (cm *CollectorManager) GetAllCollectors() map[string]Collector {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	result := make(map[string]Collector)
	for k, v := range cm.collectors {
		result[k] = v
	}
	return result
}

// CollectAll 收集所有指标
func (cm *CollectorManager) CollectAll() error {
	cm.mutex.RLock()
	collectors := make([]Collector, 0, len(cm.collectors))
	for _, collector := range cm.collectors {
		collectors = append(collectors, collector)
	}
	cm.mutex.RUnlock()

	// 收集指标
	metricChan := make(chan Metric, 100)
	go func() {
		defer close(metricChan)
		for _, collector := range collectors {
			collector.Collect(metricChan)
		}
	}()

	// 注册指标
	for metric := range metricChan {
		if err := cm.registry.Register(metric); err != nil {
			// 忽略已存在的指标错误
			if err != ErrMetricExists {
				return err
			}
		}
	}

	return nil
}

// StartPeriodicCollection 启动定期收集
func (cm *CollectorManager) StartPeriodicCollection(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := cm.CollectAll(); err != nil {
				fmt.Printf("Failed to collect metrics: %v\n", err)
			}
		}
	}
}