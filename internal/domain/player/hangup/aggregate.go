package hangup

import (
	"time"
)

// HangupAggregate 挂机聚合根
type HangupAggregate struct {
	playerID        string
	currentLocation *HangupLocation
	offlineReward   *OfflineReward
	hangupStatus    HangupStatus
	efficiencyBonus *EfficiencyBonus
	lastOnlineTime  time.Time
	lastOfflineTime time.Time
	totalHangupTime time.Duration
	dailyHangupTime time.Duration
	lastResetDate   time.Time
	updatedAt       time.Time
	version         int
}

// NewHangupAggregate 创建挂机聚合根
func NewHangupAggregate(playerID string) *HangupAggregate {
	now := time.Now()
	return &HangupAggregate{
		playerID:        playerID,
		currentLocation: nil,
		offlineReward:   NewOfflineReward(),
		hangupStatus:    HangupStatusOffline,
		efficiencyBonus: NewEfficiencyBonus(),
		lastOnlineTime:  now,
		lastOfflineTime: now,
		totalHangupTime: 0,
		dailyHangupTime: 0,
		lastResetDate:   now.Truncate(24 * time.Hour),
		updatedAt:       now,
		version:         1,
	}
}

// GetPlayerID 获取玩家ID
func (h *HangupAggregate) GetPlayerID() string {
	return h.playerID
}

// SetHangupLocation 设置挂机地点
func (h *HangupAggregate) SetHangupLocation(location *HangupLocation) error {
	if location == nil {
		return ErrInvalidHangupLocation
	}

	// 检查地点解锁条件
	if !location.IsUnlocked() {
		return ErrHangupLocationNotUnlocked
	}

	// 检查玩家等级要求
	if !h.checkLocationRequirements(location) {
		return ErrHangupLocationRequirementNotMet
	}

	h.currentLocation = location
	h.updateVersion()
	return nil
}

// GetCurrentLocation 获取当前挂机地点
func (h *HangupAggregate) GetCurrentLocation() *HangupLocation {
	return h.currentLocation
}

// StartHangup 开始挂机
func (h *HangupAggregate) StartHangup() error {
	if h.currentLocation == nil {
		return ErrNoHangupLocationSet
	}

	if h.hangupStatus == HangupStatusOnline {
		return ErrAlreadyHangingUp
	}

	h.hangupStatus = HangupStatusOnline
	h.lastOnlineTime = time.Now()
	h.updateVersion()
	return nil
}

// StopHangup 停止挂机
func (h *HangupAggregate) StopHangup() error {
	if h.hangupStatus == HangupStatusOffline {
		return ErrNotHangingUp
	}

	h.hangupStatus = HangupStatusOffline
	h.lastOfflineTime = time.Now()

	// 计算挂机时间
	hangupDuration := h.lastOfflineTime.Sub(h.lastOnlineTime)
	h.totalHangupTime += hangupDuration
	h.dailyHangupTime += hangupDuration

	h.updateVersion()
	return nil
}

// CalculateOfflineReward 计算离线奖励
func (h *HangupAggregate) CalculateOfflineReward(offlineDuration time.Duration) (*OfflineReward, error) {
	if h.currentLocation == nil {
		return nil, ErrNoHangupLocationSet
	}

	// 限制最大离线时间（例如24小时）
	maxOfflineTime := 24 * time.Hour
	if offlineDuration > maxOfflineTime {
		offlineDuration = maxOfflineTime
	}

	// 计算基础奖励
	baseReward := h.currentLocation.CalculateBaseReward(offlineDuration)

	// 应用效率加成
	finalReward := h.efficiencyBonus.ApplyBonus(baseReward)

	// 创建离线奖励
	offlineReward := &OfflineReward{
		Experience:      finalReward.Experience,
		Gold:            finalReward.Gold,
		Items:           finalReward.Items,
		OfflineDuration: offlineDuration,
		LocationID:      h.currentLocation.GetID(),
		CalculatedAt:    time.Now(),
	}

	h.offlineReward = offlineReward
	h.updateVersion()
	return offlineReward, nil
}

// ClaimOfflineReward 领取离线奖励
func (h *HangupAggregate) ClaimOfflineReward() (*OfflineReward, error) {
	if h.offlineReward == nil {
		return nil, ErrNoOfflineRewardAvailable
	}

	if h.offlineReward.IsClaimed {
		return nil, ErrOfflineRewardAlreadyClaimed
	}

	// 标记为已领取
	h.offlineReward.IsClaimed = true
	h.offlineReward.ClaimedAt = time.Now()

	reward := h.offlineReward
	h.offlineReward = nil // 清空已领取的奖励

	h.updateVersion()
	return reward, nil
}

// GetOfflineReward 获取离线奖励
func (h *HangupAggregate) GetOfflineReward() *OfflineReward {
	return h.offlineReward
}

// UpdateEfficiencyBonus 更新效率加成
func (h *HangupAggregate) UpdateEfficiencyBonus(bonus *EfficiencyBonus) {
	h.efficiencyBonus = bonus
	h.updateVersion()
}

// GetEfficiencyBonus 获取效率加成
func (h *HangupAggregate) GetEfficiencyBonus() *EfficiencyBonus {
	return h.efficiencyBonus
}

// GetHangupStatus 获取挂机状态
func (h *HangupAggregate) GetHangupStatus() HangupStatus {
	return h.hangupStatus
}

// GetTotalHangupTime 获取总挂机时间
func (h *HangupAggregate) GetTotalHangupTime() time.Duration {
	return h.totalHangupTime
}

