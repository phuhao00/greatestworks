package hangup

import (
	"math/rand"
	"time"
)

// LocationType 挂机地点类型
type LocationType int

const (
	LocationTypeUnknown LocationType = iota
	LocationTypeForest                      // 森林
	LocationTypeMountain                    // 山脉
	LocationTypeDesert                      // 沙漠
	LocationTypeOcean                       // 海洋
	LocationTypeCave                        // 洞穴
	LocationTypeDungeon                     // 地牢
	LocationTypeCity                        // 城市
	LocationTypeVillage                     // 村庄
	LocationTypeSpecial                     // 特殊地点
)

// String 返回地点类型的字符串表示
func (lt LocationType) String() string {
	switch lt {
	case LocationTypeForest:
		return "forest"
	case LocationTypeMountain:
		return "mountain"
	case LocationTypeDesert:
		return "desert"
	case LocationTypeOcean:
		return "ocean"
	case LocationTypeCave:
		return "cave"
	case LocationTypeDungeon:
		return "dungeon"
	case LocationTypeCity:
		return "city"
	case LocationTypeVillage:
		return "village"
	case LocationTypeSpecial:
		return "special"
	default:
		return "unknown"
	}
}

// GetExpMultiplier 获取经验倍率
func (lt LocationType) GetExpMultiplier() float64 {
	switch lt {
	case LocationTypeForest:
		return 1.0
	case LocationTypeMountain:
		return 1.2
	case LocationTypeDesert:
		return 1.1
	case LocationTypeOcean:
		return 1.3
	case LocationTypeCave:
		return 1.4
	case LocationTypeDungeon:
		return 1.5
	case LocationTypeCity:
		return 0.8
	case LocationTypeVillage:
		return 0.9
	case LocationTypeSpecial:
		return 2.0
	default:
		return 1.0
	}
}

// GetGoldMultiplier 获取金币倍率
func (lt LocationType) GetGoldMultiplier() float64 {
	switch lt {
	case LocationTypeForest:
		return 1.0
	case LocationTypeMountain:
		return 1.1
	case LocationTypeDesert:
		return 1.2
	case LocationTypeOcean:
		return 1.0
	case LocationTypeCave:
		return 1.3
	case LocationTypeDungeon:
		return 1.4
	case LocationTypeCity:
		return 1.5
	case LocationTypeVillage:
		return 1.2
	case LocationTypeSpecial:
		return 1.8
	default:
		return 1.0
	}
}

// RewardItem 奖励物品值对象
type RewardItem struct {
	Type     string `json:"type"`
	ItemID   string `json:"item_id"`
	Quantity int64  `json:"quantity"`
	Quality  string `json:"quality"`
}

// NewRewardItem 创建奖励物品
func NewRewardItem(itemType, itemID string, quantity int64, quality string) RewardItem {
	return RewardItem{
		Type:     itemType,
		ItemID:   itemID,
		Quantity: quantity,
		Quality:  quality,
	}
}

// IsValid 检查奖励物品是否有效
func (ri RewardItem) IsValid() bool {
	return ri.ItemID != "" && ri.Quantity > 0
}

// BaseReward 基础奖励值对象
type BaseReward struct {
	Experience int64        `json:"experience"`
	Gold       int64        `json:"gold"`
	Items      []RewardItem `json:"items"`
}

// NewBaseReward 创建基础奖励
func NewBaseReward(experience, gold int64) *BaseReward {
	return &BaseReward{
		Experience: experience,
		Gold:       gold,
		Items:      make([]RewardItem, 0),
	}
}

// AddItem 添加物品奖励
func (br *BaseReward) AddItem(item RewardItem) {
	if item.IsValid() {
		br.Items = append(br.Items, item)
	}
}

// IsEmpty 检查奖励是否为空
func (br *BaseReward) IsEmpty() bool {
	return br.Experience == 0 && br.Gold == 0 && len(br.Items) == 0
}

// Multiply 乘以倍率
func (br *BaseReward) Multiply(multiplier float64) *BaseReward {
	return &BaseReward{
		Experience: int64(float64(br.Experience) * multiplier),
		Gold:       int64(float64(br.Gold) * multiplier),
		Items:      br.Items, // 物品数量不变
	}
}

// Add 添加另一个奖励
func (br *BaseReward) Add(other *BaseReward) *BaseReward {
	result := &BaseReward{
		Experience: br.Experience + other.Experience,
		Gold:       br.Gold + other.Gold,
		Items:      make([]RewardItem, 0),
	}
	
	// 合并物品
	result.Items = append(result.Items, br.Items...)
	result.Items = append(result.Items, other.Items...)
	
	return result
}

// ItemDrop 物品掉落值对象
type ItemDrop struct {
	ItemID      string  `json:"item_id"`
	DropRate    float64 `json:"drop_rate"`    // 掉落率 (0.0-1.0)
	MinQuantity int     `json:"min_quantity"` // 最小数量
	MaxQuantity int     `json:"max_quantity"` // 最大数量
	HourlyRate  float64 `json:"hourly_rate"`  // 每小时掉落率
}

