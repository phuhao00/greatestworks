package gm

type DealGmCommandRequest struct {
	Opt          string `json:"opt"`
	OptTime      int64  `json:"opt_time"`
	Target       string `json:"target"`
	SwitchPasswd string `json:"switch_passwd"`
	ChatType     string `json:"chat_type"`
	Modules      string `json:"modules"`
	Refresh      string `json:"refresh"`
	Content      string `json:"content"`
	Switch       string `json:"switch"`
	User         string `json:"user"`
	UidType      string `json:"uid_type"`
	IfAuto       string `json:"ifAuto"`
	SwitchKey    string `json:"switch_key"`
	BannedHours  string `json:"bannedHours"`
	TagType      string `json:"tagType"`
	UserTag      string `json:"userTag"`
	AfterMin     string `json:"after_min"`
	OpsMin       string `json:"ops_min"`
	KickoffUser  string `json:"kickoff_users"`
	OpsDesc      string `json:"ops_desc"`
	RepeatTimes  string `json:"repeat_times"`
	IntervalMin  string `json:"interval_min"`
	Zones        string `json:"zones"`
	InReason     string `json:"inReason"`
	OutReason    string `json:"outReason"`
	Cond         string `json:"cond"`
	Channels     string `json:"channels"`
	Remark       string `json:"remark"`
}
