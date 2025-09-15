package email

import "errors"

// 邮件相关错误定义
var (
	// 邮件相关错误
	ErrEmailNotFound         = errors.New("邮件不存在")
	ErrEmailDeleted          = errors.New("邮件已删除")
	ErrEmailAlreadyDeleted   = errors.New("邮件已经被删除")
	ErrEmailExpired          = errors.New("邮件已过期")
	ErrMailboxFull          = errors.New("邮箱已满")
	ErrEmptySubject         = errors.New("邮件主题不能为空")
	ErrSubjectTooLong       = errors.New("邮件主题过长")
	ErrContentTooLong       = errors.New("邮件内容过长")
	
	// 附件相关错误
	ErrAttachmentNotFound        = errors.New("附件不存在")
	ErrAttachmentAlreadyClaimed  = errors.New("附件已被领取")
	ErrAttachmentExpired         = errors.New("附件已过期")
	ErrTooManyAttachments       = errors.New("附件数量过多")
	
	// 权限相关错误
	ErrInsufficientPermission   = errors.New("权限不足")
	ErrCannotSendToSelf        = errors.New("不能给自己发送邮件")
	
	// 系统相关错误
	ErrSystemError             = errors.New("系统错误")
	ErrDatabaseError           = errors.New("