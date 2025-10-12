package chat

import "context"

// ChatRepository 聊天仓储接口
type ChatRepository interface {
	// 频道相关
	SaveChannel(ctx context.Context, channel *ChatChannel) error
	GetChannelByID(ctx context.Context, channelID string) (*ChatChannel, error)
	GetChannelByName(ctx context.Context, name string) (*ChatChannel, error)
	ChannelExistsByName(ctx context.Context, name string) (bool, error)
	DeleteChannel(ctx context.Context, channelID string) error
	GetChannelsByType(ctx context.Context, channelType ChannelType) ([]*ChatChannel, error)
	GetPlayerChannels(ctx context.Context, playerID string) ([]*ChatChannel, error)

	// 消息相关
	SaveMessage(ctx context.Context, message *Message) error
	GetMessageByID(ctx context.Context, messageID string) (*Message, error)
	GetMessagesByChannelID(ctx context.Context, channelID string, limit int) ([]*Message, error)
	DeleteMessage(ctx context.Context, messageID string) error
	GetMessagesByTimeRange(ctx context.Context, channelID string, start, end int64) ([]*Message, error)

	// 成员相关
	GetChannelMembers(ctx context.Context, channelID string) ([]*Member, error)
	GetMemberByPlayerID(ctx context.Context, channelID, playerID string) (*Member, error)
	UpdateMember(ctx context.Context, channelID string, member *Member) error
}