// NewItemDrop 创建物品掉落
func NewItemDrop(itemID string, dropRate float64, minQty, maxQty int) ItemDrop {
	return ItemDrop{
		ItemID:      itemID,
		DropRate:    dropRate,
		MinQuantity: minQty,
		MaxQuantity: maxQty,
		HourlyRate:  dropRate, // 默认每小时掉落率等于基础掉落率
	}
}

// SetHourlyRate 设置每小时掉落率
func (id *ItemDrop) SetHourlyRate(rate float64) {
	id.HourlyRate = rate
}

// ShouldDrop 判断是否应该掉落
func (id ItemDrop) ShouldDrop(hours float64) bool {
	// 计算在指定小时数内的掉落概率
	totalRate := id.HourlyRate * hours
	if totalRate >= 1.0 {
		return true // 100%掉落
	}
	
	// 随机判断
	return rand.Float64() < totalRate
}

// CalculateQuantity 计算掉落数量
func (id ItemDrop) CalculateQuantity(hours float64) int {
	if id.MinQuantity == id.MaxQuantity {
		return id.MinQuantity
	}
	
	// 基于小时数调整数量
	baseQuantity := id.MinQuantity + rand.Intn(id.MaxQuantity-id.MinQuantity+1)
	
	// 长时间挂机可能获得更多物品
	if hours > 12 {
		bonusChance := (hours - 12) / 12 // 每12小时增加一次奖励机会
		if rand.Float64() < bonusChance {
			baseQuantity += rand.Intn(id.MaxQuantity-id.MinQuantity+1)
		}
	}
	
	return baseQuantity
}

// IsValid 检查物品掉落配置是否有效
func (id ItemDrop) IsValid() bool {
	return id.ItemID != "" && id.DropRate >= 0 && id.DropRate <= 1.0 && id.MinQuantity > 0 && id.MaxQuantity >= id.MinQuantity
}

// HangupConfig 挂机配置值对象
type HangupConfig struct {
	MaxOfflineHours     int     `json:"max_offline_hours"`     // 最大离线小时数
	BaseExpPerHour      int64   `json:"base_exp_per_hour"`     // 基础每小时经验
	BaseGoldPerHour     int64   `json:"base_gold_per_hour"`    // 基础每小时金币
	VipBonusMultiplier  float64 `json:"vip_bonus_multiplier"`  // VIP加成倍率
	MaxDailyHangupHours int     `json:"max_daily_hangup_hours"` // 每日最大挂机小时数
	OfflineDecayRate    float64 `json:"offline_decay_rate"`    // 离线衰减率
	MinimumLevel        int     `json:"minimum_level"`         // 最低等级要求
}

// NewHangupConfig 创建挂机配置
func NewHangupConfig() *HangupConfig {
	return &HangupConfig{
		MaxOfflineHours:     24,
		BaseExpPerHour:      100,
		BaseGoldPerHour:     50,
		VipBonusMultiplier:  1.5,
		MaxDailyHangupHours: 12,
		OfflineDecayRate:    0.8, // 离线效率为在线的80%
		MinimumLevel:        1,
	}
}

// GetMaxOfflineHours 获取最大离线小时数
func (hc *HangupConfig) GetMaxOfflineHours() int {
	return hc.MaxOfflineHours
}

// GetBaseExpPerHour 获取基础每小时经验
func (hc *HangupConfig) GetBaseExpPerHour() int64 {
	return hc.BaseExpPerHour
}

// GetBaseGoldPerHour 获取基础每小时金币
func (hc *HangupConfig) GetBaseGoldPerHour() int64 {
	return hc.BaseGoldPerHour
}

// GetVipBonusMultiplier 获取VIP加成倍率
func (hc *HangupConfig) GetVipBonusMultiplier() float64 {
	return hc.VipBonusMultiplier
}

// GetMaxDailyHangupHours 获取每日最大挂机小时数
func (hc *HangupConfig) GetMaxDailyHangupHours() int {
	return hc.MaxDailyHangupHours
}

// GetOfflineDecayRate 获取离线衰减率
func (hc *HangupConfig) GetOfflineDecayRate() float64 {
	return hc.OfflineDecayRate
}

// GetMinimumLevel 获取最低等级要求
func (hc *HangupConfig) GetMinimumLevel() int {
	return hc.MinimumLevel
}

// IsValid 检查配置是否有效
func (hc *HangupConfig) IsValid() bool {
	return hc.MaxOfflineHours > 0 &&
		hc.BaseExpPerHour >= 0 &&
		hc.BaseGoldPerHour >= 0 &&
		hc.VipBonusMultiplier >= 1.0 &&
		hc.MaxDailyHangupHours > 0 &&
		hc.OfflineDecayRate > 0 && hc.OfflineDecayRate <= 1.0 &&
		hc.MinimumLevel >= 1
}

