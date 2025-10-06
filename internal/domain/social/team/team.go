package team

import (
	"time"
)

// Team 队伍聚合根
type Team struct {
	ID         string
	Name       string
	LeaderID   string
	members    map[string]*TeamMember
	MaxMembers int
	IsPublic   bool
	Password   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Version    int64
}

// NewTeam 创建新队伍
func NewTeam(id, name, leaderID string, maxMembers int, isPublic bool) *Team {
	return &Team{
		ID:         id,
		Name:       name,
		LeaderID:   leaderID,
		members:    make(map[string]*TeamMember),
		MaxMembers: maxMembers,
		IsPublic:   isPublic,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Version:    1,
	}
}

// AddMember 添加成员
func (t *Team) AddMember(member *TeamMember) error {
	if len(t.members) >= t.MaxMembers {
		return ErrTeamFull
	}

	if _, exists := t.members[member.PlayerID]; exists {
		return ErrMemberAlreadyExists
	}

	t.members[member.PlayerID] = member
	t.UpdatedAt = time.Now()
	t.Version++

	return nil
}

// RemoveMember 移除成员
func (t *Team) RemoveMember(playerID string) error {
	if playerID == t.LeaderID {
		return ErrCannotRemoveLeader
	}

	if _, exists := t.members[playerID]; !exists {
		return ErrMemberNotFound
	}

	delete(t.members, playerID)
	t.UpdatedAt = time.Now()
	t.Version++

	return nil
}

// TransferLeadership 转让队长
func (t *Team) TransferLeadership(newLeaderID string) error {
	newLeader, exists := t.members[newLeaderID]
	if !exists {
		return ErrMemberNotFound
	}

	// 将原队长降为普通成员
	if oldLeader, exists := t.members[t.LeaderID]; exists {
		oldLeader.Role = TeamRoleMember
	}

	// 设置新队长
	t.LeaderID = newLeaderID
	newLeader.Role = TeamRoleLeader
	t.UpdatedAt = time.Now()
	t.Version++

	return nil
}

// SetPassword 设置密码
func (t *Team) SetPassword(password string) {
	t.Password = password
	t.IsPublic = password == ""
	t.UpdatedAt = time.Now()
	t.Version++
}

// GetMembers 获取所有成员
func (t *Team) GetMembers() []*TeamMember {
	members := make([]*TeamMember, 0, len(t.members))
	for _, member := range t.members {
		members = append(members, member)
	}
	return members
}

// GetMemberCount 获取成员数量
func (t *Team) GetMemberCount() int {
	return len(t.members)
}

// IsFull 检查是否已满
func (t *Team) IsFull() bool {
	return len(t.members) >= t.MaxMembers
}

// IsMember 检查是否为成员
func (t *Team) IsMember(playerID string) bool {
	_, exists := t.members[playerID]
	return exists
}
