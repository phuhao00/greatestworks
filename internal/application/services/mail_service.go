package services

import (
	"context"
	"fmt"
	"time"

	"greatestworks/internal/infrastructure/persistence"
)

// MailService 邮件服务
type MailService struct {
	mailRepo *persistence.MailRepository
}

// NewMailService 创建邮件服务
func NewMailService(mailRepo *persistence.MailRepository) *MailService {
	return &MailService{
		mailRepo: mailRepo,
	}
}

// SendMail 发送邮件
func (s *MailService) SendMail(ctx context.Context, receiverID int64, senderName, title, content string, attachments []persistence.DbAttachment, expireDays int) (int64, error) {
	// 生成邮件ID
	mailID := time.Now().UnixNano()

	hasItems := len(attachments) > 0
	expireAt := time.Now().AddDate(0, 0, expireDays)

	mail := &persistence.DbMail{
		MailID:      mailID,
		ReceiverID:  receiverID,
		SenderName:  senderName,
		Title:       title,
		Content:     content,
		IsRead:      false,
		HasItems:    hasItems,
		Attachments: attachments,
		ExpireAt:    expireAt,
	}

	if err := s.mailRepo.Create(ctx, mail); err != nil {
		return 0, fmt.Errorf("failed to send mail: %w", err)
	}

	return mailID, nil
}

// GetMails 获取邮件列表
func (s *MailService) GetMails(ctx context.Context, receiverID int64, limit int) ([]*persistence.DbMail, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // 默认50封
	}

	return s.mailRepo.FindByReceiverID(ctx, receiverID, limit)
}

// ReadMail 读取邮件
func (s *MailService) ReadMail(ctx context.Context, mailID int64) error {
	return s.mailRepo.MarkAsRead(ctx, mailID)
}

// DeleteMail 删除邮件
func (s *MailService) DeleteMail(ctx context.Context, mailID int64) error {
	return s.mailRepo.Delete(ctx, mailID)
}

// ClaimAttachments 领取附件
func (s *MailService) ClaimAttachments(ctx context.Context, mailID int64) ([]persistence.DbAttachment, error) {
	// TODO: 实现附件领取逻辑
	// 1. 检查背包空间
	// 2. 添加物品到背包
	// 3. 清空邮件附件或删除邮件

	return nil, nil
}

// SendSystemMail 发送系统邮件
func (s *MailService) SendSystemMail(ctx context.Context, receiverID int64, title, content string, attachments []persistence.DbAttachment) (int64, error) {
	return s.SendMail(ctx, receiverID, "System", title, content, attachments, 7)
}

// SendRewardMail 发送奖励邮件
func (s *MailService) SendRewardMail(ctx context.Context, receiverID int64, title string, itemRewards map[int32]int32) (int64, error) {
	// 构建附件
	attachments := make([]persistence.DbAttachment, 0, len(itemRewards))
	for itemID, count := range itemRewards {
		attachments = append(attachments, persistence.DbAttachment{
			ItemID: itemID,
			Count:  count,
		})
	}

	content := "Congratulations! You have received rewards."
	return s.SendMail(ctx, receiverID, "System", title, content, attachments, 7)
}

// CleanupExpiredMails 清理过期邮件（定时任务）
func (s *MailService) CleanupExpiredMails(ctx context.Context) error {
	return s.mailRepo.DeleteExpired(ctx)
}
