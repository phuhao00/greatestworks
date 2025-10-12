package family

import "time"

// FamilyMember 家族成员实体
type FamilyMember struct {
	PlayerID     string
	Nickname     string
	Role         FamilyRole
	Contribution int64
	JoinedAt     time.Time
	LastActive   time.Time
	Level        int
}

// FamilyRole 家族职位
type FamilyRole int

const (
	FamilyRoleMember     FamilyRole = iota // 普通成员
	FamilyRoleElite                        // 精英
	FamilyRoleViceLeader                   // 副族长
	FamilyRoleLeader                       // 族长
)

// NewFamilyMember 创建新的家族成员
func NewFamilyMember(playerID, nickname string) *FamilyMember {
	return &FamilyMember{
		PlayerID:     playerID,
		Nickname:     nickname,
		Role:         FamilyRoleMember,
		Contribution: 0,
		JoinedAt:     time.Now(),
		LastActive:   time.Now(),
		Level:        1,
	}
}

// AddContribution 增加贡献度
func (m *FamilyMember) AddContribution(amount int64) {
	m.Contribution += amount
}

// UpdateLastActive 更新最后活跃时间
func (m *FamilyMember) UpdateLastActive() {
	m.LastActive = time.Now()
}

// GetMembershipDuration 获取加入时长
func (m *FamilyMember) GetMembershipDuration() time.Duration {
	return time.Since(m.JoinedAt)
}

// IsLeader 是否为族长
func (m *FamilyMember) IsLeader() bool {
	return m.Role == FamilyRoleLeader
}

// IsViceLeader 是否为副族长
func (m *FamilyMember) IsViceLeader() bool {
	return m.Role == FamilyRoleViceLeader
}

// HasManagementPermission 是否有管理权限
func (m *FamilyMember) HasManagementPermission() bool {
	return m.Role >= FamilyRoleViceLeader
}

// CanKickMembers 是否可以踢出成员
func (m *FamilyMember) CanKickMembers() bool {
	return m.Role >= FamilyRoleViceLeader
}

// CanPromoteMembers 是否可以提升成员
func (m *FamilyMember) CanPromoteMembers() bool {
	return m.Role == FamilyRoleLeader
}
