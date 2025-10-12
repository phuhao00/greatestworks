package email

import (
	"context"
	"fmt"
)

// EmailService 邮件领域服务
type EmailService struct {
	emailRepo EmailRepository
}

// NewEmailService 创建邮件服务
func NewEmailService(emailRepo EmailRepository) *EmailService {
	return &EmailService{
		emailRepo: emailRepo,
	}
}

// SendEmail 发送邮件
func (s *EmailService) SendEmail(ctx context.Context, senderID, receiverID, subject, content string, emailType EmailType) (*Email, error) {
	// 验证邮件内容
	if err := s.validateEmail(subject, content); err != nil {
		return nil, err
	}

	// 检查接收者邮箱是否已满
	count, err := s.emailRepo.GetUnreadEmailCount(ctx, receiverID)
	if err != nil {
		return nil, fmt.Errorf("获取未读邮件数量失败: %w", err)
	}
	if count >= MaxEmailsPerPlayer {
		return nil, ErrMailboxFull
	}

	// 创建邮件
	email := NewEmail(senderID, receiverID, subject, content, emailType)

	// 保存邮件
	if err := s.emailRepo.SaveEmail(ctx, email); err != nil {
		return nil, fmt.Errorf("保存邮件失败: %w", err)
	}

	return email, nil
}

// SendSystemEmail 发送系统邮件
func (s *EmailService) SendSystemEmail(ctx context.Context, receiverID, subject, content string, attachments []*Attachment) (*Email, error) {
	// 创建系统邮件
	email := NewSystemEmail(receiverID, subject, content)

	// 添加附件
	for _, attachment := range attachments {
		if err := email.AddAttachment(attachment); err != nil {
			return nil, fmt.Errorf("添加附件失败: %w", err)
		}
	}

	// 保存邮件
	if err := s.emailRepo.SaveEmail(ctx, email); err != nil {
		return nil, fmt.Errorf("保存邮件失败: %w", err)
	}

	return email, nil
}

// ReadEmail 读取邮件
func (s *EmailService) ReadEmail(ctx context.Context, emailID, playerID string) (*Email, error) {
	// 获取邮件
	email, err := s.emailRepo.GetEmailByID(ctx, emailID)
	if err != nil {
		return nil, fmt.Errorf("获取邮件失败: %w", err)
	}
	if email == nil {
		return nil, ErrEmailNotFound
	}

	// 验证权限
	if email.ReceiverID != playerID {
		return nil, ErrInsufficientPermission
	}

	// 标记为已读
	if err := email.MarkAsRead(); err != nil {
		return nil, err
	}

	// 保存更改
	if err := s.emailRepo.SaveEmail(ctx, email); err != nil {
		return nil, fmt.Errorf("保存邮件失败: %w", err)
	}

	return email, nil
}

// ClaimAttachment 领取附件
func (s *EmailService) ClaimAttachment(ctx context.Context, emailID, attachmentID, playerID string) error {
	// 获取邮件
	email, err := s.emailRepo.GetEmailByID(ctx, emailID)
	if err != nil {
		return fmt.Errorf("获取邮件失败: %w", err)
	}
	if email == nil {
		return ErrEmailNotFound
	}

	// 验证权限
	if email.ReceiverID != playerID {
		return ErrInsufficientPermission
	}

	// 查找附件
	var targetAttachment *Attachment
	for _, attachment := range email.GetAttachments() {
		if attachment.ID == attachmentID {
			targetAttachment = attachment
			break
		}
	}

	if targetAttachment == nil {
		return ErrAttachmentNotFound
	}

	// 领取附件
	if err := targetAttachment.Claim(); err != nil {
		return err
	}

	// 保存更改
	if err := s.emailRepo.SaveEmail(ctx, email); err != nil {
		return fmt.Errorf("保存邮件失败: %w", err)
	}

	return nil
}

// DeleteEmail 删除邮件
func (s *EmailService) DeleteEmail(ctx context.Context, emailID, playerID string) error {
	// 获取邮件
	email, err := s.emailRepo.GetEmailByID(ctx, emailID)
	if err != nil {
		return fmt.Errorf("获取邮件失败: %w", err)
	}
	if email == nil {
		return ErrEmailNotFound
	}

	// 验证权限
	if email.ReceiverID != playerID {
		return ErrInsufficientPermission
	}

	// 删除邮件
	if err := email.Delete(); err != nil {
		return err
	}

	// 保存更改
	if err := s.emailRepo.SaveEmail(ctx, email); err != nil {
		return fmt.Errorf("保存邮件失败: %w", err)
	}

	return nil
}

// GetPlayerEmails 获取玩家邮件列表
func (s *EmailService) GetPlayerEmails(ctx context.Context, playerID string, limit int) ([]*Email, error) {
	return s.emailRepo.GetEmailsByReceiverID(ctx, playerID, limit)
}

// validateEmail 验证邮件
func (s *EmailService) validateEmail(subject, content string) error {
	if len(subject) == 0 {
		return ErrEmptySubject
	}
	if len(subject) > MaxSubjectLength {
		return ErrSubjectTooLong
	}
	if len(content) > MaxContentLength {
		return ErrContentTooLong
	}
	return nil
}

const (
	MaxEmailsPerPlayer = 100  // 每个玩家最大邮件数
	MaxSubjectLength   = 100  // 最大主题长度
	MaxContentLength   = 1000 // 最大内容长度
)
