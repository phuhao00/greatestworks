package events

import (
	"sync"
	"time"
)

// EventType 事件类型
type EventType string

const (
	EventTypePlayerLogin   EventType = "player_login"
	EventTypePlayerLogout  EventType = "player_logout"
	EventTypePlayerMove    EventType = "player_move"
	EventTypePlayerAction  EventType = "player_action"
	EventTypePlayerChat    EventType = "player_chat"
	EventTypePlayerMail    EventType = "player_mail"
	EventTypeGameBattle    EventType = "game_battle"
	EventTypeGameShop      EventType = "game_shop"
	EventTypeGameBag       EventType = "game_bag"
	EventTypeGamePet       EventType = "game_pet"
	EventTypeGameBuilding  EventType = "game_building"
	EventTypeSystemError   EventType = "system_error"
	EventTypeSystemWarning EventType = "system_warning"
	EventTypeSystemInfo    EventType = "system_info"
	EventTypeSystemStart   EventType = "system_start"
	EventTypeSystemStop    EventType = "system_stop"
	EventTypeSystemHealth  EventType = "system_health"
)

// EventMetrics 事件指标
type EventMetrics struct {
	eventCounts     map[EventType]uint64
	successCounts   map[EventType]uint64
	errorCounts     map[EventType]uint64
	droppedCounts   map[EventType]uint64
	processingTimes map[EventType][]time.Duration
	mu              sync.RWMutex
}

// NewEventMetrics 创建事件指标
func NewEventMetrics() *EventMetrics {
	return &EventMetrics{
		eventCounts:     make(map[EventType]uint64),
		successCounts:   make(map[EventType]uint64),
		errorCounts:     make(map[EventType]uint64),
		droppedCounts:   make(map[EventType]uint64),
		processingTimes: make(map[EventType][]time.Duration),
	}
}

// IncrementEventCount 增加事件计数
func (em *EventMetrics) IncrementEventCount(eventType EventType) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.eventCounts[eventType]++
}

// IncrementSuccessCount 增加成功计数
func (em *EventMetrics) IncrementSuccessCount(eventType EventType) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.successCounts[eventType]++
}

// IncrementErrorCount 增加错误计数
func (em *EventMetrics) IncrementErrorCount(eventType EventType) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.errorCounts[eventType]++
}

// IncrementDroppedCount 增加丢弃计数
func (em *EventMetrics) IncrementDroppedCount(eventType EventType) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.droppedCounts[eventType]++
}

// RecordProcessingTime 记录处理时间
func (em *EventMetrics) RecordProcessingTime(eventType EventType, duration time.Duration) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if _, exists := em.processingTimes[eventType]; !exists {
		em.processingTimes[eventType] = make([]time.Duration, 0)
	}

	// 保留最近100个处理时间记录
	if len(em.processingTimes[eventType]) >= 100 {
		em.processingTimes[eventType] = em.processingTimes[eventType][1:]
	}
	em.processingTimes[eventType] = append(em.processingTimes[eventType], duration)
}

// GetEventCount 获取事件计数
func (em *EventMetrics) GetEventCount(eventType EventType) uint64 {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.eventCounts[eventType]
}

// GetSuccessCount 获取成功计数
func (em *EventMetrics) GetSuccessCount(eventType EventType) uint64 {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.successCounts[eventType]
}

// GetErrorCount 获取错误计数
func (em *EventMetrics) GetErrorCount(eventType EventType) uint64 {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.errorCounts[eventType]
}

// GetDroppedCount 获取丢弃计数
func (em *EventMetrics) GetDroppedCount(eventType EventType) uint64 {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.droppedCounts[eventType]
}

// GetAverageProcessingTime 获取平均处理时间
func (em *EventMetrics) GetAverageProcessingTime(eventType EventType) time.Duration {
	em.mu.RLock()
	defer em.mu.RUnlock()

	times, exists := em.processingTimes[eventType]
	if !exists || len(times) == 0 {
		return 0
	}

	var total time.Duration
	for _, t := range times {
		total += t
	}
	return total / time.Duration(len(times))
}

// GetSuccessRate 获取成功率
func (em *EventMetrics) GetSuccessRate(eventType EventType) float64 {
	em.mu.RLock()
	defer em.mu.RUnlock()

	total := em.eventCounts[eventType]
	if total == 0 {
		return 0
	}

	success := em.successCounts[eventType]
	return float64(success) / float64(total)
}

// GetAllMetrics 获取所有指标
func (em *EventMetrics) GetAllMetrics() map[string]interface{} {
	em.mu.RLock()
	defer em.mu.RUnlock()

	metrics := map[string]interface{}{
		"event_counts":   make(map[string]uint64),
		"success_counts": make(map[string]uint64),
		"error_counts":   make(map[string]uint64),
		"dropped_counts": make(map[string]uint64),
		"success_rates":  make(map[string]float64),
		"avg_times":      make(map[string]string),
	}

	// 收集所有事件类型
	allEventTypes := make(map[EventType]bool)
	for eventType := range em.eventCounts {
		allEventTypes[eventType] = true
	}
	for eventType := range em.successCounts {
		allEventTypes[eventType] = true
	}
	for eventType := range em.errorCounts {
		allEventTypes[eventType] = true
	}
	for eventType := range em.droppedCounts {
		allEventTypes[eventType] = true
	}

	// 为每个事件类型生成指标
	for eventType := range allEventTypes {
		eventTypeStr := string(eventType)
		metrics["event_counts"].(map[string]uint64)[eventTypeStr] = em.eventCounts[eventType]
		metrics["success_counts"].(map[string]uint64)[eventTypeStr] = em.successCounts[eventType]
		metrics["error_counts"].(map[string]uint64)[eventTypeStr] = em.errorCounts[eventType]
		metrics["dropped_counts"].(map[string]uint64)[eventTypeStr] = em.droppedCounts[eventType]

		// 计算成功率
		total := em.eventCounts[eventType]
		if total > 0 {
			success := em.successCounts[eventType]
			metrics["success_rates"].(map[string]float64)[eventTypeStr] = float64(success) / float64(total)
		} else {
			metrics["success_rates"].(map[string]float64)[eventTypeStr] = 0
		}

		// 计算平均处理时间
		times, exists := em.processingTimes[eventType]
		if exists && len(times) > 0 {
			var total time.Duration
			for _, t := range times {
				total += t
			}
			avg := total / time.Duration(len(times))
			metrics["avg_times"].(map[string]string)[eventTypeStr] = avg.String()
		} else {
			metrics["avg_times"].(map[string]string)[eventTypeStr] = "0s"
		}
	}

	return metrics
}

// Reset 重置指标
func (em *EventMetrics) Reset() {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.eventCounts = make(map[EventType]uint64)
	em.successCounts = make(map[EventType]uint64)
	em.errorCounts = make(map[EventType]uint64)
	em.droppedCounts = make(map[EventType]uint64)
	em.processingTimes = make(map[EventType][]time.Duration)
}
