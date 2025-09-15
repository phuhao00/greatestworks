package family

import (
	"time"
)

// Family 家族聚合根
type Family struct {
	ID          string
	Name        string
	Description string
	Level       int
	Experience  int64
	LeaderID    string
	members     map[string]*FamilyMember
	MaxMembers  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int64
}

// NewFamily 创建新家族
func NewFamily(id, name, description, leaderID string) *Family {
	return &Family{
		ID:          id,
		Name:        name,
		Description: description,
		Level:       1,
		Experience:  0,
		LeaderID:    leaderID,
		members:     make(map[string]*FamilyMember),
		MaxMembers:  DefaultMaxMembers,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Version:     1,
	}
}

// AddMember 添加成员
func (f *Family) AddMember(member *FamilyMember) error {
	if len(f.members) >= f.MaxMembers {
		return ErrFamilyFull
	}
	
	if _, exists := f.members[member.PlayerID]; exists {
		return ErrMemberAlreadyExists
	}
	
	f.members[member.PlayerID] = member
	f.UpdatedAt = time.Now()
	f.Version++
	
	return nil
}

// RemoveMember 移除成员
func (f *Family) RemoveMember(playerID string) error {
	if playerID == f.LeaderID {
		return ErrCannotRemoveLeader
	}
	
	if _, exists := f.members[playerID]; !exists {
		return ErrMemberNotFound
	}
	
	delete(f.members, playerID)
	f.UpdatedAt = time.Now()
	f.Version++
	
	return nil
}

// PromoteMember 提升成员职位
func (f *Family) PromoteMember(playerID string, newRole FamilyRole) error {
	member, exists := f.members[playerID]
	if !exists {
		return ErrMemberNotFound
	}
	
	if newRole == FamilyRoleLeader {
		return ErrCannotPromoteToLeader
	}
	
	member.Role = newRole
	f.UpdatedAt = time.Now()
	f.Version++
	
	return nil
}

// TransferLeadership 转让族长
func (f *Family) TransferLeadership(newLeaderID string) error {
	newLeader, exists := f.members[newLeaderID]
	if !exists {
		return ErrMemberNotFound
	}
	
	// 将原族长降为副族长
	if oldLeader, exists := f.members[f.LeaderID]; exists {
		oldLeader.Role = FamilyRoleViceLeader
	}
	
	// 设置新族长
	f.LeaderID = newLeaderID
	newLeader.Role = FamilyRoleLeader
	f.UpdatedAt = time.Now()
	f.Version++
	
	return nil
}

// AddExperience 增加经验
func (f *Family) AddExperience(exp int64) {
	f.Experience += exp
	f.checkLevelUp()
	f.UpdatedAt = time.Now()
	f.Version++
}

// checkLevelUp 检查升级
func (f *Family) checkLevelUp() {
	requiredExp := f.getRequiredExperience(f.Level + 1)
	if f.Experience >= requiredExp {
		f.Level++
		f.MaxMembers += MembersPerLevel
		f.checkLevelUp() // 递归检查是否可以继续升级
	}
}

// getRequiredExperience 获取升级所需经验
func (f *Family) getRequiredExperience(level int) int64 {
	return int64(level * level * 1000)
}

// GetMembers 获取所有成员
func (f *Family) GetMembers() []*FamilyMember {
	members := make([]*FamilyMember, 0, len(f.members))
	for _, member := range f.members {
		members = append(members, member)
	}
	return members
}

// GetMemberCount 获取成员数量
func (f *Family) GetMemberCount() int {
	return len(f.members)
}

// IsFull 检查是否已满
func (f *Family) IsFull() bool {
	return len(f.members) >= f.MaxMembers
}

const (
	DefaultMaxMembers = 20  // 默认最大成员数
	MembersPerLevel   = 5