// GetDailyHangupTime 获取每日挂机时间
func (h *HangupAggregate) GetDailyHangupTime() time.Duration {
	// 检查是否需要重置每日时间
	h.checkDailyReset()
	return h.dailyHangupTime
}

// GetLastOnlineTime 获取最后在线时间
func (h *HangupAggregate) GetLastOnlineTime() time.Time {
	return h.lastOnlineTime
}

// GetLastOfflineTime 获取最后离线时间
func (h *HangupAggregate) GetLastOfflineTime() time.Time {
	return h.lastOfflineTime
}

// IsOnline 是否在线挂机
func (h *HangupAggregate) IsOnline() bool {
	return h.hangupStatus == HangupStatusOnline
}

// IsOffline 是否离线
func (h *HangupAggregate) IsOffline() bool {
	return h.hangupStatus == HangupStatusOffline
}

// GetCurrentOfflineDuration 获取当前离线时长
func (h *HangupAggregate) GetCurrentOfflineDuration() time.Duration {
	if h.IsOnline() {
		return 0
	}
	return time.Since(h.lastOfflineTime)
}

// GetVersion 获取版本
func (h *HangupAggregate) GetVersion() int {
	return h.version
}

// GetUpdatedAt 获取更新时间
func (h *HangupAggregate) GetUpdatedAt() time.Time {
	return h.updatedAt
}

// GetID 获取挂机ID (暂时使用playerID作为ID)
func (h *HangupAggregate) GetID() string {
	return h.playerID
}

// GetLocationID 获取地点ID
func (h *HangupAggregate) GetLocationID() string {
	if h.currentLocation == nil {
		return ""
	}
	return h.currentLocation.GetID()
}

// GetStartTime 获取开始时间
func (h *HangupAggregate) GetStartTime() time.Time {
	return h.lastOnlineTime
}

// GetEndTime 获取结束时间
func (h *HangupAggregate) GetEndTime() time.Time {
	return h.lastOfflineTime
}

// GetDuration 获取持续时间
func (h *HangupAggregate) GetDuration() time.Duration {
	return h.totalHangupTime
}

// GetEfficiency 获取效率
func (h *HangupAggregate) GetEfficiency() float64 {
	if h.efficiencyBonus == nil {
		return 1.0
	}
	return h.efficiencyBonus.GetTotalBonus()
}

// GetBaseRate 获取基础速率 (经验速率)
func (h *HangupAggregate) GetBaseRate() float64 {
	if h.currentLocation == nil {
		return 0.0
	}
	return h.currentLocation.GetBaseExpRate()
}

// GetStatus 获取状态
func (h *HangupAggregate) GetStatus() HangupStatus {
	return h.hangupStatus
}

// GetRewards 获取奖励列表
func (h *HangupAggregate) GetRewards() []RewardItem {
	if h.offlineReward == nil {
		return []RewardItem{}
	}
	return h.offlineReward.Items
}

// GetCreatedAt 获取创建时间
func (h *HangupAggregate) GetCreatedAt() time.Time {
	return h.lastOnlineTime // 使用第一次在线时间作为创建时间
}

// ReconstructHangupAggregate 重构挂机聚合根（用于从持久化数据恢复）
func ReconstructHangupAggregate(
	hangupID string,
	playerID string,
	locationID string,
	startTime time.Time,
	endTime time.Time,
	duration time.Duration,
	efficiency float64,
	baseRate float64,
	status HangupStatus,
	rewards []RewardItem,
	createdAt time.Time,
	updatedAt time.Time,
) *HangupAggregate {
	// 创建基础聚合根
	aggregate := &HangupAggregate{
		playerID:        playerID,
		currentLocation: nil, // 需要根据locationID重新加载
		offlineReward:   nil,
		hangupStatus:    status,
		efficiencyBonus: NewEfficiencyBonus(),
		lastOnlineTime:  startTime,
		lastOfflineTime: endTime,
		totalHangupTime: duration,
		dailyHangupTime: duration,
		lastResetDate:   createdAt.Truncate(24 * time.Hour),
		updatedAt:       updatedAt,
		version:         1,
	}

	// 如果有奖励，创建离线奖励对象
	if len(rewards) > 0 {
		aggregate.offlineReward = &OfflineReward{
			Items:           rewards,
			OfflineDuration: duration,
			LocationID:      locationID,
			CalculatedAt:    updatedAt,
			IsClaimed:       false,
		}
	}

	return aggregate
}

// 私有方法

// checkLocationRequirements 检查地点要求
func (h *HangupAggregate) checkLocationRequirements(location *HangupLocation) bool {
	// 这里需要外部提供玩家等级信息
	// 暂时返回true，实际实现中需要检查玩家等级、任务完成情况等
	return true
}

// checkDailyReset 检查每日重置
func (h *HangupAggregate) checkDailyReset() {
	now := time.Now()
	currentDate := now.Truncate(24 * time.Hour)

	if currentDate.After(h.lastResetDate) {
		h.dailyHangupTime = 0
		h.lastResetDate = currentDate
		h.updateVersion()
	}
}

// updateVersion 更新版本
func (h *HangupAggregate) updateVersion() {
	h.version++
	h.updatedAt = time.Now()
}

// HangupStatus 挂机状态
type HangupStatus int

const (
	HangupStatusOffline HangupStatus = iota // 离线
	HangupStatusOnline                      // 在线挂机
	HangupStatusPaused                      // 暂停
)

// String 返回挂机状态的字符串表示
func (hs HangupStatus) String() string {
	switch hs {
	case HangupStatusOffline:
		return "offline"
	case HangupStatusOnline:
		return "online"
	case HangupStatusPaused:
		return "paused"
	default:
		return "unknown"
	}
}
