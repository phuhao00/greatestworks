package battle

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"greatestworks/internal/domain/player"
)

// Service 战斗领域服务
type Service struct {
	battleRepository Repository
	skillRegistry    *SkillRegistry
}

// NewService 创建战斗领域服务
func NewService(battleRepository Repository) *Service {
	skillRegistry := NewSkillRegistry()
	skillRegistry.InitializeDefaultSkills()
	
	return &Service{
		battleRepository: battleRepository,
		skillRegistry:    skillRegistry,
	}
}

// CreateBattle 创建战斗
func (s *Service) CreateBattle(ctx context.Context, battleType BattleType) (*Battle, error) {
	battle := NewBattle(battleType)
	
	if err := s.battleRepository.Save(ctx, battle); err != nil {
		return nil, fmt.Errorf("save battle: %w", err)
	}
	
	return battle, nil
}

// JoinBattle 加入战斗
func (s *Service) JoinBattle(ctx context.Context, battleID BattleID, playerID player.PlayerID, team int, hp, mp int) error {
	battle, err := s.battleRepository.FindByID(ctx, battleID)
	if err != nil {
		return fmt.Errorf("find battle: %w", err)
	}
	
	if err := battle.AddParticipant(playerID, team, hp, mp); err != nil {
		return err
	}
	
	if err := s.battleRepository.Update(ctx, battle); err != nil {
		return fmt.Errorf("update battle: %w", err)
	}
	
	return nil
}

// StartBattle 开始战斗
func (s *Service) StartBattle(ctx context.Context, battleID BattleID) error {
	battle, err := s.battleRepository.FindByID(ctx, battleID)
	if err != nil {
		return fmt.Errorf("find battle: %w", err)
	}
	
	if err := battle.Start(); err != nil {
		return err
	}
	
	if err := s.battleRepository.Update(ctx, battle); err != nil {
		return fmt.Errorf("update battle: %w", err)
	}
	
	return nil
}

// ExecuteAttack 执行攻击
func (s *Service) ExecuteAttack(ctx context.Context, battleID BattleID, actorID player.PlayerID, targetID player.PlayerID) (*BattleAction, error) {
	battle, err := s.battleRepository.FindByID(ctx, battleID)
	if err != nil {
		return nil, fmt.Errorf("find battle: %w", err)
	}
	
	action, err := battle.ExecuteAction(actorID, &targetID, ActionTypeAttack, nil)
	if err != nil {
		return nil, err
	}
	
	if err := s.battleRepository.Update(ctx, battle); err != nil {
		return nil, fmt.Errorf("update battle: %w", err)
	}
	
	return action, nil
}

// ExecuteSkill 执行技能
func (s *Service) ExecuteSkill(ctx context.Context, battleID BattleID, actorID player.PlayerID, targetID *player.PlayerID, skillID string) (*BattleAction, error) {
	battle, err := s.battleRepository.FindByID(ctx, battleID)
	if err != nil {
		return nil, fmt.Errorf("find battle: %w", err)
	}
	
	// 验证技能是否存在
	skill, err := s.skillRegistry.GetSkill(skillID)
	if err != nil {
		return nil, fmt.Errorf("get skill: %w", err)
	}
	
	// 检查行动者是否有足够的魔法值
	actor := battle.findParticipant(actorID)
	if actor == nil {
		return nil, ErrPlayerNotInBattle
	}
	
	if actor.CurrentMP < skill.ManaCost() {
		return nil, ErrInsufficientMana
	}
	
	// 消耗魔法值
	actor.CurrentMP -= skill.ManaCost()
	
	// 执行技能
	action, err := s.executeSkillAction(battle, actorID, targetID, skill)
	if err != nil {
		return nil, err
	}
	
	if err := s.battleRepository.Update(ctx, battle); err != nil {
		return nil, fmt.Errorf("update battle: %w", err)
	}
	
	return action, nil
}

