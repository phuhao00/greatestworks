package team

import (
	"context"
	"fmt"
)

// TeamService 队伍领域服务
type TeamService struct {
	teamRepo TeamRepository
}

// NewTeamService 创建队伍服务
func NewTeamService(teamRepo TeamRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
	}
}

// CreateTeam 创建队伍
func (s *TeamService) CreateTeam(ctx context.Context, name, leaderID string, maxMembers int, isPublic bool, password string) (*Team, error) {
	// 验证队伍名称
	if err := s.validateTeamName(name); err != nil {
		return nil, err
	}

	// 检查玩家是否已有队伍
	hasTeam, err := s.teamRepo.PlayerHasTeam(ctx, leaderID)
	if err != nil {
		return nil, fmt.Errorf("检查玩家队伍状态失败: %w", err)
	}
	if hasTeam {
		return nil, ErrPlayerAlreadyInTeam
	}

	// 创建队伍
	teamID := generateTeamID()
	team := NewTeam(teamID, name, leaderID, maxMembers, isPublic)
	if password != "" {
		team.SetPassword(password)
	}

	// 添加队长为成员
	leader := NewTeamMember(leaderID, "", 1)
	leader.Role = TeamRoleLeader
	if err := team.AddMember(leader); err != nil {
		return nil, fmt.Errorf("添加队长失败: %w", err)
	}

	// 保存队伍
	if err := s.teamRepo.SaveTeam(ctx, team); err != nil {
		return nil, fmt.Errorf("保存队伍失败: %w", err)
	}

	return team, nil
}

// JoinTeam 加入队伍
func (s *TeamService) JoinTeam(ctx context.Context, teamID, playerID, nickname string, level int, password string) error {
	// 获取队伍
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		return fmt.Errorf("获取队伍失败: %w", err)
	}
	if team == nil {
		return ErrTeamNotFound
	}

	// 检查密码
	if !team.IsPublic && team.Password != password {
		return ErrInvalidPassword
	}

	// 检查队伍是否已满
	if team.IsFull() {
		return ErrTeamFull
	}

	// 检查玩家是否已有队伍
	hasTeam, err := s.teamRepo.PlayerHasTeam(ctx, playerID)
	if err != nil {
		return fmt.Errorf("检查玩家队伍状态失败: %w", err)
	}
	if hasTeam {
		return ErrPlayerAlreadyInTeam
	}

	// 创建新成员
	member := NewTeamMember(playerID, nickname, level)

	// 添加成员到队伍
	if err := team.AddMember(member); err != nil {
		return err
	}

	// 保存队伍
	if err := s.teamRepo.SaveTeam(ctx, team); err != nil {
		return fmt.Errorf("保存队伍失败: %w", err)
	}

	return nil
}

// LeaveTeam 离开队伍
func (s *TeamService) LeaveTeam(ctx context.Context, teamID, playerID string) error {
	// 获取队伍
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		return fmt.Errorf("获取队伍失败: %w", err)
	}
	if team == nil {
		return ErrTeamNotFound
	}

	// 如果是队长离开且队伍还有其他成员，需要先转让队长
	if team.LeaderID == playerID && team.GetMemberCount() > 1 {
		return ErrLeaderMustTransfer
	}

	// 移除成员
	if err := team.RemoveMember(playerID); err != nil {
		return err
	}

	// 如果队伍为空，删除队伍
	if team.GetMemberCount() == 0 {
		if err := s.teamRepo.DeleteTeam(ctx, teamID); err != nil {
			return fmt.Errorf("删除队伍失败: %w", err)
		}
	} else {
		// 保存队伍
		if err := s.teamRepo.SaveTeam(ctx, team); err != nil {
			return fmt.Errorf("保存队伍失败: %w", err)
		}
	}

	return nil
}

// validateTeamName 验证队伍名称
func (s *TeamService) validateTeamName(name string) error {
	if len(name) < 2 {
		return ErrTeamNameTooShort
	}
	if len(name) > 20 {
		return ErrTeamNameTooLong
	}
	return nil
}

// generateTeamID 生成队伍ID
func generateTeamID() string {
	return "team_" + randomString(16)
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
