package replication

import "time"

// PlayerJoinedEvent 玩家加入事件
type PlayerJoinedEvent struct {
	InstanceID string
	PlayerID   string
	PlayerName string
	Timestamp  time.Time
}

// PlayerLeftEvent 玩家离开事件
type PlayerLeftEvent struct {
	InstanceID string
	PlayerID   string
	Timestamp  time.Time
}

// InstanceStartedEvent 实例启动事件
type InstanceStartedEvent struct {
	InstanceID  string
	PlayerCount int
	Timestamp   time.Time
}

// InstanceFullEvent 实例满员事件
type InstanceFullEvent struct {
	InstanceID string
	Timestamp  time.Time
}

// InstanceProgressUpdatedEvent 进度更新事件
type InstanceProgressUpdatedEvent struct {
	InstanceID string
	Progress   int
	Task       string
	Timestamp  time.Time
}

// InstanceCompletedEvent 实例完成事件
type InstanceCompletedEvent struct {
	InstanceID string
	Duration   time.Duration
	Timestamp  time.Time
}

// InstanceClosingEvent 实例关闭中事件
type InstanceClosingEvent struct {
	InstanceID string
	Timestamp  time.Time
}

// InstanceClosedEvent 实例已关闭事件
type InstanceClosedEvent struct {
	InstanceID string
	Timestamp  time.Time
}