// executeSkillAction 执行技能行动
func (s *Service) executeSkillAction(battle *Battle, actorID player.PlayerID, targetID *player.PlayerID, skill *Skill) (*BattleAction, error) {
	action := &BattleAction{
		ActionID:   fmt.Sprintf("action_%d", time.Now().UnixNano()),
		ActorID:    actorID,
		TargetID:   targetID,
		ActionType: ActionTypeSkill,
		SkillID:    &skill.id.value,
		Timestamp:  time.Now(),
	}
	
	switch skill.GetSkillType() {
	case SkillTypeAttack:
		s.executeSkillAttack(battle, action, skill)
	case SkillTypeHeal:
		s.executeSkillHeal(battle, action, skill)
	case SkillTypeDefense:
		s.executeSkillDefense(battle, action, skill)
	case SkillTypeBuff:
		s.executeSkillBuff(battle, action, skill)
	case SkillTypeDebuff:
		s.executeSkillDebuff(battle, action, skill)
	}
	
	// 添加到当前回合
	battle.addActionToCurrentRound(action)
	
	// 检查战斗是否结束
	battle.checkBattleEnd()
	
	return action, nil
}

// executeSkillAttack 执行攻击技能
func (s *Service) executeSkillAttack(battle *Battle, action *BattleAction, skill *Skill) {
	if action.TargetID == nil {
		return
	}
	
	target := battle.findParticipant(*action.TargetID)
	if target == nil || !target.IsAlive {
		return
	}
	
	actor := battle.findParticipant(action.ActorID)
	if actor == nil {
		return
	}
	
	// 计算技能伤害
	damage := skill.Damage()
	
	// 暴击判断
	if s.rollCritical() {
		damage = int(float64(damage) * 1.5)
		action.Critical = true
	}
	
	action.Damage = damage
	target.CurrentHP -= damage
	target.DamageTaken += damage
	actor.DamageDealt += damage
	
	if target.CurrentHP <= 0 {
		target.CurrentHP = 0
		target.IsAlive = false
	}
	
	// 应用技能效果
	for _, effect := range skill.Effects() {
		s.applySkillEffect(target, effect)
	}
}

// executeSkillHeal 执行治疗技能
func (s *Service) executeSkillHeal(battle *Battle, action *BattleAction, skill *Skill) {
	var target *BattleParticipant
	
	if action.TargetID != nil {
		target = battle.findParticipant(*action.TargetID)
	} else {
		// 如果没有指定目标，治疗自己
		target = battle.findParticipant(action.ActorID)
	}
	
	if target == nil || !target.IsAlive {
		return
	}
	
	healAmount := skill.Healing()
	target.CurrentHP += healAmount
	action.Healing = healAmount
	
	// 不能超过最大生命值（这里简化处理，假设最大生命值为初始生命值）
	// 实际项目中应该从玩家数据中获取最大生命值
}

// executeSkillDefense 执行防御技能
func (s *Service) executeSkillDefense(battle *Battle, action *BattleAction, skill *Skill) {
	target := battle.findParticipant(action.ActorID)
	if target == nil {
		return
	}
	
	// 应用防御效果
	for _, effect := range skill.Effects() {
		s.applySkillEffect(target, effect)
	}
}

// executeSkillBuff 执行增益技能
func (s *Service) executeSkillBuff(battle *Battle, action *BattleAction, skill *Skill) {
	var target *BattleParticipant
	
	if action.TargetID != nil {
		target = battle.findParticipant(*action.TargetID)
	} else {
		target = battle.findParticipant(action.ActorID)
	}
	
	if target == nil {
		return
	}
	
	// 应用增益效果
	for _, effect := range skill.Effects() {
		s.applySkillEffect(target, effect)
	}
}

// executeSkillDebuff 执行减益技能
func (s *Service) executeSkillDebuff(battle *Battle, action *BattleAction, skill *Skill) {
	if action.TargetID == nil {
		return
	}
	
	target := battle.findParticipant(*action.TargetID)
	if target == nil || !target.IsAlive {
		return
	}
	
	// 应用减益效果
	for _, effect := range skill.Effects() {
		s.applySkillEffect(target, effect)
	}
}

// applySkillEffect 应用技能效果
func (s *Service) applySkillEffect(target *BattleParticipant, effect *SkillEffect) {
	// 这里简化处理，实际项目中应该有完整的效果系统
	switch effect.GetEffectType() {
	case EffectTypePoison:
		// 中毒效果，持续伤害
		target.CurrentHP -= effect.Value()
	case EffectTypeBurn:
		// 燃烧效果，持续伤害
		target.CurrentHP -= effect.Value()
	case EffectTypeAttackBoost:
		// 攻击力提升（这里简化处理）
	case EffectTypeDefenseBoost:
		// 防御力提升（这里简化处理）
	}
	
	if target.CurrentHP < 0 {
		target.CurrentHP = 0
		target.IsAlive = false
	}
}

