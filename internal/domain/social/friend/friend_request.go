package friend

import (
	"time"
)

// FriendRequest 好友请求实体
type FriendRequest struct {
	ID           string
	FromPlayerID string
	ToPlayerID   string
	Message      string
	Status       RequestStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ExpiresAt    time.Time
}

// RequestStatus 请求状态
type RequestStatus int

const (
	RequestStatusPending  RequestStatus = iota // 待处理
	RequestStatusAccepted                      // 已接受
	RequestStatusRejected                      // 已拒绝
	RequestStatusExpired                       // 已过期
	RequestStatusCanceled                      // 已取消
)

// NewFriendRequest 创建新的好友请求
func NewFriendRequest(fromPlayerID, toPlayerID, message string) *FriendRequest {
	return &FriendRequest{
		ID:           generateRequestID(),
		FromPlayerID: fromPlayerID,
		ToPlayerID:   toPlayerID,
		Message:      message,
		Status:       RequestStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(DefaultRequestExpiration),
	}
}

// Accept 接受好友请求
func (r *FriendRequest) Accept() error {
	if r.Status != RequestStatusPending {
		return ErrInvalidRequestStatus
	}

	if r.IsExpired() {
		return ErrRequestExpired
	}

	r.Status = RequestStatusAccepted
	r.UpdatedAt = time.Now()

	return nil
}

// Reject 拒绝好友请求
func (r *FriendRequest) Reject() error {
	if r.Status != RequestStatusPending {
		return ErrInvalidRequestStatus
	}

	if r.IsExpired() {
		return ErrRequestExpired
	}

	r.Status = RequestStatusRejected
	r.UpdatedAt = time.Now()

	return nil
}

// Cancel 取消好友请求
func (r *FriendRequest) Cancel() error {
	if r.Status != RequestStatusPending {
		return ErrInvalidRequestStatus
	}

	r.Status = RequestStatusCanceled
	r.UpdatedAt = time.Now()

	return nil
}

// IsExpired 检查请求是否已过期
func (r *FriendRequest) IsExpired() bool {
	if r.Status == RequestStatusExpired {
		return true
	}

	if time.Now().After(r.ExpiresAt) {
		r.Status = RequestStatusExpired
		r.UpdatedAt = time.Now()
		return true
	}

	return false
}

// IsPending 检查请求是否待处理
func (r *FriendRequest) IsPending() bool {
	return r.Status == RequestStatusPending && !r.IsExpired()
}

// IsAccepted 检查请求是否已接受
func (r *FriendRequest) IsAccepted() bool {
	return r.Status == RequestStatusAccepted
}

// IsRejected 检查请求是否已拒绝
func (r *FriendRequest) IsRejected() bool {
	return r.Status == RequestStatusRejected
}

// IsCanceled 检查请求是否已取消
func (r *FriendRequest) IsCanceled() bool {
	return r.Status == RequestStatusCanceled
}

// GetAge 获取请求年龄
func (r *FriendRequest) GetAge() time.Duration {
	return time.Since(r.CreatedAt)
}

// GetTimeUntilExpiration 获取距离过期的时间
func (r *FriendRequest) GetTimeUntilExpiration() time.Duration {
	if r.IsExpired() {
		return 0
	}
	return time.Until(r.ExpiresAt)
}

const (
	DefaultRequestExpiration = 7 * 24 * time.Hour // 默认请求过期时间：7天
)

// generateRequestID 生成请求ID
func generateRequestID() string {
	return "req_" + randomString(16)
}
