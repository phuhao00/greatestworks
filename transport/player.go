package transport

type Player struct {
	Uid        uint64   `bson:"uid"`
	NickName   string   `bson:"nickName"`
	Sex        int      `bson:"sex"`
	FriendList []uint64 `bson:"friendList"`
}
