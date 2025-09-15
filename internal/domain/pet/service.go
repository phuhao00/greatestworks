package pet

import (
	"fmt"
	"math/rand"
	"time"
)

// PetService 宠物领域服务
type PetService struct {
	petTemplates      map[uint32]*PetTemplate
	skillTemplates    map[string]*SkillTemplate
	skinTemplates     map[string]*SkinTemplate
	bondTemplates     map[string]*BondTemplate
	fragmentTemplates map[uint32]*FragmentTemplate
}

// NewPetService 创建宠物领域服务
func NewPetService() *PetService {
	return &PetService{
		petTemplates:      make(map[uint32]*PetTemplate),
		skillTemplates:    make(map[string]*SkillTemplate),
		skinTemplates:     make(map[string]*SkinTemplate),
		bondTemplates:     make(map[string]*BondTemplate),
		fragmentTemplates: make(map[uint32]*FragmentTemplate),
	}
}

// CreatePet 创建宠物
func (ps *PetService) CreatePet(playerID string, configID uint32, name string) (*PetAggregate, error) {
	template, exists := ps.petTemplates[configID]
	if !exists {
		return nil, ErrPetTemplateNotFound
	}
	
	if name == "" {
		name = template.DefaultName
	}
	
	pet := NewPetAggregate(playerID, configID, name, template.Category)
	
	// 应用模板属性
	ps.applyPetTemplate(pet, template)
	
	// 添加初始技能
	for _, skillID := range template.InitialSkills {
		if skillTemplate, exists := ps.skillTemplates[skillID]; exists {
			skill := ps.createSkillFromTemplate(skillTemplate)
			pet.AddSkill(skill)
		}
	}
	
	return pet, nil
}

// SummonPetFromFragments 通过碎片召唤宠物
func (ps *PetService) SummonPetFromFragments(playerID string, fragments []*PetFragment) (*PetAggregate, error) {
	if len(fragments) == 0 {
		return nil, ErrNoFragmentsProvided
	}
	
	// 检查所有碎片是否属于同一宠物
	relatedPetID := fragments[0].GetRelatedPetID()
	var totalQuantity uint64
	
	for _, fragment := range fragments {
		if fragment.GetRelatedPetID() != relatedPetID {
			return nil, ErrFragmentMismatch
		}
		totalQuantity += fragment.GetQuantity()
	}
	
	// 检查碎片数量是否足够
	fragmentTemplate, exists := ps.fragmentTemplates[fragments[0].GetFragmentID()]
	if !exists {
		return nil, ErrFragmentTemplateNotFound
	}
	
	if totalQuantity < fragmentTemplate.RequiredQuantity {
		return nil, ErrInsufficientFragments
	}
	
	// 消耗碎片
	remainingQuantity := fragmentTemplate.RequiredQuantity
	for _, fragment := range fragments {
		if remainingQuantity <= 0 {
			break
		}
		
		consumeAmount := fragment.GetQuantity()
		if consumeAmount > remainingQuantity {
			consumeAmount = remainingQuantity
		}
		
		if err := fragment.ConsumeQuantity(consumeAmount); err != nil {
			return nil, err
		}
		
		remainingQuantity -= consumeAmount
	}
	
	// 创建宠物
	return ps.CreatePet(playerID, relatedPetID, "")
}

// EvolvePet 宠物进化
func (ps *PetService) EvolvePet(pet *PetAggregate, materials []string) error {
	if pet.GetStar() >= MaxPetStar {
		return ErrMaxStarReached
	}
	
	// 检查进化材料
	if !ps.checkEvolutionMaterials(pet.GetConfigID(), pet.GetStar(), materials) {
		return ErrInsufficientEvolutionMaterials
	}
	
	// 执行进化
	if err := pet.UpgradeStar(); err != nil {
		return err
	}
	
	// 进化后可能解锁新技能
	ps.unlockEvolutionSkills(pet)
	
	return nil
}

// TrainPet 训练宠物
func (ps *PetService) TrainPet(pet *PetAggregate, trainingType TrainingType, duration time.Duration) error {
	if !trainingType.IsValid() {
		return ErrInvalidTrainingType
	}
	
	// 检查宠物状态
	if err := pet.Train(trainingType, duration); err != nil {
		return err
	}
	
	return nil
}

