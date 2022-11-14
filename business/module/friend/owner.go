package friend

import "github.com/phuhao00/network"

type Owner interface {
	GetFriendList()
	GetFriendInfo()
	AddFriend(packet *network.Message)
	DelFriend(packet *network.Message)
	GiveFriendItem()
	FriendAddApply()
	ManagerFriendApply()
}
