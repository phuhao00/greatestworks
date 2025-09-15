package family

import "errors"

// 家族相关错误定义
var (
	// 家族相关错误
	ErrFamilyNotFound         = errors.New("家族不存在")
	ErrFamilyNameExists       = errors.New("家族名称已存在")
	ErrFamilyFull            = errors.New("家族成员已满")
	ErrFamilyNameTooShort    = errors.New("家族名称太短")
	ErrFamilyNameTooLong     = errors.New("家族名称太长")
	
	// 成员相关错误
	ErrMemberNotFound         = errors.New("成员不存在")
	ErrMemberAlreadyExists    = errors.New("成员已存在")
	ErrPlayerAlreadyInFamily  = errors.New("玩家已有家族")
	ErrCannotRemoveLeader     = errors.New("不能移除族长")
	ErrLeaderCannotLeave      = errors.New("族长不能离开家族")
	ErrCannotPromoteToLeader  = errors.New("不能提升为族长")
	
	// 权限相关错误
	ErrInsufficientPermission = errors.New("权限不足")
	ErrNotFamilyMember       = errors.New("不是家族成员")
	
	// 系统相关错误
	ErrSystemError           = errors.New("系统错误")
	ErrDatabaseError         = errors.New("数据库错误")
)