// FinishPetTraining 完成宠物训练
func (ps *PetService) FinishPetTraining(pet *PetAggregate, trainingType TrainingType) (*TrainingResult, error) {
	if err := pet.FinishTraining(trainingType); err != nil {
		return nil, err
	}
	
	// 计算训练结果
	result := ps.calculateTrainingResult(pet, trainingType)
	
	return result, nil
}

// FeedPet 喂食宠物
func (ps *PetService) FeedPet(pet *PetAggregate, foodType FoodType, amount int) (*FeedingResult, error) {
	if !foodType.IsValid() {
		return nil, ErrInvalidFoodType
	}
	
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	
	// 记录喂食前的状态
	beforeExp := pet.GetExperience()
	beforeLevel := pet.GetLevel()
	beforeAttributes := pet.GetAttributes().Clone()
	
	// 执行喂食
	if err := pet.Feed(foodType, amount); err != nil {
		return nil, err
	}
	
	// 计算喂食结果
	result := &FeedingResult{
		FoodType:        foodType,
		Amount:          amount,
		ExperienceGain:  pet.GetExperience() - beforeExp,
		LevelUp:         pet.GetLevel() > beforeLevel,
		AttributeChanges: ps.calculateAttributeChanges(beforeAttributes, pet.GetAttributes()),
		Timestamp:       time.Now(),
	}
	
	return result, nil
}

// BreedPets 宠物繁殖
func (ps *PetService) BreedPets(parent1, parent2 *PetAggregate) (*PetAggregate, error) {
	// 检查繁殖条件
	if err := ps.validateBreedingConditions(parent1, parent2); err != nil {
		return nil, err
	}
	
	// 生成后代
	offspring := ps.generateOffspring(parent1, parent2)
	
	return offspring, nil
}

// CalculateBattlePower 计算战斗力
func (ps *PetService) CalculateBattlePower(pet *PetAggregate) int64 {
	basePower := pet.GetTotalPower()
	
	// 技能加成
	skillBonus := ps.calculateSkillPowerBonus(pet.GetSkills())
	
	// 等级加成
	levelBonus := int64(pet.GetLevel()) * 10
	
	// 星级加成
	starBonus := int64(pet.GetStar()) * 100
	
	return basePower + skillBonus + levelBonus + starBonus
}

// GenerateRandomPet 生成随机宠物
func (ps *PetService) GenerateRandomPet(playerID string, rarity PetRarity) (*PetAggregate, error) {
	// 根据稀有度筛选可用模板
	availableTemplates := ps.getTemplatesByRarity(rarity)
	if len(availableTemplates) == 0 {
		return nil, ErrNoAvailableTemplates
	}
	
	// 随机选择模板
	template := availableTemplates[rand.Intn(len(availableTemplates))]
	
	// 创建宠物
	pet, err := ps.CreatePet(playerID, template.ConfigID, "")
	if err != nil {
		return nil, err
	}
	
	// 应用随机属性加成
	ps.applyRandomBonus(pet, rarity)
	
	return pet, nil
}

// 私有方法

// applyPetTemplate 应用宠物模板
func (ps *PetService) applyPetTemplate(pet *PetAggregate, template *PetTemplate) {
	// 应用基础属性
	attributes := pet.GetAttributes()
	attributes.AddHealth(template.BaseHealth)
	attributes.AddAttack(template.BaseAttack)
	attributes.AddDefense(template.BaseDefense)
	attributes.AddSpeed(template.BaseSpeed)
	attributes.AddCritical(template.BaseCritical)
	attributes.AddHit(template.BaseHit)
	attributes.AddDodge(template.BaseDodge)
}

// createSkillFromTemplate 从模板创建技能
func (ps *PetService) createSkillFromTemplate(template *SkillTemplate) *PetSkill {
	skill := NewPetSkill(template.SkillID, template.Name, template.Type, template.BaseDamage, template.Cooldown)
	
	// 添加技能效果
	for _, effectTemplate := range template.Effects {
		effect := SkillEffect{
			EffectType: effectTemplate.Type,
			Value:      effectTemplate.Value,
			Duration:   effectTemplate.Duration,
		}
		skill.AddEffect(effect)
	}
	
	return skill
}

// checkEvolutionMaterials 检查进化材料
func (ps *PetService) checkEvolutionMaterials(configID uint32, currentStar uint32, materials []string) bool {
	// 简化实现，实际应该根据配置检查
	requiredMaterials := ps.getRequiredEvolutionMaterials(configID, currentStar)
	
	// 检查材料是否足够
	materialCount := make(map[string]int)
	for _, material := range materials {
		materialCount[material]++
	}
	
	for material, required := range requiredMaterials {
		if materialCount[material] < required {
			return false
		}
	}
	
	return true
}

