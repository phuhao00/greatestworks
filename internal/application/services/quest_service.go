package services

import (
	"context"
	"errors"
	"fmt"

	"greatestworks/internal/domain/quest"
	"greatestworks/internal/infrastructure/datamanager"
	"greatestworks/internal/infrastructure/persistence"
)

// QuestService 任务服务
type QuestService struct {
	questRepo *persistence.QuestRepository
}

// NewQuestService 创建任务服务
func NewQuestService(questRepo *persistence.QuestRepository) *QuestService {
	return &QuestService{
		questRepo: questRepo,
	}
}

// AcceptQuest 接受任务
func (s *QuestService) AcceptQuest(ctx context.Context, characterID int64, questID int32) error {
	// 获取任务配置
	questDefine := datamanager.GetInstance().GetQuest(questID)
	if questDefine == nil {
		return errors.New("quest not found")
	}

	// 检查前置任务
	// TODO: 实现前置任务检查

	// 创建任务进度
	objectives := make([]persistence.DbObjective, len(questDefine.Objectives))
	for i, obj := range questDefine.Objectives {
		objectives[i] = persistence.DbObjective{
			Type:     obj.Type,
			TargetID: obj.TargetID,
			Required: obj.Required,
			Current:  0,
		}
	}

	dbQuest := &persistence.DbQuest{
		CharacterID: characterID,
		QuestID:     questID,
		Status:      0, // 进行中
		Objectives:  objectives,
	}

	if err := s.questRepo.Create(ctx, dbQuest); err != nil {
		return fmt.Errorf("failed to accept quest: %w", err)
	}

	return nil
}

// GetQuests 获取角色的任务列表
func (s *QuestService) GetQuests(ctx context.Context, characterID int64) ([]*persistence.DbQuest, error) {
	return s.questRepo.FindByCharacterID(ctx, characterID)
}

// UpdateObjective 更新任务目标
func (s *QuestService) UpdateObjective(ctx context.Context, characterID int64, questID, objType, targetID, progress int32) error {
	quests, err := s.questRepo.FindByCharacterID(ctx, characterID)
	if err != nil {
		return err
	}

	// 查找对应任务
	var targetQuest *persistence.DbQuest
	for _, q := range quests {
		if q.QuestID == questID && q.Status == 0 {
			targetQuest = q
			break
		}
	}

	if targetQuest == nil {
		return errors.New("quest not found or already completed")
	}

	// 更新目标进度
	updated := false
	for i := range targetQuest.Objectives {
		obj := &targetQuest.Objectives[i]
		if obj.Type == objType && obj.TargetID == targetID {
			obj.Current += progress
			if obj.Current > obj.Required {
				obj.Current = obj.Required
			}
			updated = true
		}
	}

	if !updated {
		return errors.New("objective not found")
	}

	// 检查是否完成
	allComplete := true
	for _, obj := range targetQuest.Objectives {
		if obj.Current < obj.Required {
			allComplete = false
			break
		}
	}

	if allComplete {
		targetQuest.Status = 1 // 已完成
	}

	return s.questRepo.Update(ctx, targetQuest)
}

// SubmitQuest 提交任务
func (s *QuestService) SubmitQuest(ctx context.Context, characterID int64, questID int32) error {
	quests, err := s.questRepo.FindByCharacterID(ctx, characterID)
	if err != nil {
		return err
	}

	// 查找对应任务
	var targetQuest *persistence.DbQuest
	for _, q := range quests {
		if q.QuestID == questID && q.Status == 1 {
			targetQuest = q
			break
		}
	}

	if targetQuest == nil {
		return errors.New("quest not completed")
	}

	// 获取任务配置
	questDefine := datamanager.GetInstance().GetQuest(questID)
	if questDefine == nil {
		return errors.New("quest not found")
	}

	// TODO: 发放奖励
	// - 经验
	// - 金币
	// - 物品

	// 标记为已领取
	targetQuest.Status = 2

	return s.questRepo.Update(ctx, targetQuest)
}

// AbandonQuest 放弃任务
func (s *QuestService) AbandonQuest(ctx context.Context, characterID int64, questID int32) error {
	quests, err := s.questRepo.FindByCharacterID(ctx, characterID)
	if err != nil {
		return err
	}

	// 查找对应任务
	var targetQuest *persistence.DbQuest
	for _, q := range quests {
		if q.QuestID == questID && q.Status == 0 {
			targetQuest = q
			break
		}
	}

	if targetQuest == nil {
		return errors.New("quest not found or already completed")
	}

	// 删除任务进度（这里简化处理，实际可能需要软删除）
	// TODO: 实现任务删除
	_ = targetQuest

	return nil
}

// OnKill 击杀事件处理
func (s *QuestService) OnKill(ctx context.Context, characterID int64, targetID int32) error {
	return s.UpdateObjective(ctx, characterID, 0, int32(quest.ObjectiveTypeKill), targetID, 1)
}

// OnCollect 收集事件处理
func (s *QuestService) OnCollect(ctx context.Context, characterID int64, itemID, count int32) error {
	return s.UpdateObjective(ctx, characterID, 0, int32(quest.ObjectiveTypeCollect), itemID, count)
}

// OnReach 到达事件处理
func (s *QuestService) OnReach(ctx context.Context, characterID int64, locationID int32) error {
	return s.UpdateObjective(ctx, characterID, 0, int32(quest.ObjectiveTypeReach), locationID, 1)
}

// OnTalk 对话事件处理
func (s *QuestService) OnTalk(ctx context.Context, characterID int64, npcID int32) error {
	return s.UpdateObjective(ctx, characterID, 0, int32(quest.ObjectiveTypeTalk), npcID, 1)
}
