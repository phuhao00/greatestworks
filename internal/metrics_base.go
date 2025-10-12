package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// MetricType 指标类型
type MetricType int

const (
	MetricTypeCounter   MetricType = iota // 计数器
	MetricTypeGauge                       // 仪表盘
	MetricTypeHistogram                   // 直方图
	MetricTypeSummary                     // 摘要
)

// MetricValue 指标值
type MetricValue struct {
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// MetricsConfig 指标配置
type MetricsConfig struct {
	Enabled         bool          `json:"enabled"`
	CollectInterval time.Duration `json:"collect_interval"`
	RetentionTime   time.Duration `json:"retention_time"`
	MaxSamples      int           `json:"max_samples"`
}

// MetricsBase 指标基类
type MetricsBase struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        MetricType        `json:"type"`
	Labels      map[string]string `json:"labels"`
	Values      []MetricValue     `json:"values"`
	Config      *MetricsConfig    `json:"config"`
	mu          sync.RWMutex
	counter     int64
	gauge       int64
	lastUpdate  time.Time
	active      bool
}

// NewMetricsBase 创建指标基类实例
func NewMetricsBase(name, description string, metricType MetricType, config *MetricsConfig) *MetricsBase {
	return &MetricsBase{
		Name:        name,
		Description: description,
		Type:        metricType,
		Labels:      make(map[string]string),
		Values:      make([]MetricValue, 0),
		Config:      config,
		lastUpdate:  time.Now(),
		active:      true,
	}
}

// GetDescription 获取指标描述
func (m *MetricsBase) GetDescription() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Description
}

// SetDescription 设置指标描述
func (m *MetricsBase) SetDescription(desc string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Description = desc
}

// GetName 获取指标名称
func (m *MetricsBase) GetName() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Name
}

// SetName 设置指标名称
func (m *MetricsBase) SetName(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Name = name
}

// AddLabel 添加标签
func (m *MetricsBase) AddLabel(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Labels[key] = value
}

// RemoveLabel 移除标签
func (m *MetricsBase) RemoveLabel(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Labels, key)
}

// GetLabels 获取所有标签
func (m *MetricsBase) GetLabels() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	labels := make(map[string]string)
	for k, v := range m.Labels {
		labels[k] = v
	}
	return labels
}

// Inc 计数器递增
func (m *MetricsBase) Inc() {
	if m.Type != MetricTypeCounter {
		return
	}
	atomic.AddInt64(&m.counter, 1)
	m.recordValue(float64(atomic.LoadInt64(&m.counter)))
}

// Add 计数器增加指定值
func (m *MetricsBase) Add(value float64) {
	if m.Type != MetricTypeCounter {
		return
	}
	atomic.AddInt64(&m.counter, int64(value))
	m.recordValue(float64(atomic.LoadInt64(&m.counter)))
}

// Set 设置仪表盘值
func (m *MetricsBase) Set(value float64) {
	if m.Type != MetricTypeGauge {
		return
	}
	atomic.StoreInt64(&m.gauge, int64(value))
	m.recordValue(value)
}

// Get 获取当前值
func (m *MetricsBase) Get() float64 {
	switch m.Type {
	case MetricTypeCounter:
		return float64(atomic.LoadInt64(&m.counter))
	case MetricTypeGauge:
		return float64(atomic.LoadInt64(&m.gauge))
	default:
		return 0
	}
}

