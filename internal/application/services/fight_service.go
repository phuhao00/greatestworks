package services

import (
	"context"
	"errors"
	"math/rand/v2"

	"greatestworks/internal/domain/character"
	"greatestworks/internal/infrastructure/datamanager"
)

// SkillCastResult 技能释放结果
type SkillCastResult struct {
	CasterID   int32
	TargetID   int32
	SkillID    int32
	Damage     int32
	IsCritical bool
	Success    bool
	Message    string
}

// FightService 战斗服务
type FightService struct {
	characterService *CharacterService
	mapService       *MapService
}

// NewFightService 创建战斗服务
func NewFightService(characterService *CharacterService) *FightService {
	return &FightService{
		characterService: characterService,
	}
}

// SetMapService 注入地图服务用于查找实体
func (s *FightService) SetMapService(ms *MapService) {
	s.mapService = ms
}

// CastSkill 释放技能
func (s *FightService) CastSkill(ctx context.Context, caster *character.Actor, targetID int64, skillID int32) error {
	// 获取技能定义
	skillDefine := datamanager.GetInstance().GetSkill(skillID)
	if skillDefine == nil {
		return errors.New("skill not found")
	}

	// 检查技能是否学会
	skill := caster.GetSkillManager().GetSkill(skillID)
	if skill == nil {
		return errors.New("skill not learned")
	}

	// 检查技能状态
	if skill.State() != character.SkillStateReady {
		return errors.New("skill is not ready")
	}

	// 使用施法器释放技能（此处未解析targetID，默认无目标）
	if ok := caster.GetSpell().Cast(skillID, nil); !ok {
		return errors.New("cast failed")
	}

	return nil
}

// CastSkillByID 基于ID的施法接口：计算伤害并返回结果（不直接应用）
// handler 层负责广播结果；实际伤害应用可选地通过后续流程完成
func (s *FightService) CastSkillByID(ctx context.Context, casterEntityID int32, targetEntityID int32, skillID int32) (*SkillCastResult, error) {
	result := &SkillCastResult{
		CasterID: casterEntityID,
		TargetID: targetEntityID,
		SkillID:  skillID,
		Success:  false,
	}

	// 获取技能定义
	skillDefine := datamanager.GetInstance().GetSkill(skillID)
	if skillDefine == nil {
		result.Message = "skill not found"
		return result, errors.New("skill not found")
	}

	// 简化计算：基于技能基础伤害与固定暴击率
	baseDamage := skillDefine.BaseDamage
	if baseDamage == 0 {
		baseDamage = 100 // 默认技能伤害
	}

	// 暴击判定（固定10%暴击率，伤害1.5倍）
	isCrit := (rand.Int32N(100) < 10)
	totalDamage := baseDamage
	if isCrit {
		totalDamage = int32(float32(totalDamage) * 1.5)
	}

	result.Damage = totalDamage
	result.IsCritical = isCrit
	result.Success = true
	result.Message = "skill cast calculated"
	return result, nil
}

// ApplyDamage 应用伤害
func (s *FightService) ApplyDamage(ctx context.Context, attacker, target *character.Actor, damage int32, dmgType int32) error {
	if target == nil {
		return errors.New("target is nil")
	}

	// 应用伤害
	info := &character.DamageInfo{
		TargetID: target.ID(),
		AttackerInfo: character.AttackerInfo{
			AttackerID:   attacker.ID(),
			AttackerType: character.AttackerTypeNormal,
		},
		Amount:     damage,
		DamageType: character.DamageType(dmgType),
	}
	_ = target.OnHurt(ctx, info)

	// 检查死亡
	if target.IsDeath() {
		s.onActorDeath(ctx, target, attacker)
	}

	return nil
}

// ApplyHeal 应用治疗
func (s *FightService) ApplyHeal(ctx context.Context, caster, target *character.Actor, heal int32) error {
	if target == nil {
		return errors.New("target is nil")
	}

	currentHP := target.HP()
	maxHP := target.GetAttributeManager().Final().MaxHP
	newHP := currentHP + float32(heal)
	if newHP > maxHP {
		newHP = maxHP
	}
	target.ChangeHP(newHP - currentHP)
	return nil
}

// ApplyBuff 应用Buff
func (s *FightService) ApplyBuff(ctx context.Context, caster, target *character.Actor, buffID int32, duration float32) error {
	if target == nil {
		return errors.New("target is nil")
	}

	// 创建并添加Buff
	buff := character.NewBuff(buffID, target, caster, duration)
	target.GetBuffManager().AddBuff(buff)

	return nil
}

// RemoveBuff 移除Buff
func (s *FightService) RemoveBuff(ctx context.Context, target *character.Actor, buffID int32) error {
	if target == nil {
		return errors.New("target is nil")
	}

	target.GetBuffManager().RemoveBuffByID(buffID)
	return nil
}

// UpdateFight 更新战斗状态
func (s *FightService) UpdateFight(ctx context.Context, actor *character.Actor, deltaTime float32) {
	// 更新技能与Buff
	_ = actor.GetSkillManager().Update(ctx, deltaTime)
	_ = actor.GetBuffManager().Update(ctx, deltaTime)
}

// onActorDeath 角色死亡处理
func (s *FightService) onActorDeath(ctx context.Context, deadActor, killer *character.Actor) {
	// TODO: 掉落物品、经验奖励、复活处理等
}

// Resurrect 复活
func (s *FightService) Resurrect(ctx context.Context, actor *character.Actor) error {
	if !actor.IsDeath() {
		return errors.New("actor is not dead")
	}
	return actor.Revive(ctx)
}
