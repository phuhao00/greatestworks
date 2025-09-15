package hangup

import (
	"time"
)

// HangupLocation 挂机地点实体
type HangupLocation struct {
	id              string
	name            string
	description     string
	locationType    LocationType
	requiredLevel   int
	requiredQuests  []string
	baseExpRate     float64 // 基础经验倍率
	baseGoldRate    float64 // 基础金币倍率
	specialItems    []ItemDrop
	maxOfflineHours int // 最大离线小时数
	isUnlocked      bool
	isActive        bool
	createdAt       time.Time
	updatedAt       time.Time
}

// NewHangupLocation 创建挂机地点
func NewHangupLocation(id, name, description string, locationType LocationType) *HangupLocation {
	now := time.Now()
	return &HangupLocation{
		id:              id,
		name:            name,
		description:     description,
		locationType:    locationType,
		requiredLevel:   1,
		requiredQuests:  make([]string, 0),
		baseExpRate:     1.0,
		baseGoldRate:    1.0,
		specialItems:    make([]ItemDrop, 0),
		maxOfflineHours: 24,
		isUnlocked:      false,
		isActive:        true,
		createdAt:       now,
		updatedAt:       now,
	}
}

// GetID 获取地点ID
func (hl *HangupLocation) GetID() string {
	return hl.id
}

// GetName 获取地点名称
func (hl *HangupLocation) GetName() string {
	return hl.name
}

// GetDescription 获取地点描述
func (hl *HangupLocation) GetDescription() string {
	return hl.description
}

// GetLocationType 获取地点类型
func (hl *HangupLocation) GetLocationType() LocationType {
	return hl.locationType
}

// GetRequiredLevel 获取所需等级
func (hl *HangupLocation) GetRequiredLevel() int {
	return hl.requiredLevel
}

// SetRequiredLevel 设置所需等级
func (hl *HangupLocation) SetRequiredLevel(level int) {
	hl.requiredLevel = level
	hl.updatedAt = time.Now()
}

// GetRequiredQuests 获取所需任务
func (hl *HangupLocation) GetRequiredQuests() []string {
	return hl.requiredQuests
}

// AddRequiredQuest 添加所需任务
func (hl *HangupLocation) AddRequiredQuest(questID string) {
	hl.requiredQuests = append(hl.requiredQuests, questID)
	hl.updatedAt = time.Now()
}

// GetBaseExpRate 获取基础经验倍率
func (hl *HangupLocation) GetBaseExpRate() float64 {
	return hl.baseExpRate
}

// SetBaseExpRate 设置基础经验倍率
func (hl *HangupLocation) SetBaseExpRate(rate float64) {
	hl.baseExpRate = rate
	hl.updatedAt = time.Now()
}

// GetBaseGoldRate 获取基础金币倍率
func (hl *HangupLocation) GetBaseGoldRate() float64 {
	return hl.baseGoldRate
}

// SetBaseGoldRate 设置基础金币倍率
func (hl *HangupLocation) SetBaseGoldRate(rate float64) {
	hl.baseGoldRate = rate
	hl.updatedAt = time.Now()
}

// GetSpecialItems 获取特殊物品掉落
func (hl *HangupLocation) GetSpecialItems() []ItemDrop {
	return hl.specialItems
}

// AddSpecialItem 添加特殊物品掉落
func (hl *HangupLocation) AddSpecialItem(item ItemDrop) {
	hl.specialItems = append(hl.specialItems, item)
	hl.updatedAt = time.Now()
}

// GetMaxOfflineHours 获取最大离线小时数
func (hl *HangupLocation) GetMaxOfflineHours() int {
	return hl.maxOfflineHours
}

// SetMaxOfflineHours 设置最大离线小时数
func (hl *HangupLocation) SetMaxOfflineHours(hours int) {
	hl.maxOfflineHours = hours
	hl.updatedAt = time.Now()
}

// IsUnlocked 是否已解锁
func (hl *HangupLocation) IsUnlocked() bool {
	return hl.isUnlocked
}

// Unlock 解锁地点
func (hl *HangupLocation) Unlock() {
	hl.isUnlocked = true
	hl.updatedAt = time.Now()
}

// Lock 锁定地点
func (hl *HangupLocation) Lock() {
	hl.isUnlocked = false
	hl.updatedAt = time.Now()
}

// IsActive 是否激活
func (hl *HangupLocation) IsActive() bool {
	return hl.isActive
}

// Activate 激活地点
func (hl *HangupLocation) Activate() {
	hl.isActive = true
	hl.updatedAt = time.Now()
}

// Deactivate 停用地点
func (hl *HangupLocation) Deactivate() {
	hl.isActive = false
	hl.updatedAt = time.Now()
}

