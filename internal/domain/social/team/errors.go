package team

import "errors"

// 队伍相关错误定义
var (
	// 队伍相关错误
	ErrTeamNotFound          = errors.New("队伍不存在")
	ErrTeamFull             = errors.New("队伍已满")
	ErrTeamNameTooShort     = errors.New("队伍名称太短")
	ErrTeamNameTooLong      = errors.New("队伍名称太长")
	ErrInvalidPassword      = errors.New("密码错误")
	
	// 成员相关错误
	ErrMemberNotFound        = errors.New("成员不存在")
	ErrMemberAlreadyExists   = errors.New("成员已存在")
	ErrPlayerAlreadyInTeam   = errors.New("玩家已有队伍")
	ErrCannotRemoveLeader    = errors.New("不能移除队长")
	ErrLeaderMustTransfer    = errors.New("队长必须先转让队长职位")
	
	// 权限相关错误
	ErrInsufficientPermission = errors.New("权限不足")
	ErrNotTeamMember         = errors.New("不是队伍成员")
	
	// 系统相关错误
	ErrSystemError           = errors.New("系统错误")
	ErrDatabaseError         = errors.New("