package email

import (
	"time"
)

// Email 邮件聚合根
type Email struct {
	ID          string
	SenderID    string
	ReceiverID  string
	Subject     string
	Content     string
	Type        EmailType
	Status      EmailStatus
	attachments []*Attachment
	SentAt      time.Time
	ReadAt      *time.Time
	ExpiresAt   *time.Time
	Version     int64
}

// EmailType 邮件类型
type EmailType int

const (
	EmailTypeNormal EmailType = iota // 普通邮件
	EmailTypeSystem                  // 系统邮件
	EmailTypeReward                  // 奖励邮件
	EmailTypeNotice                  // 通知邮件
)

// EmailStatus 邮件状态
type EmailStatus int

const (
	EmailStatusUnread EmailStatus = iota // 未读
	EmailStatusRead                       // 已读
	EmailStatusDeleted                    // 已删除
	EmailStatusExpired                    // 已过期
)

// NewEmail 创建新邮件
func NewEmail(senderID, receiverID, subject, content string, emailType EmailType) *Email {
	return &Email{
		ID:          generateEmailID(),
		SenderID:    senderID,
		ReceiverID:  receiverID,
		Subject:     subject,
		Content:     content,
		Type:        emailType,
		Status:      EmailStatusUnread,
		attachments: make([]*Attachment, 0),
		SentAt:      time.Now(),
		Version:     1,
	}
}

// NewSystemEmail 创建系统邮件
func NewSystemEmail(receiverID, subject, content string) *Email {
	email := NewEmail("system", receiverID, subject, content, EmailTypeSystem)
	// 系统邮件30天后过期
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	email.ExpiresAt = &expiresAt
	return email
}

// AddAttachment 添加附件
func (e *Email) AddAttachment(attachment *Attachment) error {
	if len(e.attachments) >= MaxAttachmentsPerEmail {
		return ErrTooManyAttachments
	}
	
	e.attachments = append(e.attachments, attachment)
	e.Version++
	
	return nil
}

// MarkAsRead 标记为已读
func (e *Email) MarkAsRead() error {
	if e.Status == EmailStatusDeleted {
		return ErrEmailDeleted
	}
	
	if e.IsExpired() {
		return ErrEmailExpired
	}
	
	if e.Status == EmailStatusUnread {
		e.Status = EmailStatusRead
		readAt := time.Now()
		e.ReadAt = &readAt
		e.Version++
	}
	
	return nil
}

// Delete 删除邮件
func (e *Email) Delete() error {
	if e.Status == EmailStatusDeleted {
		return ErrEmailAlreadyDeleted
	}
	
	e.Status = EmailStatusDeleted
	e.Version++
	
	return nil
}

// IsExpired 检查是否已过期
func (e *Email) IsExpired() bool {
	if e.Status == EmailStatusExpired {
		return true
	}
	
	if e.ExpiresAt != nil && time.Now().After(*e.ExpiresAt) {
		e.Status = EmailStatusExpired
		e.Version++
		return true
	}
	
	return false
}

// IsRead 是否已读
func (e *Email) IsRead() bool {
	return e.Status == EmailStatusRead
}

// IsDeleted 是否已删除
func (e *Email) IsDeleted() bool {
	return e.Status == EmailStatusDeleted
}

// GetAttachments 获取所有附件
func (e *Email) GetAttachments() []*Attachment {
	return e.attachments
}

// HasAttachments 是否有附件
func (e *Email) HasAttachments() bool {
	return len(e.attachments) > 0
}

// GetAge 获取邮件年龄
func (e *Email) GetAge() time.Duration {
	return time.Since(e.SentAt)
}

const (
	MaxAttachmentsPerEmail = 10 // 每封邮件最大附件数
)

// generateEmailID 生成邮件ID
func generateEmailID() string {
	return "email_" + randomString(16)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}