package template

type Conf struct {
	Id          uint32
	Description string
	StartTime   string
	EndTime     string
	Reward      string
	Category    string
	Param1      string
	Param2      string
	Param3      string
}

// Verify 字段检测合法
func (c *Conf) Verify() {
	//todo 一般性检查
	//todo 关联性检查
}

// Verify 字段检测合法
func (c *Conf) AfterAllVerify() {
	//todo 业务逻辑相关检查,可以在所有配置确定都加载好了，再执行此方法
}
