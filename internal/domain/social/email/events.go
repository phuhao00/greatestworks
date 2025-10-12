package email

import "time"

// EmailEvent 邮件事件接口
type EmailEvent interface {
	GetEventType() string
	GetTimestamp() time.Time
	GetEmailID() string
}

// BaseEmailEvent 基础邮件事件
type BaseEmailEvent struct {
	EventType string
	Timestamp time.Time
	EmailID   string
}

func (e BaseEmailEvent) GetEventType() string {
	return e.EventType
}

func (e BaseEmailEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e BaseEmailEvent) GetEmailID() string {
	return e.EmailID
}

// EmailSentEvent 邮件发送事件
type EmailSentEvent struct {
	BaseEmailEvent
	SenderID   string
	ReceiverID string
	Subject    string
	Type       EmailType
}

// NewEmailSentEvent 创建邮件发送事件
func NewEmailSentEvent(emailID, senderID, receiverID, subject string, emailType EmailType) *EmailSentEvent {
	return &EmailSentEvent{
		BaseEmailEvent: BaseEmailEvent{
			EventType: "email.sent",
			Timestamp: time.Now(),
			EmailID:   emailID,
		},
		SenderID:   senderID,
		ReceiverID: receiverID,
		Subject:    subject,
		Type:       emailType,
	}
}

// EmailReadEvent 邮件阅读事件
type EmailReadEvent struct {
	BaseEmailEvent
	ReceiverID string
}

// NewEmailReadEvent 创建邮件阅读事件
func NewEmailReadEvent(emailID, receiverID string) *EmailReadEvent {
	return &EmailReadEvent{
		BaseEmailEvent: BaseEmailEvent{
			EventType: "email.read",
			Timestamp: time.Now(),
			EmailID:   emailID,
		},
		ReceiverID: receiverID,
	}
}

// EmailDeletedEvent 邮件删除事件
type EmailDeletedEvent struct {
	BaseEmailEvent
	ReceiverID string
}

// NewEmailDeletedEvent 创建邮件删除事件
func NewEmailDeletedEvent(emailID, receiverID string) *EmailDeletedEvent {
	return &EmailDeletedEvent{
		BaseEmailEvent: BaseEmailEvent{
			EventType: "email.deleted",
			Timestamp: time.Now(),
			EmailID:   emailID,
		},
		ReceiverID: receiverID,
	}
}

// AttachmentClaimedEvent 附件领取事件
type AttachmentClaimedEvent struct {
	BaseEmailEvent
	AttachmentID string
	ReceiverID   string
	Type         AttachmentType
	ItemID       string
	Quantity     int64
}

// NewAttachmentClaimedEvent 创建附件领取事件
func NewAttachmentClaimedEvent(emailID, attachmentID, receiverID, itemID string, attachmentType AttachmentType, quantity int64) *AttachmentClaimedEvent {
	return &AttachmentClaimedEvent{
		BaseEmailEvent: BaseEmailEvent{
			EventType: "email.attachment.claimed",
			Timestamp: time.Now(),
			EmailID:   emailID,
		},
		AttachmentID: attachmentID,
		ReceiverID:   receiverID,
		Type:         attachmentType,
		ItemID:       itemID,
		Quantity:     quantity,
	}
}

// SystemEmailSentEvent 系统邮件发送事件
type SystemEmailSentEvent struct {
	BaseEmailEvent
	ReceiverID      string
	Subject         string
	AttachmentCount int
}

// NewSystemEmailSentEvent 创建系统邮件发送事件
func NewSystemEmailSentEvent(emailID, receiverID, subject string, attachmentCount int) *SystemEmailSentEvent {
	return &SystemEmailSentEvent{
		BaseEmailEvent: BaseEmailEvent{
			EventType: "email.system.sent",
			Timestamp: time.Now(),
			EmailID:   emailID,
		},
		ReceiverID:      receiverID,
		Subject:         subject,
		AttachmentCount: attachmentCount,
	}
}
