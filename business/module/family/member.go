package family

type Member struct {
	Position    MemberPosition `json:"position" bson:"position"`
	Id          uint64         `json:"id" bson:"id"`
	Nick        string         `json:"nick" bson:"nick"`
	Lv          uint32         `json:"lv" bson:"lv"`
	JoinTime    int64          `json:"joinTime" bson:"joinTime"`
	OffLineTime int64          `json:"offLineTime" bson:"offLineTime"`
}
