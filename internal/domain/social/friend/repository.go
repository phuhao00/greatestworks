package friend

import "context"

// FriendRepository 好友仓储接口
type FriendRepository interface {
	// 好友关系相关
	SaveFriendship(ctx context.Context, friendship *Friendship) error
	GetFriendship(ctx context.Context, playerID, friendID string) (*Friendship, error)
	GetFriendshipByID(ctx context.Context, friendshipID string) (*Friendship, error)
	GetFriendsByPlayerID(ctx context.Context, playerID string) ([]*Friendship, error)
	GetFriendCount(ctx context.Context, playerID string) (int, error)
	DeleteFriendship(ctx context.Context, friendshipID string) error
	
	// 好友请求相关
	SaveFriendRequest(ctx context.Context, request *FriendRequest) error
	GetFriendRequestByID(ctx context.Context, requestID string) (*FriendRequest, error)
	GetPendingRequest(ctx context.Context, fromPlayerID, toPlayerID string) (*FriendRequest, error)
	GetPendingRequestsByPlayerID(ctx context.Context, playerID string) ([]*FriendRequest, error)
	GetSentRequestsByPlayerID(ctx context.Context, playerID string) ([]*FriendRequest, error)
	DeleteFriendRequest(ctx context.Context, requestID string) error
	
	// 查询相关
	IsFriend(ctx context.Context, playerID, friendID string) (bool, error)
	IsBlocked(ctx context.Context, playerID, blockedPlayerID string) (bool, error)
	GetMutualFriends(ctx context.Context, playerID1, playerID2 string) ([]*Friendship, error)
}