package family

import (
	"context"
	"fmt"
)

// FamilyService 家族领域服务
type FamilyService struct {
	familyRepo FamilyRepository
}

// NewFamilyService 创建家族服务
func NewFamilyService(familyRepo FamilyRepository) *FamilyService {
	return &FamilyService{
		familyRepo: familyRepo,
	}
}

// CreateFamily 创建家族
func (s *FamilyService) CreateFamily(ctx context.Context, name, description, leaderID string) (*Family, error) {
	// 验证家族名称
	if err := s.validateFamilyName(name); err != nil {
		return nil, err
	}

	// 检查名称是否已存在
	exists, err := s.familyRepo.FamilyExistsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("检查家族名称失败: %w", err)
	}
	if exists {
		return nil, ErrFamilyNameExists
	}

	// 检查玩家是否已有家族
	hasFamily, err := s.familyRepo.PlayerHasFamily(ctx, leaderID)
	if err != nil {
		return nil, fmt.Errorf("检查玩家家族状态失败: %w", err)
	}
	if hasFamily {
		return nil, ErrPlayerAlreadyInFamily
	}

	// 创建家族
	familyID := generateFamilyID()
	family := NewFamily(familyID, name, description, leaderID)

	// 添加族长为成员
	leader := NewFamilyMember(leaderID, "")
	leader.Role = FamilyRoleLeader
	if err := family.AddMember(leader); err != nil {
		return nil, fmt.Errorf("添加族长失败: %w", err)
	}

	// 保存家族
	if err := s.familyRepo.SaveFamily(ctx, family); err != nil {
		return nil, fmt.Errorf("保存家族失败: %w", err)
	}

	return family, nil
}

// JoinFamily 加入家族
func (s *FamilyService) JoinFamily(ctx context.Context, familyID, playerID, nickname string) error {
	// 获取家族
	family, err := s.familyRepo.GetFamilyByID(ctx, familyID)
	if err != nil {
		return fmt.Errorf("获取家族失败: %w", err)
	}
	if family == nil {
		return ErrFamilyNotFound
	}

	// 检查家族是否已满
	if family.IsFull() {
		return ErrFamilyFull
	}

	// 检查玩家是否已有家族
	hasFamily, err := s.familyRepo.PlayerHasFamily(ctx, playerID)
	if err != nil {
		return fmt.Errorf("检查玩家家族状态失败: %w", err)
	}
	if hasFamily {
		return ErrPlayerAlreadyInFamily
	}

	// 创建新成员
	member := NewFamilyMember(playerID, nickname)

	// 添加成员到家族
	if err := family.AddMember(member); err != nil {
		return err
	}

	// 保存家族
	if err := s.familyRepo.SaveFamily(ctx, family); err != nil {
		return fmt.Errorf("保存家族失败: %w", err)
	}

	return nil
}

// LeaveFamily 离开家族
func (s *FamilyService) LeaveFamily(ctx context.Context, familyID, playerID string) error {
	// 获取家族
	family, err := s.familyRepo.GetFamilyByID(ctx, familyID)
	if err != nil {
		return fmt.Errorf("获取家族失败: %w", err)
	}
	if family == nil {
		return ErrFamilyNotFound
	}

	// 族长不能直接离开，需要先转让族长
	if family.LeaderID == playerID {
		return ErrLeaderCannotLeave
	}

	// 移除成员
	if err := family.RemoveMember(playerID); err != nil {
		return err
	}

	// 保存家族
	if err := s.familyRepo.SaveFamily(ctx, family); err != nil {
		return fmt.Errorf("保存家族失败: %w", err)
	}

	return nil
}

// validateFamilyName 验证家族名称
func (s *FamilyService) validateFamilyName(name string) error {
	if len(name) < 2 {
		return ErrFamilyNameTooShort
	}
	if len(name) > 20 {
		return ErrFamilyNameTooLong
	}
	return nil
}

// generateFamilyID 生成家族ID
func generateFamilyID() string {
	return "fam_" + randomString(16)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
