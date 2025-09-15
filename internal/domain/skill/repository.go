package skill

import (
	"context"
	"time"
)

// Repository 技能仓储接口
type Repository interface {
	// 基础CRUD操作
	Save(ctx context.Context, skillTree *SkillTree) error
	FindByPlayerID(ctx context.Context, playerID string) (*SkillTree, error)
	Delete(ctx context.Context, playerID string) error
	Exists(ctx context.Context, playerID string) (bool, error)
	
	// 批量操作
	SaveBatch(ctx context.Context, skillTrees []*SkillTree) error
	FindByPlayerIDs(ctx context.Context, playerIDs []string) ([]*SkillTree, error)
	
	// 技能查询
	FindSkillsByType(ctx context.Context, playerID string, skillType SkillType) ([]*Skill, error)
	FindLearnedSkills(ctx context.Context, playerID string) ([]*Skill, error)
	FindAvailableSkills(ctx context.Context, playerID string) ([]*Skill, error)
	GetSkillLevel(ctx context.Context, playerID string, skillID string) (int, error)
	
	// 技能统计
	GetSkillStats(ctx context.Context, playerID string) (*SkillStats, error)
	GetSkillUsageHistory(ctx context.Context, playerID string, limit int) ([]*SkillUsageRecord, error)
	GetTopSkills(ctx context.Context, playerID string, limit int) ([]*SkillRanking, error)
	
	// 技能配置
	GetSkillConfig(ctx context.Context, skillID string) (*SkillConfig, error)
	GetAllSkillConfigs(ctx context.Context) ([]*SkillConfig, error)
	SaveSkillConfig(ctx context.Context, config *SkillConfig) error
}

// SkillStats 技能统计信息
type SkillStats struct {
	PlayerID        string            `json:"player_id"`
	TotalSkills     int               `json:"total_skills"`
	LearnedSkills   int               `json:"learned_skills"`
	SkillPoints     int64             `json:"skill_points"`
	TotalPoints     int64             `json:"total_points"`
	SkillsByType    map[SkillType]int `json:"skills_by_type"`
	HighestLevel    int               `json:"highest_level"`
	AverageLevel    float64           `json:"average_level"`
	LastUpdate      time.Time         `json:"last_update"`
}

// SkillUsageRecord 技能使用记录
type SkillUsageRecord struct {
	ID         string    `json:"id"`
	PlayerID   string    `json:"player_id"`
	SkillID    string    `json:"skill_id"`
	SkillName  string    `json:"skill_name"`
	TargetID   string    `json:"target_id"`
	Damage     int64     `json:"damage"`
	Success    bool      `json:"success"`
	UsedAt     time.Time `json:"used_at"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// SkillRanking 技能排行
type SkillRanking struct {
	SkillID    string `json:"skill_id"`
	SkillName  string `json:"skill_name"`
	Level      int    `json:"level"`
	UsageCount int64  `json:"usage_count"`
	TotalDamage int64 `json:"total_damage"`
	SuccessRate float64 `json:"success_rate"`
}

// SkillConfig 技能配置
type SkillConfig struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	SkillType     SkillType              `json:"skill_type"`
	MaxLevel      int                    `json:"max_level"`
	Prerequisites []string               `json:"prerequisites"`
	BaseDamage    int64                  `json:"base_damage"`
	DamageType    DamageType             `json:"damage_type"`
	ManaCost      int64                  `json:"mana_cost"`
	Cooldown      time.Duration          `json:"cooldown"`
	CastTime      time.Duration          `json:"cast_time"`
	Range         float64                `json:"range"`
	Effects       []*SkillEffectConfig   `json:"effects"`
	Scaling       map[AttributeType]float64 `json:"scaling"`
	LevelRequirement int                 `json:"level_requirement"`
	ClassRestrictions []string           `json:"class_restrictions"`
	RaceRestrictions  []string           `json:"race_restrictions"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// SkillEffectConfig 技能效果配置
type SkillEffectConfig struct {
	EffectType EffectType    `json:"effect_type"`
	Value      float64       `json:"value"`
	Duration   time.Duration `json:"duration"`
	Target     TargetType    `json:"target"`
	Condition  *EffectConditionConfig `json:"condition,omitempty"`
}

// EffectConditionConfig 效果条件配置
type EffectConditionConfig struct {
	ConditionType ConditionType `json:"condition_type"`
	Value         interface{}   `json:"value"`
}

// SkillQueryFilter 技能查询过滤器
type SkillQueryFilter struct {
	PlayerID          string      `json:"player_id"`
	SkillTypes        []SkillType `json:"skill_types,omitempty"`
	MinLevel          *int        `json:"min_level,omitempty"`
	MaxLevel          *int        `json:"max_level,omitempty"`
	LearnedOnly       bool        `json:"learned_only"`
	AvailableOnly     bool        `json:"available_only"`
	UsableOnly        bool        `json:"usable_only"`
	IncludePassive    bool        `json:"include_passive"`
	SortBy            string      `json:"sort_by"` // level, usage_count, damage
	SortOrder         string      `json:"sort_order"` // asc, desc
	Limit             int         `json:"limit"`
	Offset            int         `json:"offset"`
}