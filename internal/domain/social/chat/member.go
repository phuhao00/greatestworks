package chat

import (
	"time"
)

// Member 聊天频道成员实体
type Member struct {
	PlayerID    string
	Nickname    string
	Role        MemberRole
	JoinedAt    time.Time
	LastActive  time.Time
	Permissions []Permission
	IsMuted     bool
	MutedUntil  *time.Time
}

// MemberRole 成员角色
type MemberRole int

const (
	MemberRoleNormal    MemberRole = iota // 普通成员
	MemberRoleModerator                   // 管理员
	MemberRoleOwner                       // 频道所有者
)

// Permission 权限
type Permission int

const (
	PermissionSendMessage   Permission = iota // 发送消息
	PermissionDeleteMessage                   // 删除消息
	PermissionMuteMembers                     // 禁言成员
	PermissionKickMembers                     // 踢出成员
	PermissionManageChannel                   // 管理频道
)

// NewMember 创建新成员
func NewMember(playerID, nickname string) *Member {
	return &Member{
		PlayerID:    playerID,
		Nickname:    nickname,
		Role:        MemberRoleNormal,
		JoinedAt:    time.Now(),
		LastActive:  time.Now(),
		Permissions: getDefaultPermissions(MemberRoleNormal),
		IsMuted:     false,
	}
}

// HasPermission 检查是否有权限
func (m *Member) HasPermission(permission Permission) bool {
	for _, p := range m.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// SetRole 设置角色
func (m *Member) SetRole(role MemberRole) {
	m.Role = role
	m.Permissions = getDefaultPermissions(role)
}

// Mute 禁言成员
func (m *Member) Mute(duration time.Duration) {
	m.IsMuted = true
	mutedUntil := time.Now().Add(duration)
	m.MutedUntil = &mutedUntil
}

// Unmute 解除禁言
func (m *Member) Unmute() {
	m.IsMuted = false
	m.MutedUntil = nil
}

// IsCurrentlyMuted 检查当前是否被禁言
func (m *Member) IsCurrentlyMuted() bool {
	if !m.IsMuted {
		return false
	}

	if m.MutedUntil != nil && time.Now().After(*m.MutedUntil) {
		m.Unmute()
		return false
	}

	return true
}

// UpdateLastActive 更新最后活跃时间
func (m *Member) UpdateLastActive() {
	m.LastActive = time.Now()
}

// GetMembershipDuration 获取加入时长
func (m *Member) GetMembershipDuration() time.Duration {
	return time.Since(m.JoinedAt)
}

// CanSendMessage 检查是否可以发送消息
func (m *Member) CanSendMessage() bool {
	return m.HasPermission(PermissionSendMessage) && !m.IsCurrentlyMuted()
}

// getDefaultPermissions 获取角色默认权限
func getDefaultPermissions(role MemberRole) []Permission {
	switch role {
	case MemberRoleNormal:
		return []Permission{PermissionSendMessage}
	case MemberRoleModerator:
		return []Permission{
			PermissionSendMessage,
			PermissionDeleteMessage,
			PermissionMuteMembers,
		}
	case MemberRoleOwner:
		return []Permission{
			PermissionSendMessage,
			PermissionDeleteMessage,
			PermissionMuteMembers,
			PermissionKickMembers,
			PermissionManageChannel,
		}
	default:
		return []Permission{PermissionSendMessage}
	}
}