// rollCritical 暴击判断
func (s *Service) rollCritical() bool {
	// 20% 暴击率
	return rand.Intn(100) < 20
}

// GetBattleStatus 获取战斗状态
func (s *Service) GetBattleStatus(ctx context.Context, battleID BattleID) (*Battle, error) {
	battle, err := s.battleRepository.FindByID(ctx, battleID)
	if err != nil {
		return nil, fmt.Errorf("find battle: %w", err)
	}
	
	return battle, nil
}

// GetPlayerBattles 获取玩家的战斗列表
func (s *Service) GetPlayerBattles(ctx context.Context, playerID player.PlayerID, limit int) ([]*Battle, error) {
	battles, err := s.battleRepository.FindByPlayerID(ctx, playerID, limit)
	if err != nil {
		return nil, fmt.Errorf("find player battles: %w", err)
	}
	
	return battles, nil
}

// EndBattle 结束战斗
func (s *Service) EndBattle(ctx context.Context, battleID BattleID) error {
	battle, err := s.battleRepository.FindByID(ctx, battleID)
	if err != nil {
		return fmt.Errorf("find battle: %w", err)
	}
	
	if battle.IsFinished() {
		return ErrBattleAlreadyFinished
	}
	
	// 强制结束战斗
	battle.status = BattleStatusCancelled
	now := time.Now()
	battle.endTime = &now
	battle.updatedAt = now
	battle.version++
	
	if err := s.battleRepository.Update(ctx, battle); err != nil {
		return fmt.Errorf("update battle: %w", err)
	}
	
	return nil
}

// CalculateBattleRewards 计算战斗奖励
func (s *Service) CalculateBattleRewards(ctx context.Context, battleID BattleID) (map[player.PlayerID]*BattleReward, error) {
	battle, err := s.battleRepository.FindByID(ctx, battleID)
	if err != nil {
		return nil, fmt.Errorf("find battle: %w", err)
	}
	
	if !battle.IsFinished() {
		return nil, ErrBattleNotFinished
	}
	
	rewards := make(map[player.PlayerID]*BattleReward)
	
	for _, participant := range battle.Participants() {
		reward := &BattleReward{
			PlayerID: participant.PlayerID,
			Exp:      s.calculateExpReward(participant, battle),
			Gold:     s.calculateGoldReward(participant, battle),
			Items:    s.calculateItemRewards(participant, battle),
		}
		
		// 获胜者额外奖励
		if battle.Winner() != nil && *battle.Winner() == participant.PlayerID {
			reward.Exp = int64(float64(reward.Exp) * 1.5)
			reward.Gold = int64(float64(reward.Gold) * 1.2)
		}
		
		rewards[participant.PlayerID] = reward
	}
	
	return rewards, nil
}

// calculateExpReward 计算经验奖励
func (s *Service) calculateExpReward(participant *BattleParticipant, battle *Battle) int64 {
	baseExp := int64(100)
	
	// 根据造成的伤害调整
	damageBonus := int64(participant.DamageDealt / 10)
	
	// 根据战斗时长调整
	duration := battle.endTime.Sub(battle.startTime)
	timeBonus := int64(duration.Minutes() * 5)
	
	return baseExp + damageBonus + timeBonus
}

// calculateGoldReward 计算金币奖励
func (s *Service) calculateGoldReward(participant *BattleParticipant, battle *Battle) int64 {
	baseGold := int64(50)
	
	// 根据造成的伤害调整
	damageBonus := int64(participant.DamageDealt / 20)
	
	return baseGold + damageBonus
}

// calculateItemRewards 计算物品奖励
func (s *Service) calculateItemRewards(participant *BattleParticipant, battle *Battle) []string {
	items := make([]string, 0)
	
	// 简单的随机掉落系统
	if rand.Intn(100) < 30 { // 30% 概率掉落物品
		items = append(items, "health_potion")
	}
	
	if rand.Intn(100) < 20 { // 20% 概率掉落装备
		items = append(items, "basic_sword")
	}
	
	return items
}

// BattleReward 战斗奖励
type BattleReward struct {
	PlayerID player.PlayerID `json:"player_id"`
	Exp      int64           `json:"exp"`
	Gold     int64           `json:"gold"`
	Items    []string        `json:"items"`
}