// HangupSession 挂机会话值对象
type HangupSession struct {
	SessionID   string        `json:"session_id"`
	PlayerID    string        `json:"player_id"`
	LocationID  string        `json:"location_id"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Duration    time.Duration `json:"duration"`
	IsOnline    bool          `json:"is_online"`
	Reward      *BaseReward   `json:"reward,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
}

// NewHangupSession 创建挂机会话
func NewHangupSession(sessionID, playerID, locationID string, isOnline bool) *HangupSession {
	now := time.Now()
	return &HangupSession{
		SessionID:  sessionID,
		PlayerID:   playerID,
		LocationID: locationID,
		StartTime:  now,
		IsOnline:   isOnline,
		CreatedAt:  now,
	}
}

// End 结束会话
func (hs *HangupSession) End(reward *BaseReward) {
	hs.EndTime = time.Now()
	hs.Duration = hs.EndTime.Sub(hs.StartTime)
	hs.Reward = reward
}

// GetSessionID 获取会话ID
func (hs *HangupSession) GetSessionID() string {
	return hs.SessionID
}

// GetPlayerID 获取玩家ID
func (hs *HangupSession) GetPlayerID() string {
	return hs.PlayerID
}

// GetLocationID 获取地点ID
func (hs *HangupSession) GetLocationID() string {
	return hs.LocationID
}

// GetStartTime 获取开始时间
func (hs *HangupSession) GetStartTime() time.Time {
	return hs.StartTime
}

// GetEndTime 获取结束时间
func (hs *HangupSession) GetEndTime() time.Time {
	return hs.EndTime
}

// GetDuration 获取持续时间
func (hs *HangupSession) GetDuration() time.Duration {
	return hs.Duration
}

// IsOnlineSession 是否在线会话
func (hs *HangupSession) IsOnlineSession() bool {
	return hs.IsOnline
}

// GetReward 获取奖励
func (hs *HangupSession) GetReward() *BaseReward {
	return hs.Reward
}

// IsActive 是否活跃会话
func (hs *HangupSession) IsActive() bool {
	return hs.EndTime.IsZero()
}

// GetCreatedAt 获取创建时间
func (hs *HangupSession) GetCreatedAt() time.Time {
	return hs.CreatedAt
}

// HangupRank 挂机排行值对象
type HangupRank struct {
	PlayerID        string        `json:"player_id"`
	PlayerName      string        `json:"player_name"`
	TotalHangupTime time.Duration `json:"total_hangup_time"`
	TotalExperience int64         `json:"total_experience"`
	TotalGold       int64         `json:"total_gold"`
	Rank            int           `json:"rank"`
	LastUpdated     time.Time     `json:"last_updated"`
}

// NewHangupRank 创建挂机排行
func NewHangupRank(playerID, playerName string, totalTime time.Duration, totalExp, totalGold int64, rank int) *HangupRank {
	return &HangupRank{
		PlayerID:        playerID,
		PlayerName:      playerName,
		TotalHangupTime: totalTime,
		TotalExperience: totalExp,
		TotalGold:       totalGold,
		Rank:            rank,
		LastUpdated:     time.Now(),
	}
}

// GetPlayerID 获取玩家ID
func (hr *HangupRank) GetPlayerID() string {
	return hr.PlayerID
}

// GetPlayerName 获取玩家名称
func (hr *HangupRank) GetPlayerName() string {
	return hr.PlayerName
}

// GetTotalHangupTime 获取总挂机时间
func (hr *HangupRank) GetTotalHangupTime() time.Duration {
	return hr.TotalHangupTime
}

// GetTotalExperience 获取总经验
func (hr *HangupRank) GetTotalExperience() int64 {
	return hr.TotalExperience
}

// GetTotalGold 获取总金币
func (hr *HangupRank) GetTotalGold() int64 {
	return hr.TotalGold
}

// GetRank 获取排名
func (hr *HangupRank) GetRank() int {
	return hr.Rank
}

// GetLastUpdated 获取最后更新时间
func (hr *HangupRank) GetLastUpdated() time.Time {
	return hr.LastUpdated
}

// UpdateRank 更新排名
func (hr *HangupRank) UpdateRank(newRank int) {
	hr.Rank = newRank
	hr.LastUpdated = time.Now()
}

// CalculateEfficiency 计算挂机效率
func (hr *HangupRank) CalculateEfficiency() float64 {
	if hr.TotalHangupTime == 0 {
		return 0
	}
	
	hours := hr.TotalHangupTime.Hours()
	expPerHour := float64(hr.TotalExperience) / hours
	goldPerHour := float64(hr.TotalGold) / hours
	
	// 综合效率计算（经验和金币的加权平均）
	return (expPerHour + goldPerHour*0.5) / 100 // 归一化到合理范围
}