// getRequiredEvolutionMaterials 获取进化所需材料
func (ps *PetService) getRequiredEvolutionMaterials(configID uint32, currentStar uint32) map[string]int {
	// 简化实现
	return map[string]int{
		"evolution_stone": int(currentStar),
		"gold":            int(currentStar) * 1000,
	}
}

// unlockEvolutionSkills 解锁进化技能
func (ps *PetService) unlockEvolutionSkills(pet *PetAggregate) {
	// 根据星级解锁新技能
	star := pet.GetStar()
	if star >= 3 {
		// 3星解锁特殊技能
		if skillTemplate, exists := ps.skillTemplates["special_skill_1"]; exists {
			skill := ps.createSkillFromTemplate(skillTemplate)
			pet.AddSkill(skill)
		}
	}
	if star >= 5 {
		// 5星解锁终极技能
		if skillTemplate, exists := ps.skillTemplates["ultimate_skill_1"]; exists {
			skill := ps.createSkillFromTemplate(skillTemplate)
			pet.AddSkill(skill)
		}
	}
}

// calculateTrainingResult 计算训练结果
func (ps *PetService) calculateTrainingResult(pet *PetAggregate, trainingType TrainingType) *TrainingResult {
	result := &TrainingResult{
		TrainingType: trainingType,
		Timestamp:    time.Now(),
	}
	
	switch trainingType {
	case TrainingTypeExperience:
		result.ExperienceGain = TrainingExperienceGain
	case TrainingTypeAttribute:
		result.AttributeGain = TrainingAttributeGain
	case TrainingTypeSkill:
		result.SkillExpGain = TrainingSkillExpGain
	}
	
	return result
}

// calculateAttributeChanges 计算属性变化
func (ps *PetService) calculateAttributeChanges(before, after *PetAttributes) map[string]int64 {
	changes := make(map[string]int64)
	
	changes["health"] = after.GetHealth() - before.GetHealth()
	changes["attack"] = after.GetAttack() - before.GetAttack()
	changes["defense"] = after.GetDefense() - before.GetDefense()
	changes["speed"] = after.GetSpeed() - before.GetSpeed()
	changes["critical"] = after.GetCritical() - before.GetCritical()
	changes["hit"] = after.GetHit() - before.GetHit()
	changes["dodge"] = after.GetDodge() - before.GetDodge()
	
	return changes
}

// validateBreedingConditions 验证繁殖条件
func (ps *PetService) validateBreedingConditions(parent1, parent2 *PetAggregate) error {
	if parent1.GetPlayerID() != parent2.GetPlayerID() {
		return ErrDifferentOwners
	}
	
	if parent1.GetLevel() < MinBreedingLevel || parent2.GetLevel() < MinBreedingLevel {
		return ErrInsufficientBreedingLevel
	}
	
	if !parent1.IsAlive() || !parent2.IsAlive() {
		return ErrPetIsDead
	}
	
	if !parent1.IsIdle() || !parent2.IsIdle() {
		return ErrPetNotIdle
	}
	
	return nil
}

// generateOffspring 生成后代
func (ps *PetService) generateOffspring(parent1, parent2 *PetAggregate) *PetAggregate {
	// 简化实现：随机选择一个父母的配置
	var configID uint32
	if rand.Float32() < 0.5 {
		configID = parent1.GetConfigID()
	} else {
		configID = parent2.GetConfigID()
	}
	
	// 创建后代
	offspring, _ := ps.CreatePet(parent1.GetPlayerID(), configID, "")
	
	// 继承部分属性
	ps.inheritAttributes(offspring, parent1, parent2)
	
	return offspring
}

// inheritAttributes 继承属性
func (ps *PetService) inheritAttributes(offspring, parent1, parent2 *PetAggregate) {
	offspringAttrs := offspring.GetAttributes()
	parent1Attrs := parent1.GetAttributes()
	parent2Attrs := parent2.GetAttributes()
	
	// 取父母属性的平均值作为基础
	avgHealth := (parent1Attrs.GetHealth() + parent2Attrs.GetHealth()) / 2
	avgAttack := (parent1Attrs.GetAttack() + parent2Attrs.GetAttack()) / 2
	avgDefense := (parent1Attrs.GetDefense() + parent2Attrs.GetDefense()) / 2
	avgSpeed := (parent1Attrs.GetSpeed() + parent2Attrs.GetSpeed()) / 2
	
	// 应用继承属性（有一定随机性）
	offspringAttrs.AddHealth(int64(float64(avgHealth) * (0.8 + rand.Float64()*0.4)))
	offspringAttrs.AddAttack(int64(float64(avgAttack) * (0.8 + rand.Float64()*0.4)))
	offspringAttrs.AddDefense(int64(float64(avgDefense) * (0.8 + rand.Float64()*0.4)))
	offspringAttrs.AddSpeed(int64(float64(avgSpeed) * (0.8 + rand.Float64()*0.4)))
}

