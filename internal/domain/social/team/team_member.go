package team

import "time"

// TeamMember 队伍成员实体
type TeamMember struct {
	PlayerID   string
	Nickname   string
	Role       TeamRole
	JoinedAt   time.Time
	LastActive time.Time
	Level      int
	IsReady    bool
}

// TeamRole 队伍角色
type TeamRole int

const (
	TeamRoleMember TeamRole = iota // 普通成员
	TeamRoleLeader                 // 队长
)

// NewTeamMember 创建新的队伍成员
func NewTeamMember(playerID, nickname string, level int) *TeamMember {
	return &TeamMember{
		PlayerID:   playerID,
		Nickname:   nickname,
		Role:       TeamRoleMember,
		JoinedAt:   time.Now(),
		LastActive: time.Now(),
		Level:      level,
		IsReady:    false,
	}
}

// SetReady 设置准备状态
func (m *TeamMember) SetReady(ready bool) {
	m.IsReady = ready
}

// UpdateLastActive 更新最后活跃时间
func (m *TeamMember) UpdateLastActive() {
	m.LastActive = time.Now()
}

// IsLeader 是否为队长
func (m *TeamMember) IsLeader() bool {
	return m.Role == TeamRoleLeader
}

// GetMembershipDuration 获取加入时长
func (m *TeamMember) GetMembershipDuration() time.Duration {
	return time.Since(m.JoinedAt)
}