// CalculateBaseReward 计算基础奖励
func (hl *HangupLocation) CalculateBaseReward(duration time.Duration) *BaseReward {
	hours := duration.Hours()
	
	// 限制最大离线时间
	if hours > float64(hl.maxOfflineHours) {
		hours = float64(hl.maxOfflineHours)
	}
	
	// 计算基础奖励（这里使用简单的线性计算）
	baseExp := int64(hours * 100 * hl.baseExpRate)   // 每小时100经验
	baseGold := int64(hours * 50 * hl.baseGoldRate)  // 每小时50金币
	
	// 计算物品掉落
	items := make([]RewardItem, 0)
	for _, itemDrop := range hl.specialItems {
		if itemDrop.ShouldDrop(hours) {
			items = append(items, RewardItem{
				ItemID:   itemDrop.ItemID,
				Quantity: itemDrop.CalculateQuantity(hours),
			})
		}
	}
	
	return &BaseReward{
		Experience: baseExp,
		Gold:       baseGold,
		Items:      items,
	}
}

// GetCreatedAt 获取创建时间
func (hl *HangupLocation) GetCreatedAt() time.Time {
	return hl.createdAt
}

// GetUpdatedAt 获取更新时间
func (hl *HangupLocation) GetUpdatedAt() time.Time {
	return hl.updatedAt
}

// OfflineReward 离线奖励实体
type OfflineReward struct {
	Experience      int64         `json:"experience"`
	Gold            int64         `json:"gold"`
	Items           []RewardItem  `json:"items"`
	OfflineDuration time.Duration `json:"offline_duration"`
	LocationID      string        `json:"location_id"`
	CalculatedAt    time.Time     `json:"calculated_at"`
	IsClaimed       bool          `json:"is_claimed"`
	ClaimedAt       time.Time     `json:"claimed_at,omitempty"`
}

// NewOfflineReward 创建离线奖励
func NewOfflineReward() *OfflineReward {
	return &OfflineReward{
		Experience:   0,
		Gold:         0,
		Items:        make([]RewardItem, 0),
		CalculatedAt: time.Now(),
		IsClaimed:    false,
	}
}

// IsEmpty 是否为空奖励
func (or *OfflineReward) IsEmpty() bool {
	return or.Experience == 0 && or.Gold == 0 && len(or.Items) == 0
}

// GetTotalValue 获取总价值（用于显示）
func (or *OfflineReward) GetTotalValue() int64 {
	// 简单的价值计算：经验 + 金币 + 物品价值
	totalValue := or.Experience + or.Gold
	
	for _, item := range or.Items {
		// 假设每个物品价值10金币
		totalValue += int64(item.Quantity) * 10
	}
	
	return totalValue
}

// EfficiencyBonus 效率加成实体
type EfficiencyBonus struct {
	vipBonus       float64            `json:"vip_bonus"`        // VIP加成
	equipmentBonus float64            `json:"equipment_bonus"`  // 装备加成
	skillBonus     float64            `json:"skill_bonus"`      // 技能加成
	guildBonus     float64            `json:"guild_bonus"`      // 公会加成
	eventBonus     float64            `json:"event_bonus"`      // 活动加成
	specialBonus   map[string]float64 `json:"special_bonus"`    // 特殊加成
	updatedAt      time.Time          `json:"updated_at"`
}

// NewEfficiencyBonus 创建效率加成
func NewEfficiencyBonus() *EfficiencyBonus {
	return &EfficiencyBonus{
		vipBonus:       0.0,
		equipmentBonus: 0.0,
		skillBonus:     0.0,
		guildBonus:     0.0,
		eventBonus:     0.0,
		specialBonus:   make(map[string]float64),
		updatedAt:      time.Now(),
	}
}

// GetVipBonus 获取VIP加成
func (eb *EfficiencyBonus) GetVipBonus() float64 {
	return eb.vipBonus
}

// SetVipBonus 设置VIP加成
func (eb *EfficiencyBonus) SetVipBonus(bonus float64) {
	eb.vipBonus = bonus
	eb.updatedAt = time.Now()
}

// GetEquipmentBonus 获取装备加成
func (eb *EfficiencyBonus) GetEquipmentBonus() float64 {
	return eb.equipmentBonus
}

// SetEquipmentBonus 设置装备加成
func (eb *EfficiencyBonus) SetEquipmentBonus(bonus float64) {
	eb.equipmentBonus = bonus
	eb.updatedAt = time.Now()
}

// GetSkillBonus 获取技能加成
func (eb *EfficiencyBonus) GetSkillBonus() float64 {
	return eb.skillBonus
}

// SetSkillBonus 设置技能加成
func (eb *EfficiencyBonus) SetSkillBonus(bonus float64) {
	eb.skillBonus = bonus
	eb.updatedAt = time.Now()
}

// GetGuildBonus 获取公会加成
func (eb *EfficiencyBonus) GetGuildBonus() float64 {
	return eb.guildBonus
}

