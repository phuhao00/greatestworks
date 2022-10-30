package rank

type BlackList struct {
	Rank      uint32                `bson:"rank"`
	BlackList map[string]*BlackInfo `bson:"blacklist"`
}

type BlackInfo struct {
	PlayerId uint64 `bson:"pid"`
	SetTime  int64  `bson:"stm"`
	Reason   string `bson:"reason"`
}
