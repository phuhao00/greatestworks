package chat

import "errors"

// 聊天相关错误定义
var (
	// 频道相关错误
	ErrChannelNotFound      = errors.New("频道不存在")
	ErrChannelAlreadyExists = errors.New("频道已存在")
	ErrChannelFull          = errors.New("频道已满")
	ErrChannelNameTooShort  = errors.New("频道名称太短")
	ErrChannelNameTooLong   = errors.New("频道名称太长")
	ErrInvalidChannelID     = errors.New("无效的频道ID")

	// 成员相关错误
	ErrMemberNotFound         = errors.New("成员不存在")
	ErrMemberAlreadyExists    = errors.New("成员已存在")
	ErrMemberMuted            = errors.New("成员被禁言")
	ErrInsufficientPermission = errors.New("权限不足")
	ErrSenderNotInChannel     = errors.New("发送者不在频道中")

	// 消息相关错误
	ErrMessageNotFound    = errors.New("消息不存在")
	ErrEmptyContent       = errors.New("消息内容为空")
	ErrMessageTooLong     = errors.New("消息内容过长")
	ErrInvalidSenderID    = errors.New("无效的发送者ID")
	ErrInvalidMessageType = errors.New("无效的消息类型")

	// 系统相关错误
	ErrSystemError   = errors.New("系统错误")
	ErrDatabaseError = errors.New("数据库错误")
	ErrNetworkError  = errors.New("网络错误")
)