// calculateSkillPowerBonus 计算技能战力加成
func (ps *PetService) calculateSkillPowerBonus(skills []*PetSkill) int64 {
	var bonus int64
	for _, skill := range skills {
		bonus += skill.GetDamage() * int64(skill.GetLevel())
	}
	return bonus
}

// getTemplatesByRarity 根据稀有度获取模板
func (ps *PetService) getTemplatesByRarity(rarity PetRarity) []*PetTemplate {
	var templates []*PetTemplate
	for _, template := range ps.petTemplates {
		if template.Rarity == rarity {
			templates = append(templates, template)
		}
	}
	return templates
}

// applyRandomBonus 应用随机加成
func (ps *PetService) applyRandomBonus(pet *PetAggregate, rarity PetRarity) {
	multiplier := rarity.GetAttributeMultiplier()
	attributes := pet.GetAttributes()
	
	// 随机加成范围：0.8-1.2倍
	randomFactor := 0.8 + rand.Float64()*0.4
	finalMultiplier := multiplier * randomFactor
	
	bonusHealth := int64(float64(attributes.GetHealth()) * (finalMultiplier - 1.0))
	bonusAttack := int64(float64(attributes.GetAttack()) * (finalMultiplier - 1.0))
	bonusDefense := int64(float64(attributes.GetDefense()) * (finalMultiplier - 1.0))
	bonusSpeed := int64(float64(attributes.GetSpeed()) * (finalMultiplier - 1.0))
	
	attributes.AddHealth(bonusHealth)
	attributes.AddAttack(bonusAttack)
	attributes.AddDefense(bonusDefense)
	attributes.AddSpeed(bonusSpeed)
}

// 模板结构体定义

// PetTemplate 宠物模板
type PetTemplate struct {
	ConfigID      uint32
	Name          string
	DefaultName   string
	Category      PetCategory
	Rarity        PetRarity
	Size          PetSize
	Gender        PetGender
	Personality   PetPersonality
	BaseHealth    int64
	BaseAttack    int64
	BaseDefense   int64
	BaseSpeed     int64
	BaseCritical  int64
	BaseHit       int64
	BaseDodge     int64
	InitialSkills []string
	Description   string
}

// SkillTemplate 技能模板
type SkillTemplate struct {
	SkillID    string
	Name       string
	Type       SkillType
	BaseDamage int64
	Cooldown   time.Duration
	Effects    []SkillEffectTemplate
	Description string
}

// SkillEffectTemplate 技能效果模板
type SkillEffectTemplate struct {
	Type     string
	Value    float64
	Duration time.Duration
}

// SkinTemplate 皮肤模板
type SkinTemplate struct {
	SkinID         string
	Name           string
	Rarity         PetRarity
	PowerBonus     int64
	AttributeBonus map[string]float64
	UnlockCondition string
	Description    string
}

// BondTemplate 羁绊模板
type BondTemplate struct {
	BondID      string
	Name        string
	Description string
	PowerBonus  int64
	Attributes  map[string]float64
	RequiredPets []uint32
}

// FragmentTemplate 碎片模板
type FragmentTemplate struct {
	FragmentID       uint32
	Name             string
	RelatedPetID     uint32
	RequiredQuantity uint64
	Description      string
}

// 结果结构体定义

// TrainingResult 训练结果
type TrainingResult struct {
	TrainingType   TrainingType
	ExperienceGain uint64
	AttributeGain  int64
	SkillExpGain   uint64
	Timestamp      time.Time
}

// FeedingResult 喂食结果
type FeedingResult struct {
	FoodType         FoodType
	Amount           int
	ExperienceGain   uint64
	LevelUp          bool
	AttributeChanges map[string]int64
	Timestamp        time.Time
}

// 常量定义
const (
	MinBreedingLevel = 10
)