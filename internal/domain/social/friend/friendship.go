package friend

import (
	"time"
)

// Friendship 好友关系聚合根
type Friendship struct {
	ID          string
	PlayerID    string
	FriendID    string
	Status      FriendshipStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RequestedBy string // 发起请求的玩家ID
	Notes       string // 备注
	Version     int64
}

// FriendshipStatus 好友关系状态
type FriendshipStatus int

const (
	FriendshipStatusPending  FriendshipStatus = iota // 待确认
	FriendshipStatusAccepted                         // 已接受
	FriendshipStatusBlocked                          // 已屏蔽
	FriendshipStatusDeleted                          // 已删除
)

// NewFriendship 创建新的好友关系
func NewFriendship(playerID, friendID, requestedBy string) *Friendship {
	return &Friendship{
		ID:          generateFriendshipID(),
		PlayerID:    playerID,
		FriendID:    friendID,
		Status:      FriendshipStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		RequestedBy: requestedBy,
		Version:     1,
	}
}

// Accept 接受好友请求
func (f *Friendship) Accept() error {
	if f.Status != FriendshipStatusPending {
		return ErrInvalidFriendshipStatus
	}

	f.Status = FriendshipStatusAccepted
	f.UpdatedAt = time.Now()
	f.Version++

	return nil
}

// Block 屏蔽好友
func (f *Friendship) Block() error {
	if f.Status == FriendshipStatusDeleted {
		return ErrFriendshipDeleted
	}

	f.Status = FriendshipStatusBlocked
	f.UpdatedAt = time.Now()
	f.Version++

	return nil
}

// Unblock 解除屏蔽
func (f *Friendship) Unblock() error {
	if f.Status != FriendshipStatusBlocked {
		return ErrNotBlocked
	}

	f.Status = FriendshipStatusAccepted
	f.UpdatedAt = time.Now()
	f.Version++

	return nil
}

// Delete 删除好友关系
func (f *Friendship) Delete() error {
	if f.Status == FriendshipStatusDeleted {
		return ErrFriendshipAlreadyDeleted
	}

	f.Status = FriendshipStatusDeleted
	f.UpdatedAt = time.Now()
	f.Version++

	return nil
}

// SetNotes 设置备注
func (f *Friendship) SetNotes(notes string) {
	f.Notes = notes
	f.UpdatedAt = time.Now()
	f.Version++
}

// IsActive 检查好友关系是否活跃
func (f *Friendship) IsActive() bool {
	return f.Status == FriendshipStatusAccepted
}

// IsPending 检查是否为待确认状态
func (f *Friendship) IsPending() bool {
	return f.Status == FriendshipStatusPending
}

// IsBlocked 检查是否被屏蔽
func (f *Friendship) IsBlocked() bool {
	return f.Status == FriendshipStatusBlocked
}

// IsDeleted 检查是否已删除
func (f *Friendship) IsDeleted() bool {
	return f.Status == FriendshipStatusDeleted
}

// GetOtherPlayerID 获取对方玩家ID
func (f *Friendship) GetOtherPlayerID(currentPlayerID string) string {
	if f.PlayerID == currentPlayerID {
		return f.FriendID
	}
	return f.PlayerID
}

// GetDuration 获取好友关系持续时间
func (f *Friendship) GetDuration() time.Duration {
	return time.Since(f.CreatedAt)
}

// GetVersion 获取版本号
func (f *Friendship) GetVersion() int64 {
	return f.Version
}

// generateFriendshipID 生成好友关系ID
func generateFriendshipID() string {
	return "fs_" + randomString(16)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
