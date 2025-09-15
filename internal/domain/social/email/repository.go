package email

import "context"

// EmailRepository 邮件仓储接口
type EmailRepository interface {
	// 邮件相关
	SaveEmail(ctx context.Context, email *Email) error
	GetEmailByID(ctx context.Context, emailID string) (*Email, error)
	DeleteEmail(ctx context.Context, emailID string) error
	GetEmailsByReceiverID(ctx context.Context, receiverID string, limit int) ([]*Email, error)
	GetEmailsBySenderID(ctx context.Context, senderID string, limit int) ([]*Email, error)
	
	// 统计相关
	GetUnreadEmailCount(ctx context.Context, playerID string) (int, error)
	GetTotalEmailCount(ctx context.Context, playerID string) (int, error)
	
	// 查询相关
	GetEmailsByType(ctx context.Context, receiverID string, emailType EmailType, limit int) ([]*Email, error)
	GetEmailsByStatus(ctx context.Context, receiverID string, status EmailStatus, limit int) ([]*Email, error)
	GetExpiredEmails(ctx context.Context, limit int) ([]*Email, error)
	
	// 附件相关
	GetEmailsWithUnclaimed Attachments(ctx context.Context, receiverID string, limit int) ([]*Email, error)
	GetAttachmentByID(ctx context.Context, attach