// SetGuildBonus 设置公会加成
func (eb *EfficiencyBonus) SetGuildBonus(bonus float64) {
	eb.guildBonus = bonus
	eb.updatedAt = time.Now()
}

// GetEventBonus 获取活动加成
func (eb *EfficiencyBonus) GetEventBonus() float64 {
	return eb.eventBonus
}

// SetEventBonus 设置活动加成
func (eb *EfficiencyBonus) SetEventBonus(bonus float64) {
	eb.eventBonus = bonus
	eb.updatedAt = time.Now()
}

// GetSpecialBonus 获取特殊加成
func (eb *EfficiencyBonus) GetSpecialBonus(key string) float64 {
	return eb.specialBonus[key]
}

// SetSpecialBonus 设置特殊加成
func (eb *EfficiencyBonus) SetSpecialBonus(key string, bonus float64) {
	eb.specialBonus[key] = bonus
	eb.updatedAt = time.Now()
}

// GetTotalBonus 获取总加成
func (eb *EfficiencyBonus) GetTotalBonus() float64 {
	total := 1.0 + eb.vipBonus + eb.equipmentBonus + eb.skillBonus + eb.guildBonus + eb.eventBonus
	
	for _, bonus := range eb.specialBonus {
		total += bonus
	}
	
	return total
}

// ApplyBonus 应用加成到基础奖励
func (eb *EfficiencyBonus) ApplyBonus(baseReward *BaseReward) *BaseReward {
	totalBonus := eb.GetTotalBonus()
	
	return &BaseReward{
		Experience: int64(float64(baseReward.Experience) * totalBonus),
		Gold:       int64(float64(baseReward.Gold) * totalBonus),
		Items:      baseReward.Items, // 物品不受加成影响
	}
}

// GetUpdatedAt 获取更新时间
func (eb *EfficiencyBonus) GetUpdatedAt() time.Time {
	return eb.updatedAt
}

// HangupStatistics 挂机统计实体
type HangupStatistics struct {
	playerID           string        `json:"player_id"`
	totalHangupTime    time.Duration `json:"total_hangup_time"`
	totalExperience    int64         `json:"total_experience"`
	totalGold          int64         `json:"total_gold"`
	totalItemsObtained int           `json:"total_items_obtained"`
	favoriteLocation   string        `json:"favorite_location"`
	longestSession     time.Duration `json:"longest_session"`
	lastHangupDate     time.Time     `json:"last_hangup_date"`
	updatedAt          time.Time     `json:"updated_at"`
}

// NewHangupStatistics 创建挂机统计
func NewHangupStatistics(playerID string) *HangupStatistics {
	return &HangupStatistics{
		playerID:           playerID,
		totalHangupTime:    0,
		totalExperience:    0,
		totalGold:          0,
		totalItemsObtained: 0,
		favoriteLocation:   "",
		longestSession:     0,
		lastHangupDate:     time.Time{},
		updatedAt:          time.Now(),
	}
}

// UpdateStatistics 更新统计数据
func (hs *HangupStatistics) UpdateStatistics(sessionDuration time.Duration, reward *OfflineReward, locationID string) {
	hs.totalHangupTime += sessionDuration
	hs.totalExperience += reward.Experience
	hs.totalGold += reward.Gold
	hs.totalItemsObtained += len(reward.Items)
	
	if sessionDuration > hs.longestSession {
		hs.longestSession = sessionDuration
	}
	
	hs.favoriteLocation = locationID // 简化实现，实际应该统计最常用的地点
	hs.lastHangupDate = time.Now()
	hs.updatedAt = time.Now()
}

// GetPlayerID 获取玩家ID
func (hs *HangupStatistics) GetPlayerID() string {
	return hs.playerID
}

// GetTotalHangupTime 获取总挂机时间
func (hs *HangupStatistics) GetTotalHangupTime() time.Duration {
	return hs.totalHangupTime
}

// GetTotalExperience 获取总经验
func (hs *HangupStatistics) GetTotalExperience() int64 {
	return hs.totalExperience
}

// GetTotalGold 获取总金币
func (hs *HangupStatistics) GetTotalGold() int64 {
	return hs.totalGold
}

// GetTotalItemsObtained 获取总物品数量
func (hs *HangupStatistics) GetTotalItemsObtained() int {
	return hs.totalItemsObtained
}

// GetFavoriteLocation 获取最喜欢的地点
func (hs *HangupStatistics) GetFavoriteLocation() string {
	return hs.favoriteLocation
}

// GetLongestSession 获取最长会话时间
func (hs *HangupStatistics) GetLongestSession() time.Duration {
	return hs.longestSession
}

// GetLastHangupDate 获取最后挂机日期
func (hs *HangupStatistics) GetLastHangupDate() time.Time {
	return hs.lastHangupDate
}

// GetUpdatedAt 获取更新时间
func (hs *HangupStatistics) GetUpdatedAt() time.Time {
	return hs.updatedAt
}