// Package replication 副本/实例领域模型
// 负责管理游戏副本（Dungeon）和实例（Instance）的生命周期
package replication

import (
	"fmt"
	"sync"
	"time"
)

// InstanceStatus 实例状态
type InstanceStatus int

const (
	InstanceStatusPending  InstanceStatus = iota // 等待创建
	InstanceStatusCreating                       // 创建中
	InstanceStatusActive                         // 活跃中
	InstanceStatusFull                           // 已满员
	InstanceStatusClosing                        // 关闭中
	InstanceStatusClosed                         // 已关闭
)

// InstanceType 实例类型
type InstanceType int

const (
	InstanceTypeDungeon InstanceType = iota // 副本
	InstanceTypeRaid                        // 团队副本
	InstanceTypePVP                         // PVP竞技场
	InstanceTypeEvent                       // 活动副本
)

// Instance 副本实例聚合根
type Instance struct {
	mu sync.RWMutex

	// 标识
	instanceID   string       // 实例ID
	templateID   string       // 模板ID（副本配置ID）
	instanceType InstanceType // 实例类型
	sceneID      string       // 关联的场景ID

	// 玩家
	players       map[string]*PlayerInfo // 玩家列表
	maxPlayers    int                    // 最大玩家数
	minPlayers    int                    // 最小玩家数
	ownerPlayerID string                 // 创建者/队长

	// 状态
	status     InstanceStatus
	difficulty int // 难度等级

	// 时间
	createdAt time.Time
	startedAt time.Time
	expireAt  time.Time
	closedAt  time.Time
	lifetime  time.Duration // 生命周期

	// 进度
	progress        int               // 进度百分比
	completedTasks  []string          // 已完成任务
	metadata        map[string]string // 元数据
	scoreMultiplier float64           // 分数倍率

	// 领域事件
	events []interface{}
}

// PlayerInfo 玩家信息
type PlayerInfo struct {
	PlayerID   string
	PlayerName string
	Level      int
	JoinedAt   time.Time
	IsReady    bool
	Role       string // tank, healer, dps
}

// NewInstance 创建新实例
func NewInstance(
	instanceID string,
	templateID string,
	instanceType InstanceType,
	ownerPlayerID string,
	maxPlayers int,
	lifetime time.Duration,
) *Instance {
	now := time.Now()
	return &Instance{
		instanceID:      instanceID,
		templateID:      templateID,
		instanceType:    instanceType,
		ownerPlayerID:   ownerPlayerID,
		maxPlayers:      maxPlayers,
		minPlayers:      1,
		players:         make(map[string]*PlayerInfo),
		status:          InstanceStatusPending,
		difficulty:      1,
		createdAt:       now,
		expireAt:        now.Add(lifetime),
		lifetime:        lifetime,
		progress:        0,
		completedTasks:  []string{},
		metadata:        make(map[string]string),
		scoreMultiplier: 1.0,
		events:          []interface{}{},
	}
}

// ID 获取实例ID
func (i *Instance) ID() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.instanceID
}

// TemplateID 获取模板ID
func (i *Instance) TemplateID() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.templateID
}

// Type 获取实例类型
func (i *Instance) Type() InstanceType {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.instanceType
}

// Status 获取状态
func (i *Instance) Status() InstanceStatus {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.status
}

// PlayerCount 获取当前玩家数
func (i *Instance) PlayerCount() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return len(i.players)
}

// MaxPlayers 获取最大玩家数
func (i *Instance) MaxPlayers() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.maxPlayers
}

// Progress 获取进度
func (i *Instance) Progress() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.progress
}

// SceneID 获取场景ID
func (i *Instance) SceneID() string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.sceneID
}

// SetSceneID 设置场景ID
func (i *Instance) SetSceneID(sceneID string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.sceneID = sceneID
}

// AddPlayer 添加玩家
func (i *Instance) AddPlayer(playerID, playerName string, level int, role string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// 检查状态
	if i.status == InstanceStatusClosed || i.status == InstanceStatusClosing {
		return fmt.Errorf("instance is closed or closing")
	}

	// 检查是否已满
	if len(i.players) >= i.maxPlayers {
		return fmt.Errorf("instance is full")
	}

	// 检查是否已存在
	if _, exists := i.players[playerID]; exists {
		return fmt.Errorf("player already in instance")
	}

	// 添加玩家
	i.players[playerID] = &PlayerInfo{
		PlayerID:   playerID,
		PlayerName: playerName,
		Level:      level,
		JoinedAt:   time.Now(),
		IsReady:    false,
		Role:       role,
	}

	// 发布事件
	i.addEvent(&PlayerJoinedEvent{
		InstanceID: i.instanceID,
		PlayerID:   playerID,
		PlayerName: playerName,
		Timestamp:  time.Now(),
	})

	// 检查是否满员
	if len(i.players) >= i.maxPlayers {
		i.status = InstanceStatusFull
		i.addEvent(&InstanceFullEvent{
			InstanceID: i.instanceID,
			Timestamp:  time.Now(),
		})
	}

	return nil
}

