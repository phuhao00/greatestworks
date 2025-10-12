package email

import "time"

// Attachment 邮件附件实体
type Attachment struct {
	ID        string
	Type      AttachmentType
	ItemID    string
	Quantity  int64
	Claimed   bool
	ClaimedAt *time.Time
	ExpiresAt *time.Time
}

// AttachmentType 附件类型
type AttachmentType int

const (
	AttachmentTypeItem     AttachmentType = iota // 物品
	AttachmentTypeCurrency                       // 货币
	AttachmentTypeExp                            // 经验
	AttachmentTypeVIP                            // VIP时间
)

// NewAttachment 创建新附件
func NewAttachment(attachmentType AttachmentType, itemID string, quantity int64) *Attachment {
	return &Attachment{
		ID:       generateAttachmentID(),
		Type:     attachmentType,
		ItemID:   itemID,
		Quantity: quantity,
		Claimed:  false,
	}
}

// NewItemAttachment 创建物品附件
func NewItemAttachment(itemID string, quantity int64) *Attachment {
	return NewAttachment(AttachmentTypeItem, itemID, quantity)
}

// NewCurrencyAttachment 创建货币附件
func NewCurrencyAttachment(currencyType string, amount int64) *Attachment {
	return NewAttachment(AttachmentTypeCurrency, currencyType, amount)
}

// NewExpAttachment 创建经验附件
func NewExpAttachment(amount int64) *Attachment {
	return NewAttachment(AttachmentTypeExp, "exp", amount)
}

// Claim 领取附件
func (a *Attachment) Claim() error {
	if a.Claimed {
		return ErrAttachmentAlreadyClaimed
	}

	if a.IsExpired() {
		return ErrAttachmentExpired
	}

	a.Claimed = true
	claimedAt := time.Now()
	a.ClaimedAt = &claimedAt

	return nil
}

// IsExpired 检查附件是否已过期
func (a *Attachment) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*a.ExpiresAt)
}

// SetExpiration 设置过期时间
func (a *Attachment) SetExpiration(expiresAt time.Time) {
	a.ExpiresAt = &expiresAt
}

// IsClaimed 是否已领取
func (a *Attachment) IsClaimed() bool {
	return a.Claimed
}

// IsItem 是否为物品附件
func (a *Attachment) IsItem() bool {
	return a.Type == AttachmentTypeItem
}

// IsCurrency 是否为货币附件
func (a *Attachment) IsCurrency() bool {
	return a.Type == AttachmentTypeCurrency
}

// IsExp 是否为经验附件
func (a *Attachment) IsExp() bool {
	return a.Type == AttachmentTypeExp
}

// generateAttachmentID 生成附件ID
func generateAttachmentID() string {
	return "att_" + randomString(12)
}
