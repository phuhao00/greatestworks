package friend

type Info struct {
	UId      uint64
	ChatTime int64
	AddTime  int64
	Tag      string //备注
}

type Request struct {
	Userid  uint64 // 玩家ID
	OpTime  int64  // 操作时间
	AddType int32  // 申请加好友的途径
}
