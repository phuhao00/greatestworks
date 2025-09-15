package friend

import "errors"

// 好友相关错误定义
var (
	// 好友关系相关错误
	ErrFriendshipNotFound        = errors.New("好友关系不存在")
	ErrFriendshipAlreadyExists   = errors.New("好友关系已存在")
	ErrFriendshipDeleted         = errors.New("好友关系已删除")
	ErrFriendshipAlreadyDeleted  = errors.New("好友关系已经被删除")
	ErrInvalidFriendshipStatus   = errors.New("无效的好友关系状态")
	ErrAlreadyFriends           = errors.New("已经是好友关系")
	ErrNotBlocked               = errors.New("未被屏蔽")
	
	// 好友请求相关错误
	ErrRequestNotFound          = errors.New("好友请求不存在")
	ErrRequestAlreadyExists     = errors.New("好友请求已存在")
	ErrRequestExpired           = errors.New("好友请求已过期")
	ErrInvalidRequestStatus     = errors.New("无效的请求状态")
	
	// 权限相关错误
	ErrInsufficientPermission   = errors.New("权限不足")
	ErrCannotAddSelf           = errors.New("不能添加自己为好友")
	
	// 限制相关错误
	ErrTooManyFriends          = errors.New("好友数量已达上限")
	ErrTooManyRequests         = errors.New("请求数量已达上限")
	
	// 系统相关错误
	ErrSystemError             = errors.New("系统错误")
	ErrDatabaseError           = errors.New("数据库错误")
	ErrNetworkError            = errors.New("网络错误")