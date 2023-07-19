package rank

type BlackList struct {
	RankId    uint32                `bson:"rankid"`
	BlackList map[string]*BlackInfo `bson:"blacklist"`
}

type BlackInfo struct {
	PlayerId uint64 `bson:"pid"`
	SetTime  int64  `bson:"stm"`
	Reason   string `bson:"reason"`
}
