package friend

type Info struct {
	UId      uint64
	ChatTime int64
	AddTime  int64
	Tag      string //备注
}

type Request struct {
	Userid  uint64 `json:"userid" bson:"userid"`   // 玩家ID
	OpTime  int64  `json:"opTime" bson:"opTime"`   // 操作时间
	AddType int32  `json:"addType" bson:"addType"` // 申请加好友的途径
}
