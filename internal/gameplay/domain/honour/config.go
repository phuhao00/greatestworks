package honour

type Config struct {
	Id               uint32 `json:"id"`
	UpgradeCondition uint32 `json:"upgradeCondition"`
	ActiveCondition  uint32 `json:"activeCondition"`
	DeleteCondition  uint32 `json:"deleteCondition"`
	Desc             string `json:"desc"`
}