// recordValue 记录指标值
func (m *MetricsBase) recordValue(value float64) {
	if !m.Config.Enabled || !m.active {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 创建新的指标值
	metricValue := MetricValue{
		Value:     value,
		Timestamp: time.Now(),
		Labels:    m.getLabelsSnapshot(),
	}

	// 添加到值列表
	m.Values = append(m.Values, metricValue)
	m.lastUpdate = time.Now()

	// 清理过期数据
	m.cleanupExpiredValues()

	// 限制最大样本数
	if len(m.Values) > m.Config.MaxSamples {
		m.Values = m.Values[len(m.Values)-m.Config.MaxSamples:]
	}

	log.Printf("[MetricsBase] 记录指标 %s: %f", m.Name, value)
}

// getLabelsSnapshot 获取标签快照
func (m *MetricsBase) getLabelsSnapshot() map[string]string {
	labels := make(map[string]string)
	for k, v := range m.Labels {
		labels[k] = v
	}
	return labels
}

// cleanupExpiredValues 清理过期值
func (m *MetricsBase) cleanupExpiredValues() {
	if m.Config.RetentionTime <= 0 {
		return
	}

	cutoff := time.Now().Add(-m.Config.RetentionTime)
	validValues := make([]MetricValue, 0)

	for _, value := range m.Values {
		if value.Timestamp.After(cutoff) {
			validValues = append(validValues, value)
		}
	}

	m.Values = validValues
}

// GetValues 获取所有指标值
func (m *MetricsBase) GetValues() []MetricValue {
	m.mu.RLock()
	defer m.mu.RUnlock()

	values := make([]MetricValue, len(m.Values))
	copy(values, m.Values)
	return values
}

// GetLatestValue 获取最新指标值
func (m *MetricsBase) GetLatestValue() *MetricValue {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.Values) == 0 {
		return nil
	}

	latest := m.Values[len(m.Values)-1]
	return &latest
}

// GetValuesSince 获取指定时间以来的指标值
func (m *MetricsBase) GetValuesSince(since time.Time) []MetricValue {
	m.mu.RLock()
	defer m.mu.RUnlock()

	values := make([]MetricValue, 0)
	for _, value := range m.Values {
		if value.Timestamp.After(since) {
			values = append(values, value)
		}
	}
	return values
}

// GetStatistics 获取统计信息
func (m *MetricsBase) GetStatistics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.Values) == 0 {
		return map[string]interface{}{
			"count": 0,
			"min":   0,
			"max":   0,
			"avg":   0,
			"sum":   0,
		}
	}

	var min, max, sum float64
	min = m.Values[0].Value
	max = m.Values[0].Value

	for _, value := range m.Values {
		if value.Value < min {
			min = value.Value
		}
		if value.Value > max {
			max = value.Value
		}
		sum += value.Value
	}

	avg := sum / float64(len(m.Values))

	return map[string]interface{}{
		"count": len(m.Values),
		"min":   min,
		"max":   max,
		"avg":   avg,
		"sum":   sum,
	}
}

// Reset 重置指标
func (m *MetricsBase) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	atomic.StoreInt64(&m.counter, 0)
	atomic.StoreInt64(&m.gauge, 0)
	m.Values = make([]MetricValue, 0)
	m.lastUpdate = time.Now()

	log.Printf("[MetricsBase] 重置指标 %s", m.Name)
}

// Enable 启用指标收集
func (m *MetricsBase) Enable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.active = true
	log.Printf("[MetricsBase] 启用指标 %s", m.Name)
}

// Disable 禁用指标收集
func (m *MetricsBase) Disable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.active = false
	log.Printf("[MetricsBase] 禁用指标 %s", m.Name)
}

// IsActive 检查指标是否活跃
func (m *MetricsBase) IsActive() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.active
}

// GetLastUpdate 获取最后更新时间
func (m *MetricsBase) GetLastUpdate() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastUpdate
}

// ToJSON 转换为JSON格式
func (m *MetricsBase) ToJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data := map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"type":        m.Type,
		"labels":      m.Labels,
		"current":     m.Get(),
		"statistics":  m.GetStatistics(),
		"last_update": m.lastUpdate,
		"active":      m.active,
	}

	return json.Marshal(data)
}

// String 字符串表示
func (m *MetricsBase) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return fmt.Sprintf("Metric{name=%s, type=%d, value=%f, active=%t}",
		m.Name, m.Type, m.Get(), m.active)
}