// RemovePlayer 移除玩家
func (i *Instance) RemovePlayer(playerID string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exists := i.players[playerID]; !exists {
		return fmt.Errorf("player not in instance")
	}

	delete(i.players, playerID)

	// 发布事件
	i.addEvent(&PlayerLeftEvent{
		InstanceID: i.instanceID,
		PlayerID:   playerID,
		Timestamp:  time.Now(),
	})

	// 更新状态
	if i.status == InstanceStatusFull && len(i.players) < i.maxPlayers {
		i.status = InstanceStatusActive
	}

	// 如果没有玩家了，标记为关闭中
	if len(i.players) == 0 {
		i.MarkForClosing()
	}

	return nil
}

// Start 启动实例
func (i *Instance) Start() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.status != InstanceStatusPending && i.status != InstanceStatusCreating {
		return fmt.Errorf("instance cannot be started in current status: %d", i.status)
	}

	// 检查最小玩家数
	if len(i.players) < i.minPlayers {
		return fmt.Errorf("not enough players: %d/%d", len(i.players), i.minPlayers)
	}

	i.status = InstanceStatusActive
	i.startedAt = time.Now()

	// 发布事件
	i.addEvent(&InstanceStartedEvent{
		InstanceID:  i.instanceID,
		PlayerCount: len(i.players),
		Timestamp:   time.Now(),
	})

	return nil
}

// UpdateProgress 更新进度
func (i *Instance) UpdateProgress(progress int, completedTask string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.progress = progress
	if completedTask != "" {
		i.completedTasks = append(i.completedTasks, completedTask)
	}

	// 发布事件
	i.addEvent(&InstanceProgressUpdatedEvent{
		InstanceID: i.instanceID,
		Progress:   progress,
		Task:       completedTask,
		Timestamp:  time.Now(),
	})

	// 如果完成了
	if progress >= 100 {
		i.status = InstanceStatusClosing
		i.addEvent(&InstanceCompletedEvent{
			InstanceID: i.instanceID,
			Duration:   time.Since(i.startedAt),
			Timestamp:  time.Now(),
		})
	}
}

// MarkForClosing 标记为关闭中
func (i *Instance) MarkForClosing() {
	if i.status != InstanceStatusClosed {
		i.status = InstanceStatusClosing
		i.addEvent(&InstanceClosingEvent{
			InstanceID: i.instanceID,
			Timestamp:  time.Now(),
		})
	}
}

// Close 关闭实例
func (i *Instance) Close() {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.status = InstanceStatusClosed
	i.closedAt = time.Now()

	// 发布事件
	i.addEvent(&InstanceClosedEvent{
		InstanceID: i.instanceID,
		Timestamp:  time.Now(),
	})
}

// IsExpired 检查是否过期
func (i *Instance) IsExpired() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return time.Now().After(i.expireAt)
}

// GetPlayers 获取玩家列表
func (i *Instance) GetPlayers() []*PlayerInfo {
	i.mu.RLock()
	defer i.mu.RUnlock()

	players := make([]*PlayerInfo, 0, len(i.players))
	for _, p := range i.players {
		players = append(players, p)
	}
	return players
}

// addEvent 添加领域事件
func (i *Instance) addEvent(event interface{}) {
	i.events = append(i.events, event)
}

// GetEvents 获取并清空事件
func (i *Instance) GetEvents() []interface{} {
	i.mu.Lock()
	defer i.mu.Unlock()

	events := i.events
	i.events = []interface{}{}
	return events
}

// GetMetadata 获取元数据
func (i *Instance) GetMetadata(key string) (string, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	val, ok := i.metadata[key]
	return val, ok
}

// SetMetadata 设置元数据
func (i *Instance) SetMetadata(key, value string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.metadata[key] = value
}

// CreatedAt 获取创建时间
func (i *Instance) CreatedAt() time.Time {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.createdAt
}

// Difficulty 获取难度
func (i *Instance) Difficulty() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.difficulty
}
