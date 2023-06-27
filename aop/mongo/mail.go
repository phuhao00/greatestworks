package mongo

type MailItem struct {
	ItemId uint32 `bson:"iid"`
	Num    int32  `bson:"num"`
}

type MailInfo struct {
	MUuid     uint64     `bson:"uid"`      //
	MailID    uint32     `bson:"mid"`      // 邮件表ID
	MType     int32      `bson:"type"`     // 邮件类型
	MSender   uint64     `bson:"send"`     // 邮件发送人
	MSNick    string     `bson:"nick"`     // 邮件发送人
	MSHead    uint32     `bson:"head"`     // 头像
	MStatus   uint32     `bson:"stat"`     // 邮件状态
	MItems    []MailItem `bson:"items"`    // 邮件发放物品
	MContent  string     `bson:"con"`      // 邮件内容
	MTime     int64      `bson:"time"`     // 邮件发送时间
	Template  uint32     `bson:"tpl"`      // 邮件模板
	Topic     string     `bson:"topic"`    // 邮件主题
	Paster    []uint32   `bson:"past"`     // 邮件贴纸
	Decorator []string   `bson:"deco"`     // 装饰
	Color     []float32  `bson:"color"`    // 颜色
	Validity  []int64    `bson:"validity"` // 生效日期
	RegTm     []int64    `bson:"reg"`      // 注册日期
	Chan      []string   `bson:"chan"`     // 渠道
	Version   string     `bson:"ver"`      // 版本号
}

type MailSystem struct {
	OwnerID uint64     `bson:"uid"`     // 玩家ID
	SendCnt int16      `bson:"scnt"`    // 当天发送次数
	Normal  []MailInfo `bson:"normal"`  // 正常邮件
	Collect []MailInfo `bson:"collect"` // 收藏邮件
	Recycle []MailInfo `bson:"recycle"` // 回收站
	History []uint64   `bson:"history"` // 历史邮件
}

func (t *MailSystem) C() string {
	return "MailSystem"
}

func (t *MailSystem) DB() string {
	return "greatest-work"
}
